// Up command.
//
// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"fmt"
	"strings"
	"time"

	"caravan-cli/cli"
	"caravan-cli/health"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

// upCmd represents the up command.
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Deploy the caravan infra",
	Long:  `This commands applies the generated terraform configs and provision the needed infrastructure to deploy a caravan instance`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		c, err := cli.NewConfigFromFile()
		if err != nil {
			log.Error().Msgf("ERR: %s\n", err)
			if strings.Contains(err.Error(), "no such file or directory") {
				log.Info().Msgf("please run init")
				return nil
			}
			return err
		}

		prv, err := getProvider(c)
		if err != nil {
			return err
		}
		if c.Status < cli.InfraDeployDone {
			c.Status = cli.InfraDeployRunning
			if err := c.Save(); err != nil {
				return fmt.Errorf("error persisting state: %w", err)
			}

			err := prv.Deploy(cli.Infrastructure)
			if err != nil {
				return err
			}

			c.Status = cli.InfraDeployDone
			if err := c.Save(); err != nil {
				return fmt.Errorf("error persisting state: %w", err)
			}
		}
		log.Info().Msgf("[%s] deployment of infrastructure completed\n", c.Status)

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
		if c.Status < cli.PlatformDeployDone {
			c.Status = cli.PlatformDeployRunning
			if err := c.Save(); err != nil {
				return fmt.Errorf("error persisting state: %w", err)
			}

			err := prv.Deploy(cli.Platform)
			if err != nil {
				return err
			}

			c.Status = cli.PlatformDeployDone
			if err := c.Save(); err != nil {
				return fmt.Errorf("error persisting state: %w", err)
			}
		}
		log.Info().Msgf("[%s] deployment of platform completed\n", c.Status)
		if err := checkStatus(c, "consul", "/v1/connect/ca/roots", 20); err != nil {
			return err
		}
		if c.Status < cli.ApplicationDeployDone {
			c.Status = cli.ApplicationDeployRunning
			if err := c.Save(); err != nil {
				return fmt.Errorf("error persisting state: %w", err)
			}

			err := prv.Deploy(cli.ApplicationSupport)
			if err != nil {
				return err
			}

			c.Status = cli.ApplicationDeployDone
			if err := c.Save(); err != nil {
				return fmt.Errorf("error persisting state: %w", err)
			}
		}
		log.Info().Msgf("[%s] deployment of application completed\n", c.Status)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}

func checkStatus(c *cli.Config, tool string, path string, count int) error {
	log.Info().Msgf("checking %s status:", tool)

	h := health.NewHealth("https://"+tool+"."+c.Name+"."+c.Domain+path, c.CAPath)
	for i := 0; i <= count; i++ {
		if h.Check() {
			log.Info().Msgf("OK\n")
			break
		}
		if i >= count {
			log.Warn().Msgf("KO\n")
			return fmt.Errorf("timeout waiting for %s to be available", tool)
		}
		time.Sleep(6 * time.Second)
		log.Info().Msgf(".")
	}
	return nil
}
