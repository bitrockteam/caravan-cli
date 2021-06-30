/*
Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>

*/
package cmd

import (
	"fmt"
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

		// TODO is this really needed ? is tf apply idempotent ?
		if c.Status < caravan.InfraDeployDone {
			err := deployInfra(c)
			if err != nil {
				return err
			}
		}

		if c.Status >= caravan.InfraDeployDone {
			err := deployPlatform(c)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}

func deployInfra(c *caravan.Config) error {
	// Infra
	t := &tf.Terraform{}
	if err := t.Init(c.WorkdirInfra); err != nil {
		return err
	}
	c.Status = caravan.InfraDeployRunning
	if err := c.SaveConfig(); err != nil {
		return fmt.Errorf("error persisting state: %w", err)
	}

	if err := t.ApplyVarFile(c.Name+"-infra.tfvars", 600*time.Second); err != nil {
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
	t := tf.Terraform{}
	if err := t.Init(c.WorkdirPlatform); err != nil {
		return err
	}

	c.Status = caravan.PlatformDeployRunning
	if err := c.SaveConfig(); err != nil {
		return fmt.Errorf("error persisting state: %w", err)
	}
	if err := t.ApplyVarFile(c.Name+"-"+c.Provider+".tfvars", 600*time.Second); err != nil {
		return fmt.Errorf("error doing terraform apply: %w", err)
	}

	c.Status = caravan.PlatformDeployDone
	if err := c.SaveConfig(); err != nil {
		return fmt.Errorf("error persisting state: %w", err)
	}
	return nil
}
