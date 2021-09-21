package caravan

import (
	"caravan/internal/terraform"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type GenericProvider struct {
	Caravan *Config
}

func (g GenericProvider) Bake() error {
	if _, err := os.Stat(g.Caravan.WorkdirProject); os.IsNotExist(err) {
		return fmt.Errorf("please run init before bake")
	}

	t := terraform.New()
	if err := t.Init(g.Caravan.WorkdirBaking); err != nil {
		return err
	}
	env := map[string]string{}
	if err := t.ApplyVarFile(g.Caravan.WorkdirBakingVars, 1200*time.Second, env, "*"); err != nil {
		return err
	}
	return nil
}

func (g GenericProvider) Deploy(layer DeployLayer) error {
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
	tf := terraform.New()
	if err := tf.Init(c.WorkdirInfra); err != nil {
		return err
	}
	env := map[string]string{}
	for _, target := range targets {
		if err := tf.ApplyVarFile(filepath.Base(c.WorkdirInfraVars), 600*time.Second, env, target); err != nil {
			return fmt.Errorf("error doing terraform apply: %w", err)
		}
	}
	return nil
}

func GenericDeployPlatform(c *Config, targets []string) error {
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

func GenericDeployApplicationSupport(c *Config, targets []string) error {
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

func (g GenericProvider) Destroy(layer DeployLayer) error {
	switch layer {
	case Infrastructure:
		return g.cleanInfra()
	case Platform:
		return g.cleanPlatform()
	case ApplicationSupport:
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
			return nil
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
			return nil
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
			return nil
		}
	}
	return nil
}

func (g GenericProvider) Status() error {
	panic("implement me")
}
