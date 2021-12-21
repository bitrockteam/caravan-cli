// Amazon Web Services provider.
package aws

import (
	"caravan-cli/cli"
	"caravan-cli/provider"
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type AWS struct {
	provider.GenericProvider
	AWSConfig aws.Config
}

func New(ctx context.Context, c *cli.Config) (a AWS, err error) {
	a = AWS{}
	a.Caravan = c

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return a, err
	}
	if a.Caravan.Region == "" {
		a.Caravan.Region = cfg.Region
	}

	if err := a.ValidateConfiguration(ctx); err != nil {
		return a, err
	}

	cfg, err = config.LoadDefaultConfig(ctx,
		config.WithRegion(a.Caravan.Region),
	)
	if err != nil {
		return a, err
	}

	a.AWSConfig = cfg

	return a, nil
}

func (a AWS) GetTemplates(ctx context.Context) ([]cli.Template, error) {
	return []cli.Template{
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
			Text: platformTfVarsTmpl,
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

func (a AWS) ValidateConfiguration(ctx context.Context) error {
	// check project name
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
	// check valid region
	if a.Caravan.Region == "" {
		return fmt.Errorf("please provide a region configuration")
	}
	if _, err := net.LookupIP(fmt.Sprintf("ec2.%s.amazonaws.com", a.Caravan.Region)); err != nil {
		return fmt.Errorf("region %s not allowed: %w", a.Caravan.Region, err)
	}
	return nil
}

func (a AWS) InitProvider(ctx context.Context) error {
	if err := a.CreateStateStore(ctx, a.Caravan.StateStoreName); err != nil {
		log.Error().Msgf("failed to create state store")
		return err
	}
	if err := a.CreateLock(ctx, a.Caravan.LockName); err != nil {
		log.Error().Msgf("failed to create lock")
		return err
	}
	return nil
}

func (a AWS) CleanProvider(ctx context.Context) error {
	log.Info().Msgf("removing terraform state and locking structures")

	if a.Caravan.Force {
		log.Info().Msgf("emptying bucket %s", a.Caravan.StateStoreName)
		err := a.EmptyStateStore(ctx, a.Caravan.StateStoreName)
		if err != nil {
			return fmt.Errorf("error emptying: %w", err)
		}
	}

	if err := a.DeleteStateStore(ctx, a.Caravan.StateStoreName); err != nil {
		return err
	}

	if err := a.DeleteLock(ctx, a.Caravan.Name+"-caravan-terraform-state-lock"); err != nil {
		return err
	}

	return nil
}

func (a AWS) Deploy(ctx context.Context, layer cli.DeployLayer) error {
	switch layer {
	case cli.Infrastructure:
		return provider.GenericDeployInfra(ctx, a.Caravan, []string{"aws_lb.hashicorp_alb", "*"})
	case cli.Platform:
		return provider.GenericDeployPlatform(ctx, a.Caravan, []string{"*"})
	case cli.ApplicationSupport:
		return provider.GenericDeployApplicationSupport(ctx, a.Caravan, []string{"*"})
	default:
		return fmt.Errorf("unknown Deploy Layer")
	}
}
