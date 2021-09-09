// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"fmt"
	"strings"
	"time"

	"caravan/internal/caravan"

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

		provider, err := getProvider(c)
		if err != nil {
			return err
		}

		if c.Status < caravan.InfraDeployDone {
			err := provider.Deploy(caravan.Infrastructure)
			if err != nil {
				return err
			}
		}
		fmt.Printf("[%s] deployment of infrastructure completed\n", c.Status)

		if err := checkStatus(c, "vault", "/v1/sys/leader", 30); err != nil {
			return err
		}
		if err := checkStatus(c, "consul", "/v1/status/leader", 30); err != nil {
			return err
		}
		if err := checkStatus(c, "nomad", "/v1/status/leader", 30); err != nil {
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
		if err := c.Save(); err != nil {
			return fmt.Errorf("error persisting state: %w", err)
		}
		if c.Status < caravan.PlatformDeployDone {
			err := provider.Deploy(caravan.Platform)
			if err != nil {
				return err
			}
		}
		fmt.Printf("[%s] deployment of platform completed\n", c.Status)
		if err := checkStatus(c, "consul", "/v1/connect/ca/roots", 20); err != nil {
			return err
		}
		if c.Status < caravan.ApplicationDeployDone {
			err := provider.Deploy(caravan.ApplicationSupport)
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
