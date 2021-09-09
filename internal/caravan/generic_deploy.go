package caravan

import (
	tf "caravan/internal/terraform"
	"fmt"
	"path/filepath"
	"time"
)

type GenericDeploy struct {
	GenericProvider
}

func (g GenericDeploy) Deploy(layer DeployLayer) error {
	switch layer {
	case Infrastructure:
		return GenericDeployInfra(g.Caravan, []string{"*"})
	case Platform:
		return GenericDeployPlatform(g.Caravan, []string{"*"})
	case ApplicationSupport:
		return GenericDeployApplicationSupport(g.Caravan, []string{"*"})
	default:
		return fmt.Errorf("unknown Deploy Layer")
	}
}

func GenericDeployInfra(c *Config, targets []string) error {
	// Infra
	fmt.Println("deploying platform")
	t := &tf.Terraform{}
	if err := t.Init(c.WorkdirInfra); err != nil {
		return err
	}
	c.Status = InfraDeployRunning
	if err := c.Save(); err != nil {
		return fmt.Errorf("error persisting state: %w", err)
	}
	env := map[string]string{}
	for _, target := range targets {
		if err := t.ApplyVarFile(filepath.Base(c.WorkdirInfraVars), 600*time.Second, env, target); err != nil {
			return fmt.Errorf("error doing terraform apply: %w", err)
		}
	}

	c.Status = InfraDeployDone
	if err := c.Save(); err != nil {
		return fmt.Errorf("error persisting state: %w", err)
	}

	return nil
}

func GenericDeployPlatform(c *Config, targets []string) error {
	// Platform
	fmt.Printf("deployng platform\n")
	t := tf.Terraform{}
	if err := t.Init(c.WorkdirPlatform); err != nil {
		return err
	}

	c.Status = PlatformDeployRunning
	if err := c.Save(); err != nil {
		return fmt.Errorf("error persisting state: %w", err)
	}

	env := map[string]string{
		"VAULT_TOKEN": c.VaultRootToken,
		"NOMAD_TOKEN": c.NomadToken,
	}
	for _, target := range targets {
		if err := t.ApplyVarFile(filepath.Base(c.WorkdirPlatformVars), 600*time.Second, env, target); err != nil {
			return fmt.Errorf("error doing terraform apply: %w", err)
		}
	}

	c.Status = PlatformDeployDone
	if err := c.Save(); err != nil {
		return fmt.Errorf("error persisting state: %w", err)
	}
	return nil
}

func GenericDeployApplicationSupport(c *Config, targets []string) error {
	// Application support
	t := tf.Terraform{}
	fmt.Printf("deployng application\n")
	if err := t.Init(c.WorkdirApplication); err != nil {
		return err
	}

	c.Status = ApplicationDeployRunning
	if err := c.Save(); err != nil {
		return fmt.Errorf("error persisting state: %w", err)
	}

	env := map[string]string{
		"VAULT_TOKEN": c.VaultRootToken,
		"NOMAD_TOKEN": c.NomadToken,
	}

	for _, target := range targets {
		if err := t.ApplyVarFile(filepath.Base(c.WorkdirApplicationVars), 600*time.Second, env, target); err != nil {
			return fmt.Errorf("error doing terraform apply: %w", err)
		}
	}

	c.Status = ApplicationDeployDone
	if err := c.Save(); err != nil {
		return fmt.Errorf("error persisting state: %w", err)
	}
	return nil
}
