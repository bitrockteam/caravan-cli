/*
Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>

*/
package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"caravan/internal/caravan"
	tf "caravan/internal/terraform"

	"github.com/spf13/cobra"
)

// upCmd represents the up command.
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Deploy the caravan infra",
	Long:  `This commands applies the generated terraform configs and provision the needed infrastructure to deploy a caravan instance`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		c, err := caravan.NewConfigFromFile()
		if err != nil {
			fmt.Printf("ERR: %s\n", err)
			if strings.Contains(err.Error(), "no such file or directory") {
				fmt.Printf("please run init")
				return nil
			}
			return err
		}

		if c.Status < caravan.InfraDeployDone {
			err := deployInfra(c)
			if err != nil {
				return err
			}
		}
		fmt.Printf("[%s] deployment of infrastructure completed\n", c.Status)

		if c.VaultRootToken == "" {
			if err := c.SetVaultRootToken(); err != nil {
				return fmt.Errorf("error setting Vault Root Token: %w", err)
			}
		}
		if c.NomadToken == "" {
			if err := c.SetNomadToken(); err != nil {
				return fmt.Errorf("error setting Nomad Token: %w", err)
			}
		}
		if err := c.SaveConfig(); err != nil {
			return fmt.Errorf("error persisting state: %w", err)
		}
		if c.Status < caravan.PlatformDeployDone {
			err := deployPlatform(c)
			if err != nil {
				return err
			}
		}
		fmt.Printf("[%s] deployment of platform completed\n", c.Status)
		fmt.Printf("checking caravan status:")
		co := caravan.NewConsulHealth("https://consul."+c.Name+"."+c.Domain+"/v1/connect/ca/roots", c.CApath)
		for i := 0; i <= 10; i++ {
			if co.Check() {
				break
			}
			if i >= 10 {
				return fmt.Errorf("timeout waiting for consul to be available")
			}
			time.Sleep(6 * time.Second)
			fmt.Printf(".")
		}

		if c.Status < caravan.ApplicationDeployDone {
			err := deployApplication(c)
			if err != nil {
				return err
			}
		}
		fmt.Printf("[%s] deployment of application completed\n", c.Status)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}

func deployInfra(c *caravan.Config) error {
	// Infra
	fmt.Printf("deploying platform")
	t := &tf.Terraform{}
	if err := t.Init(c.WorkdirInfra); err != nil {
		return err
	}
	c.Status = caravan.InfraDeployRunning
	if err := c.SaveConfig(); err != nil {
		return fmt.Errorf("error persisting state: %w", err)
	}
	env := map[string]string{}
	if err := t.ApplyVarFile(filepath.Base(c.WorkdirInfraVars), 600*time.Second, env); err != nil {
		return fmt.Errorf("error doing terraform apply: %w", err)
	}

	c.Status = caravan.InfraDeployDone
	if err := c.SaveConfig(); err != nil {
		return fmt.Errorf("error persisting state: %w", err)
	}

	return nil
}

func deployPlatform(c *caravan.Config) error {
	// Platform
	fmt.Printf("deployng platform\n")
	t := tf.Terraform{}
	if err := t.Init(c.WorkdirPlatform); err != nil {
		return err
	}

	c.Status = caravan.PlatformDeployRunning
	if err := c.SaveConfig(); err != nil {
		return fmt.Errorf("error persisting state: %w", err)
	}

	env := map[string]string{
		"VAULT_TOKEN": c.VaultRootToken,
		"NOMAD_TOKEN": c.NomadToken,
	}
	if err := t.ApplyVarFile(filepath.Base(c.WorkdirPlatformVars), 600*time.Second, env); err != nil {
		return fmt.Errorf("error doing terraform apply: %w", err)
	}

	c.Status = caravan.PlatformDeployDone
	if err := c.SaveConfig(); err != nil {
		return fmt.Errorf("error persisting state: %w", err)
	}
	return nil
}

func deployApplication(c *caravan.Config) error {
	// Application support
	t := tf.Terraform{}
	fmt.Printf("deployng application\n")
	if err := t.Init(c.WorkdirApplication); err != nil {
		return err
	}

	c.Status = caravan.ApplicationDeployRunning
	if err := c.SaveConfig(); err != nil {
		return fmt.Errorf("error persisting state: %w", err)
	}

	env := map[string]string{
		"VAULT_TOKEN": c.VaultRootToken,
		"NOMAD_TOKEN": c.NomadToken,
	}
	if err := t.ApplyVarFile(filepath.Base(c.WorkdirApplicationVars), 600*time.Second, env); err != nil {
		return fmt.Errorf("error doing terraform apply: %w", err)
	}

	c.Status = caravan.PlatformDeployDone
	if err := c.SaveConfig(); err != nil {
		return fmt.Errorf("error persisting state: %w", err)
	}
	return nil
}
