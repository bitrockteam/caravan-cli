package caravan

import (
	"caravan/internal/terraform"
	"fmt"
	"path/filepath"
)

type GenericDestroy struct {
	GenericProvider
}

func (g GenericDestroy) Destroy(layer DeployLayer) error {
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

func (g GenericDestroy) cleanInfra() (err error) {
	fmt.Printf("removing terraform infrastructure\n")
	tf := terraform.Terraform{}
	err = tf.Init(g.Caravan.WorkdirInfra)
	if err != nil {
		return err
	}
	g.Caravan.Status = InfraCleanRunning
	if err := g.Caravan.Save(); err != nil {
		fmt.Printf("error during config update of config: %s\n", err)
		return nil
	}
	env := map[string]string{}
	if err := tf.Destroy(filepath.Base(g.Caravan.WorkdirInfraVars), env); err != nil {
		fmt.Printf("error during destroy of cloud resources: %s\n", err)
		if !g.Caravan.Force {
			return nil
		}
	}
	g.Caravan.Status = InfraCleanDone
	if err := g.Caravan.Save(); err != nil {
		fmt.Printf("error during config update of config: %s\n", err)
		return nil
	}
	return nil
}

func (g GenericDestroy) cleanPlatform() (err error) {
	fmt.Printf("removing terraform platform\n")
	tf := terraform.Terraform{}
	err = tf.Init(g.Caravan.WorkdirPlatform)
	if err != nil {
		return err
	}
	env := map[string]string{
		"VAULT_TOKEN": g.Caravan.VaultRootToken,
		"NOMAD_TOKEN": g.Caravan.NomadToken,
	}
	g.Caravan.Status = PlatformCleanRunning
	if err := g.Caravan.Save(); err != nil {
		fmt.Printf("error during config update of config: %s\n", err)
		return nil
	}
	if err := tf.Destroy(filepath.Base(g.Caravan.WorkdirPlatformVars), env); err != nil {
		fmt.Printf("error during destroy of cloud resources: %s\n", err)
		if !g.Caravan.Force {
			return nil
		}
	}
	g.Caravan.Status = PlatformCleanDone
	if err := g.Caravan.Save(); err != nil {
		fmt.Printf("error during config update of config: %s\n", err)
		return nil
	}
	return nil
}

func (g GenericDestroy) cleanApplication() (err error) {
	fmt.Printf("removing terraform application\n")
	tf := terraform.Terraform{}
	err = tf.Init(g.Caravan.WorkdirApplicationVars)
	if err != nil {
		return err
	}
	g.Caravan.Status = ApplicationCleanRunning
	if err := g.Caravan.Save(); err != nil {
		fmt.Printf("error during config update of config: %s\n", err)
		return nil
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
	g.Caravan.Status = ApplicationCleanDone
	if err := g.Caravan.Save(); err != nil {
		fmt.Printf("error during config update of config: %s\n", err)
		return nil
	}
	return nil
}
