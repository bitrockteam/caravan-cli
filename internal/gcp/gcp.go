package gcp

import (
	"caravan/internal/caravan"
	"fmt"
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
	panic("implement me")
	// return []caravan.Template{}, nil
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
	panic("implement me")
}

func (g GCP) CleanProvider() error {
	//TODO: delete resources
	panic("implement me")
}
