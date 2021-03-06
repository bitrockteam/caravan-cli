// Up command.
//
// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"fmt"
	"strings"
	"time"

	"caravan-cli/cli"
	"caravan-cli/cli/checker"

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
			log.Error().Msgf("ERR: %s", err)
			if strings.Contains(err.Error(), "no such file or directory") {
				log.Info().Msgf("please run init")
				return nil
			}
			return err
		}

		prv, err := getProvider(ctx, c)
		if err != nil {
			return err
		}
		log.Info().Msgf("[%s] current status", c.Status)
		if c.Status < cli.InfraDeployDone {
			c.Status = cli.InfraDeployRunning
			if err := c.Save(); err != nil {
				return fmt.Errorf("error persisting state: %w", err)
			}

			err := prv.Deploy(ctx, cli.Infrastructure)
			if err != nil {
				return err
			}

			c.Status = cli.InfraDeployDone
			if err := c.Save(); err != nil {
				return fmt.Errorf("error persisting state: %w", err)
			}
			log.Info().Msgf("[%s] deployment of infrastructure completed", c.Status)

			if err := checkStatus(c, "vault", 30); err != nil {
				return err
			}
			if err := checkStatus(c, "consul", 30); err != nil {
				return err
			}
			if c.DeployNomad {
				if err := checkStatus(c, "nomad", 30); err != nil {
					return err
				}
			}
			log.Debug().Msgf("setting Vault root token")
			if c.VaultRootToken == "" {
				if err := c.SetVaultRootToken(); err != nil {
					return fmt.Errorf("error setting Vault root Token: %w", err)
				}
			}
			if c.DeployNomad {
				if c.NomadToken == "" {
					log.Debug().Msgf("setting Nomad token")
					if err := c.SetNomadToken(); err != nil {
						return fmt.Errorf("error setting Nomad token: %w", err)
					}
				}
			}
			if err := c.Save(); err != nil {
				return fmt.Errorf("error persisting state: %w", err)
			}
		}
		log.Info().Msgf("[%s] deployment of infrastructure completed", c.Status)
		if c.Status < cli.PlatformDeployDone {
			c.Status = cli.PlatformDeployRunning
			if err := c.Save(); err != nil {
				return fmt.Errorf("error persisting state: %w", err)
			}

			err := prv.Deploy(ctx, cli.Platform)
			if err != nil {
				return err
			}

			c.Status = cli.PlatformDeployDone
			if err := c.Save(); err != nil {
				return fmt.Errorf("error persisting state: %w", err)
			}
			log.Info().Msgf("[%s] deployment of platform completed", c.Status)
			if err := checkURL(c, "consul", "/v1/connect/ca/roots", 60); err != nil {
				return err
			}
		}
		log.Info().Msgf("[%s] deployment of platform completed", c.Status)
		if c.DeployNomad {
			if c.Status < cli.ApplicationDeployDone {
				c.Status = cli.ApplicationDeployRunning
				if err := c.Save(); err != nil {
					return fmt.Errorf("error persisting state: %w", err)
				}

				err := prv.Deploy(ctx, cli.ApplicationSupport)
				if err != nil {
					return err
				}

				c.Status = cli.ApplicationDeployDone
				if err := c.Save(); err != nil {
					return fmt.Errorf("error persisting state: %w", err)
				}
			}
			log.Info().Msgf("[%s] deployment of application completed", c.Status)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}

func checkURL(c *cli.Config, tool, path string, count int) (err error) {
	tls, err := checker.TLSClient(c.CAPath)
	if err != nil {
		return err
	}

	checker := checker.NewGenericChecker(fmt.Sprintf("https://%s.%s.%s", tool, c.Name, c.Domain), tls)
	for i := 0; i <= count; i++ {
		if checker.CheckURL(ctx, path) {
			log.Info().Msgf("OK")
			break
		}
		if i >= count {
			log.Warn().Msgf("KO")
			return fmt.Errorf("timeout waiting for %s to be available", tool)
		}
		time.Sleep(6 * time.Second)
		log.Info().Msgf(".")
	}
	return nil
}

func checkStatus(c *cli.Config, tool string, count int) (err error) {
	log.Info().Msgf("checking %s status:", tool)

	var check checker.Checker
	switch tool {
	case cli.Nomad:
		check, err = checker.NewNomadChecker(fmt.Sprintf("https://%s.%s.%s", tool, c.Name, c.Domain), c.CAPath)
		if err != nil {
			return err
		}
	case cli.Consul:
		check, err = checker.NewConsulChecker(fmt.Sprintf("https://%s.%s.%s", tool, c.Name, c.Domain), c.CAPath)
		if err != nil {
			return err
		}
	case cli.Vault:
		check, err = checker.NewVaultChecker(fmt.Sprintf("https://%s.%s.%s", tool, c.Name, c.Domain), c.CAPath)
		if err != nil {
			return err
		}
	default:
		log.Error().Msgf("tool not supported: %s", tool)
	}

	for i := 0; i < count; i++ {
		if check.Status(ctx) {
			log.Info().Msgf("checking %s status: OK", tool)
			break
		}
		if i >= count {
			log.Warn().Msgf("checking %s status: KO", tool)
			return fmt.Errorf("timeout waiting for %s to be available", tool)
		}
		time.Sleep(6 * time.Second)
		fmt.Printf(".")
	}
	return nil
}
