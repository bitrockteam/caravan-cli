package gcp

import (
	"caravan/internal/caravan"
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"google.golang.org/api/cloudresourcemanager/v3"
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
	return GCP{Caravan: c}, nil
}

func (g GCP) Init() error {
	fmt.Printf("creating project: %s,%s\n", g.Caravan.GCPOrgID, g.Caravan.Name)
	if err := g.CreateProject(g.Caravan.Name, g.Caravan.GCPOrgID); err != nil {
		return err
	}
	return nil
}

func (g GCP) Clean() error {
	if err := g.DeleteProject(g.Caravan.Name, g.Caravan.GCPOrgID); err != nil {
		return err
	}
	return nil
}

func (g GCP) CreateBucket(name string) error {
	fmt.Printf("NOP\n")
	return nil
}

func (g GCP) DeleteBucket(name string) error {
	fmt.Printf("NOP\n")
	return nil
}
func (g GCP) EmptyBucket(name string) error {
	fmt.Printf("NOP\n")
	return nil
}
func (g GCP) CreateLockTable(name string) error {
	fmt.Printf("NOP\n")
	return nil
}

func (g GCP) DeleteLockTable(name string) error {
	fmt.Printf("NOP\n")
	return nil
}

func (g GCP) GenerateConfig() error {
	fmt.Printf("NOP\n")
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
			return nil
		}
		return err
	}

	for i := 0; i < 10; i++ {
		resp, err := cloudresourcemanagerService.Operations.Get(op.Name).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("error getting operation %s/%s: %w", id, orgID, err)
		}
		if resp.Done {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("timed out creating project %s", id)
}

// DeleteProject deletes a project from its project id.
func (g GCP) DeleteProject(name, organization string) error {
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

	return resp.Projects[0], nil
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
