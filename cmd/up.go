/*
Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>

*/
package cmd

import (
	"fmt"
	"strings"

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
		name, _ := cmd.Flags().GetString("project")

		c, err := caravan.NewConfigFromFile(name)
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

		if c.Status > caravan.InfraDeployDone {
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	upCmd.PersistentFlags().String("project", "", "name of project")
	_ = upCmd.MarkPersistentFlagRequired("project")
	upCmd.PersistentFlags().String("provider", "", "name of provider")
	_ = upCmd.MarkPersistentFlagRequired("provider")
}

func deployInfra(c *caravan.Config) error {
	tf := tf.NewTerraform(c.WorkdirInfra)
	err := tf.Init()
	if err != nil {
		return err
	}
	c.Status = caravan.InfraDeployRunning
	if err := c.SaveConfig(); err != nil {
		return fmt.Errorf("error persisting state: %w", err)
	}

	if err := tf.ApplyVarFile(c.Name + "-infra.tfvars"); err != nil {
		return fmt.Errorf("error doing terraform apply: %w", err)
	}

	// TODO remove and use status as marker for destruction
	c.Destroy = true
	if err := c.SaveConfig(); err != nil {
		return fmt.Errorf("error persisting state: %w", err)
	}

	c.Status = caravan.InfraDeployDone
	if err := c.SaveConfig(); err != nil {
		return fmt.Errorf("error persisting state: %w", err)
	}
	return nil
}

func deployPlatform(c *caravan.Config) error {
	tf := tf.NewTerraform(c.WorkdirPlatform)
	if err := tf.Init(); err != nil {
		return err
	}
	return nil
}
