package gcp

import (
	"caravan/internal/caravan"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
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

	if c.UserEmail == "" {
		home := os.Getenv("HOME")
		u, err := g.GetUserEmail(filepath.Join(home, ".config/gcloud/configurations/config_default"))
		if err != nil {
			return g, err
		}
		c.UserEmail = u
	}
	g.Caravan = c
	if err := g.ValidateConfiguration(); err != nil {
		return g, err
	}

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
	// check project name
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

	// check valid region
	if g.Caravan.Region != "europe-west6" {
		return fmt.Errorf("gcp region %s not supported", g.Caravan.Region)
	}
	return nil
}

func (g GCP) InitProvider() error {
	if err := g.CreateServiceAccount(g.Caravan.ServiceAccount); err != nil {
		return err
	}

	// permissions for the terraform service account on the current project
	if err := g.AddPolicyBinding("projects", g.Caravan.Name, g.Caravan.ServiceAccount, "roles/owner"); err != nil {
		return err
	}
	if err := g.AddPolicyBinding("projects", g.Caravan.Name, g.Caravan.ServiceAccount, "roles/storage.admin"); err != nil {
		return err
	}

	// permission for the terraform service account on the parent project
	if err := g.AddPolicyBinding("projects", g.Caravan.ParentProject, g.Caravan.ServiceAccount, "roles/compute.imageUser"); err != nil {
		return err
	}
	if err := g.AddPolicyBinding("projects", g.Caravan.ParentProject, g.Caravan.ServiceAccount, "roles/dns.admin"); err != nil {
		return err
	}
	if err := g.AddPolicyBinding("projects", g.Caravan.ParentProject, g.Caravan.ServiceAccount, "roles/compute.networkAdmin"); err != nil {
		return err
	}
	if err := g.AddPolicyBinding("projects", g.Caravan.ParentProject, g.Caravan.ServiceAccount, "roles/iam.serviceAccountUser"); err != nil {
		return err
	}

	// permission for the current user on the parent project
	if err := g.AddPolicyBinding("projects", g.Caravan.ParentProject, g.Caravan.UserEmail, "roles/iam.serviceAccountUser"); err != nil {
		return err
	}

	// create keys for service account
	kb64, err := g.CreateServiceAccountKeys(g.Caravan.ServiceAccount, g.Caravan.ServiceAccount+"-sa-keys")
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
	if err := g.DeleteServiceAccount(g.Caravan.ServiceAccount); err != nil {
		return err
	}
	if err := g.DeleteStateStore(g.Caravan.StateStoreName); err != nil {
		return err
	}

	return nil
}
