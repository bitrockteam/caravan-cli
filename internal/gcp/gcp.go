package gcp

import (
	"caravan/internal/caravan"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type GCP struct {
	caravan.GenericProvider
	caravan.GenericBake
	caravan.GenericStatus
	caravan.GenericDestroy
	caravan.GenericDeploy
}

func New(c *caravan.Config) (g GCP, err error) {
	g = GCP{}
	g.Caravan = c
	if err := g.ValidateConfiguration(); err != nil {
		return g, err
	}

	// TODO: more setup?

	return g, nil
}

func (g GCP) GetTemplates() ([]caravan.Template, error) {
	return []caravan.Template{
		{
			Name: "baking-vars",
			Text: bakingTfVarsTmpl,
			Path: g.Caravan.WorkdirBakingVars,
		},
		{
			Name: "infra-vars",
			Text: infraTfVarsTmpl,
			Path: g.Caravan.WorkdirInfraVars,
		},
		{
			Name: "platform-vars",
			Text: platformTfVarsTmpl,
			Path: g.Caravan.WorkdirPlatformVars,
		},
		{
			Name: "application-vars",
			Text: applicationTfVarsTmpl,
			Path: g.Caravan.WorkdirApplicationVars,
		},
		{
			Name: "infra-backend",
			Text: infraBackendTmpl,
			Path: g.Caravan.WorkdirInfraBackend,
		},
		{
			Name: "platform-backend",
			Text: platformBackendTmpl,
			Path: g.Caravan.WorkdirPlatformBackend,
		},
		{
			Name: "application-backend",
			Text: applicationSupportBackendTmpl,
			Path: g.Caravan.WorkdirApplicationBackend,
		},
	}, nil
}

func (g GCP) ValidateConfiguration() error {
	m, err := regexp.MatchString("^[-0-9A-Za-z]{6,15}$", g.Caravan.Name)
	if err != nil {
		return err
	}
	if !m {
		return fmt.Errorf("project name not compliant: must be between 6 and 15 characters long, only alphanumerics and hypens (-) are allowed: %s", g.Caravan.Name)
	}
	if strings.Index(g.Caravan.Name, "-") == 0 {
		return fmt.Errorf("project name not compliant: cannot start with hyphen (-): %s", g.Caravan.Name)
	}
	return nil
}

func (g GCP) InitProvider() error {
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

	if err := g.CreateStateStore(g.Caravan.StateStoreName); err != nil {
		return err
	}

	return nil
}

func (g GCP) CleanProvider() error {
	if err := g.DeleteServiceAccount("terraform"); err != nil {
		return err
	}
	if err := g.DeleteStateStore(g.Caravan.StateStoreName); err != nil {
		return err
	}

	return nil
}
