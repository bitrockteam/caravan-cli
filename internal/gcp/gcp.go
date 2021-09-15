package gcp

import (
	"caravan/internal/caravan"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/iam/v1"
	"google.golang.org/grpc/status"
)

type GCP struct {
	Caravan   caravan.Config
	Templates []caravan.Template
}

func New(c caravan.Config) (g GCP, err error) {
	if err := validate(c); err != nil {
		return g, err
	}
	g = GCP{Caravan: c}

	g.Templates = loadTemplates(g)
	return g, nil
}

func (g GCP) Init() error {
	// assume that the project and billing account are already available
	/*
		if err := g.CreateProject(g.Caravan.Name, g.Caravan.GCPOrgID); err != nil {
			return err
		}

		if err := g.SetBillingAccount(g.Caravan.Name, g.Caravan.GCPOrgID); err != nil {
			return err
		}
	*/

	if err := g.CreateServiceAccount("terraform"); err != nil {
		return err
	}

	// permissions for the terraform service account on the current project
	if err := g.AddPolicyBinding("projects", g.Caravan.Name, "terraform", "roles/owner"); err != nil {
		return err
	}
	if err := g.AddPolicyBinding("projects", g.Caravan.Name, "terraform", "roles/storage.admin"); err != nil {
		return err
	}

	// permission for the terraform service account on the parent project
	if err := g.AddPolicyBinding("projects", g.Caravan.ParentProject, "terraform", "roles/compute.imageUser"); err != nil {
		return err
	}
	if err := g.AddPolicyBinding("projects", g.Caravan.ParentProject, "terraform", "roles/dns.admin"); err != nil {
		return err
	}
	if err := g.AddPolicyBinding("projects", g.Caravan.ParentProject, "terraform", "roles/compute.networkAdmin"); err != nil {
		return err
	}
	if err := g.AddPolicyBinding("projects", g.Caravan.ParentProject, "terraform", "roles/iam.serviceAccountUser"); err != nil {
		return err
	}

	// permission for the current user on the parent project
	if err := g.AddPolicyBinding("projects", g.Caravan.ParentProject, "andrea.simonini@bitrock.it", "roles/iam.serviceAccountUser"); err != nil {
		return err
	}

	// create keys for service account
	kb64, err := g.CreateServiceAccountKeys("terraform", "terraform-sa-keys")
	if err != nil {
		return err
	}
	k, err := base64.StdEncoding.DecodeString(kb64)
	if err != nil {
		return err
	}
	if err := os.WriteFile(g.Caravan.WorkdirInfra+"/."+g.Caravan.Name+"-terraform-sa-key.json", k, 0600); err != nil {
		return err
	}

	return nil
}

func (g GCP) Clean() error {
	fmt.Printf("cleaning: NOP\n")
	return nil
}

func (g GCP) CreateStateStore(name string) error {
	fmt.Printf("creating bucket %s on project: %s\n", name, g.Caravan.Name)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	storageLocation := &storage.BucketAttrs{
		Location: g.Caravan.Region,
	}
	bucket := client.Bucket(name)
	if err := bucket.Create(ctx, g.Caravan.Name, storageLocation); err != nil {
		s, _ := status.FromError(err)
		if strings.Contains(s.Message(), "You already own this bucket") {
			return nil
		}
		return fmt.Errorf("error during bucket %s create: %w", name, err)
	}
	return nil
}

func (g GCP) DeleteStateStore(name string) error {
	fmt.Printf("deleting bucket %s on project: %s\n", name, g.Caravan.Name)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	bucket := client.Bucket(name)
	if err := bucket.Delete(ctx); err != nil {
		s, _ := status.FromError(err)
		if strings.Contains(s.Message(), "notFound") {
			return nil
		}
		return fmt.Errorf("error during bucket %s delete: %w", name, err)
	}
	return nil
}

func (g GCP) EmptyStateStore(name string) error {
	fmt.Printf("empty state store: NOP\n")
	return nil
}
func (g GCP) CreateLock(name string) error {
	fmt.Printf("create lock: NOP\n")
	return nil
}

func (g GCP) DeleteLock(name string) error {
	fmt.Printf("delete lock: NOP\n")
	return nil
}

// CreateProject creates a project in GCP and waits for the completion.
func (g GCP) CreateProject(id, orgID string) error {
	ctx := context.Background()
	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return fmt.Errorf("unable to get resourcemanager: %w", err)
	}

	p := &cloudresourcemanager.Project{
		ProjectId:   id,
		Parent:      orgID,
		DisplayName: id,
	}
	op, err := cloudresourcemanagerService.Projects.Create(p).Context(ctx).Do()
	if err != nil {
		s, _ := status.FromError(err)
		if strings.Contains(s.Message(), "alreadyExists") {
			p, err := g.GetProject(id, orgID)
			if err != nil || p == nil {
				return fmt.Errorf("unable to find already existing project %s - %w", id, err)
			}
			if p.State == "ACTIVE" {
				return nil
			}
			return fmt.Errorf("project %s already existing and not in active state: %s", id, p.State)
		}
		return err
	}

	for i := 0; i < 10; i++ {
		resp, err := cloudresourcemanagerService.Operations.Get(op.Name).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("error getting operation %s/%s: %w", id, orgID, err)
		}
		if resp.Done {
			_, _ = g.GetProject(id, orgID)
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("timed out creating project %s", id)
}

// DeleteProject deletes a project from its project id.
func (g GCP) DeleteProject(name, organization string) error {
	fmt.Printf("deleting project: %s\n", name)
	ctx := context.Background()

	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return fmt.Errorf("unable to get resourcemanager: %w", err)
	}

	p, err := g.GetProject(name, organization)
	if err != nil {
		return err
	}

	if p == nil {
		return nil
	}

	if p.State == "DELETE_REQUESTED" {
		return nil
	}

	_, err = cloudresourcemanagerService.Projects.Delete(p.Name).Context(ctx).Do()
	if err != nil {
		return err
	}

	// check project name
	p, err = g.GetProject(name, organization)
	if err != nil {
		return err
	}
	if p == nil {
		return nil
	}

	if p.State == "DELETE_REQUESTED" {
		return nil
	}
	return nil
}

func (g GCP) GetProject(name, organization string) (p *cloudresourcemanager.Project, err error) {
	ctx := context.Background()

	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return p, fmt.Errorf("unable to get resourcemanager: %w", err)
	}
	// check project name
	q := fmt.Sprintf("id:%s parent:%s", name, organization)
	resp, err := cloudresourcemanagerService.Projects.Search().Query(q).Context(ctx).Do()
	if err != nil {
		return p, err
	}

	if len(resp.Projects) == 0 {
		return p, nil
	}

	if len(resp.Projects) != 1 {
		return p, fmt.Errorf("unable to uniquely identify project: %s (%d)", name, len(resp.Projects))
	}
	// fmt.Printf("found project: %v\n", resp.Projects[0])
	return resp.Projects[0], nil
}

func (g GCP) SetBillingAccount(name, bai string) (err error) {
	fmt.Printf("Setting Billing account: %s\n", bai)
	ctx := context.Background()

	cloudbillingservice, err := cloudbilling.NewService(ctx)
	if err != nil {
		return fmt.Errorf("unable to get billingservice: %w", err)
	}

	pbi, err := cloudbillingservice.Projects.GetBillingInfo(fmt.Sprintf("projects/%s", name)).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("unable to get project billing info: %w", err)
	}

	pbi.BillingAccountName = fmt.Sprintf("billingAccounts/%s", bai)
	_, err = cloudbillingservice.Projects.UpdateBillingInfo(fmt.Sprintf("projects/%s", name), pbi).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("unable to update project billing info: %w", err)
	}
	return nil
}

func (g GCP) CreateServiceAccount(name string) (err error) {
	fmt.Printf("Create service account: %s\n", name)
	ctx := context.Background()

	iamservice, err := iam.NewService(ctx)
	if err != nil {
		return fmt.Errorf("unable to get iam service: %w", err)
	}

	sar := iam.CreateServiceAccountRequest{
		AccountId: name,
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: name,
		},
	}
	_, err = iamservice.Projects.ServiceAccounts.Create(fmt.Sprintf("projects/%s", g.Caravan.Name), &sar).Context(ctx).Do()
	if err != nil {
		s, _ := status.FromError(err)
		if strings.Contains(s.Message(), "alreadyExists") {
			return nil
		}
		return fmt.Errorf("unable to create service account %s: %w", name, err)
	}

	return nil
}

func (g GCP) DeleteServiceAccount(name string) (err error) {
	fmt.Printf("delete service account: %s\n", name)
	ctx := context.Background()

	iamservice, err := iam.NewService(ctx)
	if err != nil {
		return fmt.Errorf("unable to get iam service: %w", err)
	}

	_, err = iamservice.Projects.ServiceAccounts.Delete(fmt.Sprintf("projects/%s/serviceAccounts/%s@%s.iam.gserviceaccount.com", g.Caravan.Name, name, g.Caravan.Name)).Context(ctx).Do()
	if err != nil {
		s, _ := status.FromError(err)
		if strings.Contains(s.Message(), "notFound") {
			return nil
		}
		return fmt.Errorf("unable to delete service account %s: %w", name, err)
	}

	return nil
}

func (g GCP) CreateServiceAccountKeys(sa, name string) (key string, err error) {
	fmt.Printf("create service account keys: %s\n", name)
	ctx := context.Background()

	iamservice, err := iam.NewService(ctx)
	if err != nil {
		return key, fmt.Errorf("unable to get iam service: %w", err)
	}

	sak, err := iamservice.Projects.ServiceAccounts.Keys.Create(fmt.Sprintf("projects/%s/serviceAccounts/%s@%s.iam.gserviceaccount.com", g.Caravan.Name, sa, g.Caravan.Name), &iam.CreateServiceAccountKeyRequest{}).Context(ctx).Do()
	if err != nil {
		s, _ := status.FromError(err)
		if strings.Contains(s.Message(), "alreadyExists") {
			return sak.PrivateKeyData, nil
		}
		return key, fmt.Errorf("unable to create service account keys %s: %w", name, err)
	}

	return sak.PrivateKeyData, nil
}

func (g GCP) AddPolicyBinding(resource, name, sa, role string) error {
	fmt.Printf("add policy binding: %s %s@%s/%s\n", sa, role, resource, name)
	ctx := context.Background()

	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return err
	}

	policy, err := g.GetPolicyBinding(resource, name, sa)
	if err != nil {
		return err
	}
	policy.Bindings = append(policy.Bindings, &cloudresourcemanager.Binding{Role: role})
	rb := &cloudresourcemanager.SetIamPolicyRequest{
		Policy: policy,
	}

	_, err = cloudresourcemanagerService.Projects.SetIamPolicy(fmt.Sprintf("%s/%s", resource, name), rb).Context(ctx).Do()
	if err != nil {
		return err
	}

	return nil
}

func (g GCP) GetPolicyBinding(resource, name, sa string) (policy *cloudresourcemanager.Policy, err error) {
	fmt.Printf("get policy binding: %s@%s/%s\n", sa, resource, name)
	ctx := context.Background()

	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return policy, err
	}

	gir := &cloudresourcemanager.GetIamPolicyRequest{}

	policy, err = cloudresourcemanagerService.Projects.GetIamPolicy(fmt.Sprintf("%s/%s", resource, name), gir).Context(ctx).Do()
	if err != nil {
		return policy, err
	}

	return policy, nil
}

func validate(c caravan.Config) error {
	m, err := regexp.MatchString("^[-0-9A-Za-z]{6,15}$", c.Name)
	if err != nil {
		return err
	}
	if !m {
		return fmt.Errorf("project name not compliant: must be between 6 and 15 characters long, only alphanumerics and hypens (-) are allowed: %s", c.Name)
	}
	if strings.Index(c.Name, "-") == 0 {
		return fmt.Errorf("project name not compliant: cannot start with hyphen (-): %s", c.Name)
	}
	return nil
}
