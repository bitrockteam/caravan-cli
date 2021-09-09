package aws

import (
	"caravan/internal/caravan"
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type AWS struct {
	caravan.GenericProvider
	caravan.GenericBake
	caravan.GenericStatus
	caravan.GenericDestroy
	AWSConfig aws.Config
}

func New(c *caravan.Config) (a AWS, err error) {
	a = AWS{}
	a.Caravan = c
	if err := a.ValidateConfiguration(); err != nil {
		return a, err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if a.Caravan.Region != "" {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(a.Caravan.Region),
		)
	}
	if err != nil {
		return a, err
	}
	if cfg.Region == "" {
		return a, fmt.Errorf("please provide a region")
	}
	a.Caravan.Region = cfg.Region
	a.AWSConfig = cfg

	return a, nil
}

func (a AWS) GetTemplates() ([]caravan.Template, error) {
	return []caravan.Template{
		{
			Name: "baking-vars",
			Text: bakingTfVarsTmpl,
			Path: a.Caravan.WorkdirBakingVars,
		},
		{
			Name: "infra-vars",
			Text: infraTfVarsTmpl,
			Path: a.Caravan.WorkdirInfraVars,
		},
		{
			Name: "infra-backend",
			Text: infraBackendTmpl,
			Path: a.Caravan.WorkdirInfraBackend,
		},
		{
			Name: "platform-backend",
			Text: platformBackendTmpl,
			Path: a.Caravan.WorkdirPlatformBackend,
		},
		{
			Name: "platform-vars",
			Text: paltformTvVarsTmpl,
			Path: a.Caravan.WorkdirPlatformVars,
		},
		{
			Name: "application-vars",
			Text: applicationTfVarsTmpl,
			Path: a.Caravan.WorkdirApplicationVars,
		},
		{
			Name: "application-backend",
			Text: applicationSupportBackendTmpl,
			Path: a.Caravan.WorkdirApplicationBackend,
		},
	}, nil
}

func (a AWS) ValidateConfiguration() error {
	m, err := regexp.MatchString("^[-0-9A-Za-z]{3,12}$", a.Caravan.Name)
	if err != nil {
		return err
	}
	if !m {
		return fmt.Errorf("project name not compliant: must be between 3 and 12 character long, only alphanumerics and hypens (-) are allowed: %s", a.Caravan.Name)
	}
	if strings.Index(a.Caravan.Name, "-") == 0 {
		return fmt.Errorf("project name not compliant: cannot start with hyphen (-): %s", a.Caravan.Name)
	}
	return nil
}

func (a AWS) InitProvider() error {
	if err := a.CreateStateStore(a.Caravan.StateStoreName); err != nil {
		fmt.Printf("failed to create state store")
		return err
	}
	if err := a.CreateLock(a.Caravan.LockName); err != nil {
		fmt.Printf("failed to create lock")
		return err
	}
	return nil
}

func (a AWS) Init() error {
	return nil
}

func (a AWS) CleanProvider() error {
	fmt.Printf("removing terraform state and locking structures\n")

	if a.Caravan.Force {
		fmt.Printf("emptying bucket %s\n", a.Caravan.Name+"-caravan-terraform-state")
		err := a.EmptyStateStore(a.Caravan.Name + "-caravan-terraform-state")
		if err != nil {
			return fmt.Errorf("error emptying: %w", err)
		}
	}

	if err := a.DeleteStateStore(a.Caravan.Name + "-caravan-terraform-state"); err != nil {
		return err
	}

	if err := a.DeleteLock(a.Caravan.Name + "-caravan-terraform-state-lock"); err != nil {
		return err
	}

	return nil
}

func (a AWS) GenerateConfig() error {
	return nil
}

func (a AWS) Deploy(layer caravan.DeployLayer) error {
	switch layer {
	case caravan.Infrastructure:
		return caravan.GenericDeployInfra(a.Caravan, []string{"aws_lb.hashicorp_alb", "*"})
	case caravan.Platform:
		return caravan.GenericDeployPlatform(a.Caravan, []string{"*"})
	case caravan.ApplicationSupport:
		return caravan.GenericDeployApplicationSupport(a.Caravan, []string{"*"})
	default:
		return fmt.Errorf("unknown Deploy Layer")
	}
}
