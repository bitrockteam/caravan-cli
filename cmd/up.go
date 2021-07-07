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

		if err := checkStatus(c, "vault", "/v1/sys/leader", 20); err != nil {
			return err
		}
		if err := checkStatus(c, "consul", "/v1/status/leader", 20); err != nil {
			return err
		}
		if err := checkStatus(c, "nomad", "/v1/status/leader", 20); err != nil {
			return err
		}
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
		if err := checkStatus(c, "consul", "/v1/connect/ca/roots", 30); err != nil {
			return err
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

func checkStatus(c *caravan.Config, tool string, path string, count int) error {
	fmt.Printf("checking %s status:", tool)

	h := caravan.NewHealth("https://"+tool+"."+c.Name+"."+c.Domain+path, c.CApath)
	for i := 0; i <= count; i++ {
		if h.Check() {
			fmt.Printf("OK\n")
			break
		}
		if i >= count {
			fmt.Printf("KO\n")
			return fmt.Errorf("timeout waiting for %s to be available", tool)
		}
		time.Sleep(6 * time.Second)
		fmt.Printf(".")
	}
	return nil
}
