package provider

import (
	"caravan-cli/cli"
	"caravan-cli/terraform"
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
)

// GenericProvider is the generic implementation of the Provider interface and holds the Caravan config.
type GenericProvider struct {
	Caravan *cli.Config
}

// Bake performs the terraform apply to the caravan-baking repo.
func (g GenericProvider) Bake(ctx context.Context) error {
	t := terraform.New(g.Caravan.LogLevel)
	if err := t.Init(ctx, g.Caravan.WorkdirBaking); err != nil {
		return err
	}
	env := map[string]string{}
	if err := t.ApplyVarFile(ctx, filepath.Base(g.Caravan.WorkdirBakingVars), 1800*time.Second, env, "*"); err != nil {
		return err
	}
	return nil
}

// Depoly executes the corresponding terraform apply for the given layers/caravan repo.
func (g GenericProvider) Deploy(ctx context.Context, layer cli.DeployLayer) error {
	switch layer {
	case cli.Infrastructure:
		return GenericDeployInfra(ctx, g.Caravan, []string{"*"})
	case cli.Platform:
		return GenericDeployPlatform(ctx, g.Caravan, []string{"*"})
	case cli.ApplicationSupport:
		return GenericDeployApplicationSupport(ctx, g.Caravan, []string{"*"})
	default:
		return fmt.Errorf("unknown Deploy Layer")
	}
}

func GenericDeployInfra(ctx context.Context, c *cli.Config, targets []string) error {
	// Infra
	log.Info().Msgf("deploying infra")
	tf := terraform.New(c.LogLevel)
	if err := tf.Init(ctx, c.WorkdirInfra); err != nil {
		return err
	}
	env := map[string]string{}
	for _, target := range targets {
		if err := tf.ApplyVarFile(ctx, filepath.Base(c.WorkdirInfraVars), 1200*time.Second, env, target); err != nil {
			return fmt.Errorf("error doing terraform apply: %w", err)
		}
	}
	return nil
}

func GenericDeployPlatform(ctx context.Context, c *cli.Config, targets []string) error {
	// Platform
	log.Info().Msgf("deploying platform")
	tf := terraform.New(c.LogLevel)
	if err := tf.Init(ctx, c.WorkdirPlatform); err != nil {
		return err
	}
	env := map[string]string{
		"VAULT_TOKEN": c.VaultRootToken,
		"NOMAD_TOKEN": c.NomadToken,
	}
	for _, target := range targets {
		if err := tf.ApplyVarFile(ctx, filepath.Base(c.WorkdirPlatformVars), 600*time.Second, env, target); err != nil {
			return fmt.Errorf("error doing terraform apply: %w", err)
		}
	}
	return nil
}

func GenericDeployApplicationSupport(ctx context.Context, c *cli.Config, targets []string) error {
	// Application support
	log.Info().Msgf("deploying application")
	tf := terraform.New(c.LogLevel)
	if err := tf.Init(ctx, c.WorkdirApplication); err != nil {
		return err
	}
	env := map[string]string{
		"VAULT_TOKEN": c.VaultRootToken,
		"NOMAD_TOKEN": c.NomadToken,
	}
	for _, target := range targets {
		if err := tf.ApplyVarFile(ctx, filepath.Base(c.WorkdirApplicationVars), 600*time.Second, env, target); err != nil {
			return fmt.Errorf("error doing terraform apply: %w", err)
		}
	}
	return nil
}

func (g GenericProvider) Destroy(ctx context.Context, layer cli.DeployLayer) error {
	switch layer {
	case cli.Infrastructure:
		return g.cleanInfra(ctx)
	case cli.Platform:
		return g.cleanPlatform(ctx)
	case cli.ApplicationSupport:
		return g.cleanApplication(ctx)
	default:
		return fmt.Errorf("cannot destroy unknown deploy layer: %d", layer)
	}
}

func (g GenericProvider) cleanInfra(ctx context.Context) (err error) {
	log.Info().Msgf("removing terraform infrastructure")
	tf := terraform.New(g.Caravan.LogLevel)
	err = tf.Init(ctx, g.Caravan.WorkdirInfra)
	if err != nil {
		return err
	}
	env := map[string]string{}
	if err := tf.Destroy(ctx, filepath.Base(g.Caravan.WorkdirInfraVars), env); err != nil {
		log.Error().Msgf("error during destroy of cloud resources: %s", err)
		if !g.Caravan.Force {
			return err
		}
	}
	return nil
}

func (g GenericProvider) cleanPlatform(ctx context.Context) (err error) {
	log.Info().Msgf("removing terraform platform")
	tf := terraform.New(g.Caravan.LogLevel)
	err = tf.Init(ctx, g.Caravan.WorkdirPlatform)
	if err != nil {
		return err
	}
	env := map[string]string{
		"VAULT_TOKEN": g.Caravan.VaultRootToken,
		"NOMAD_TOKEN": g.Caravan.NomadToken,
	}
	if err := tf.Destroy(ctx, filepath.Base(g.Caravan.WorkdirPlatformVars), env); err != nil {
		log.Error().Msgf("error during destroy of cloud resources: %s", err)
		if !g.Caravan.Force {
			return err
		}
	}
	return nil
}

func (g GenericProvider) cleanApplication(ctx context.Context) (err error) {
	log.Info().Msgf("removing terraform application")
	tf := terraform.New(g.Caravan.LogLevel)
	err = tf.Init(ctx, g.Caravan.WorkdirApplication)
	if err != nil {
		return err
	}
	env := map[string]string{
		"VAULT_TOKEN": g.Caravan.VaultRootToken,
		"NOMAD_TOKEN": g.Caravan.NomadToken,
	}
	if err := tf.Destroy(ctx, filepath.Base(g.Caravan.WorkdirApplicationVars), env); err != nil {
		log.Error().Msgf("error during destroy of cloud resources: %s", err)
		if !g.Caravan.Force {
			return err
		}
	}
	return nil
}

func (g GenericProvider) Status(ctx context.Context) error {
	panic("implement me")
}
