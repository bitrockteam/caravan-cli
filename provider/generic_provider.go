package provider

import (
	"caravan-cli/cli"
	"caravan-cli/terraform"
	"fmt"
	"path/filepath"
	"time"
)

// GenericProvider is the generic implementation of the Provider interface and holds the Caravan config.
type GenericProvider struct {
	Caravan *cli.Config
}

// Bake performs the terraform apply to the caravan-baking repo.
func (g GenericProvider) Bake() error {
	t := terraform.New()
	if err := t.Init(g.Caravan.WorkdirBaking); err != nil {
		return err
	}
	env := map[string]string{}
	if err := t.ApplyVarFile(filepath.Base(g.Caravan.WorkdirBakingVars), 1200*time.Second, env, "*"); err != nil {
		return err
	}
	return nil
}

// Depoly executes the corresponding terraform apply for the given layers/caravan repo.
func (g GenericProvider) Deploy(layer cli.DeployLayer) error {
	switch layer {
	case cli.Infrastructure:
		return GenericDeployInfra(g.Caravan, []string{"*"})
	case cli.Platform:
		return GenericDeployPlatform(g.Caravan, []string{"*"})
	case cli.ApplicationSupport:
		return GenericDeployApplicationSupport(g.Caravan, []string{"*"})
	default:
		return fmt.Errorf("unknown Deploy Layer")
	}
}

func GenericDeployInfra(c *cli.Config, targets []string) error {
	// Infra
	fmt.Println("deploying infra")
	tf := terraform.New()
	if err := tf.Init(c.WorkdirInfra); err != nil {
		return err
	}
	env := map[string]string{}
	for _, target := range targets {
		if err := tf.ApplyVarFile(filepath.Base(c.WorkdirInfraVars), 1200*time.Second, env, target); err != nil {
			return fmt.Errorf("error doing terraform apply: %w", err)
		}
	}
	return nil
}

func GenericDeployPlatform(c *cli.Config, targets []string) error {
	// Platform
	fmt.Printf("deployng platform\n")
	tf := terraform.New()
	if err := tf.Init(c.WorkdirPlatform); err != nil {
		return err
	}
	env := map[string]string{
		"VAULT_TOKEN": c.VaultRootToken,
		"NOMAD_TOKEN": c.NomadToken,
	}
	for _, target := range targets {
		if err := tf.ApplyVarFile(filepath.Base(c.WorkdirPlatformVars), 600*time.Second, env, target); err != nil {
			return fmt.Errorf("error doing terraform apply: %w", err)
		}
	}
	return nil
}

func GenericDeployApplicationSupport(c *cli.Config, targets []string) error {
	// Application support
	tf := terraform.New()
	fmt.Printf("deployng application\n")
	if err := tf.Init(c.WorkdirApplication); err != nil {
		return err
	}
	env := map[string]string{
		"VAULT_TOKEN": c.VaultRootToken,
		"NOMAD_TOKEN": c.NomadToken,
	}
	for _, target := range targets {
		if err := tf.ApplyVarFile(filepath.Base(c.WorkdirApplicationVars), 600*time.Second, env, target); err != nil {
			return fmt.Errorf("error doing terraform apply: %w", err)
		}
	}
	return nil
}

func (g GenericProvider) Destroy(layer cli.DeployLayer) error {
	switch layer {
	case cli.Infrastructure:
		return g.cleanInfra()
	case cli.Platform:
		return g.cleanPlatform()
	case cli.ApplicationSupport:
		return g.cleanApplication()
	default:
		return fmt.Errorf("cannot destroy unknown deploy layer: %d", layer)
	}
}

func (g GenericProvider) cleanInfra() (err error) {
	fmt.Printf("removing terraform infrastructure\n")
	tf := terraform.New()
	err = tf.Init(g.Caravan.WorkdirInfra)
	if err != nil {
		return err
	}
	env := map[string]string{}
	if err := tf.Destroy(filepath.Base(g.Caravan.WorkdirInfraVars), env); err != nil {
		fmt.Printf("error during destroy of cloud resources: %s\n", err)
		if !g.Caravan.Force {
			return err
		}
	}
	return nil
}

func (g GenericProvider) cleanPlatform() (err error) {
	fmt.Printf("removing terraform platform\n")
	tf := terraform.New()
	err = tf.Init(g.Caravan.WorkdirPlatform)
	if err != nil {
		return err
	}
	env := map[string]string{
		"VAULT_TOKEN": g.Caravan.VaultRootToken,
		"NOMAD_TOKEN": g.Caravan.NomadToken,
	}
	if err := tf.Destroy(filepath.Base(g.Caravan.WorkdirPlatformVars), env); err != nil {
		fmt.Printf("error during destroy of cloud resources: %s\n", err)
		if !g.Caravan.Force {
			return err
		}
	}
	return nil
}

func (g GenericProvider) cleanApplication() (err error) {
	fmt.Printf("removing terraform application\n")
	tf := terraform.New()
	err = tf.Init(g.Caravan.WorkdirApplicationVars)
	if err != nil {
		return err
	}
	env := map[string]string{
		"VAULT_TOKEN": g.Caravan.VaultRootToken,
		"NOMAD_TOKEN": g.Caravan.NomadToken,
	}
	if err := tf.Destroy(filepath.Base(g.Caravan.WorkdirApplicationVars), env); err != nil {
		fmt.Printf("error during destroy of cloud resources: %s\n", err)
		if !g.Caravan.Force {
			return err
		}
	}
	return nil
}

func (g GenericProvider) Status() error {
	panic("implement me")
}
