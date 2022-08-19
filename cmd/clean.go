// Clean command.
//
// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"caravan-cli/cli"
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command.
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleanup the needed config and terraform state store",
	Long:  `Deletion of the config files and supporting state stores/locking for terraform created during either up or init.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var c *cli.Config
		c, err = cli.NewConfigFromFile()
		if err != nil {
			if errors.As(err, &cli.ConfigFileNotFound{}) {
				log.Info().Msgf("all clean")
				return nil
			}
			return fmt.Errorf("error reading config: %s", err)
		}
		if force {
			c.Force = true
		}

		if c.Status == cli.InitMissing {
			os.RemoveAll(c.WorkdirProject)
			os.RemoveAll(c.Workdir + "/caravan.state")
			return nil
		}

		prv, err := getProvider(ctx, c)
		if err != nil {
			return fmt.Errorf("error getting provider: %s", err)
		}

		if c.Status > cli.ApplicationCleanDone {
			c.Status = cli.ApplicationCleanRunning
			if err := c.Save(); err != nil {
				log.Error().Msgf("error during config update of config: %s", err)
				return nil
			}

			if err := prv.Destroy(ctx, cli.ApplicationSupport); err != nil {
				return err
			}

			c.Status = cli.ApplicationCleanDone
			if err := c.Save(); err != nil {
				log.Error().Msgf("error during config update of config: %s", err)
				return nil
			}
		}

		if c.Status > cli.PlatformCleanDone {
			c.Status = cli.PlatformCleanRunning
			if err := c.Save(); err != nil {
				log.Error().Msgf("error during config update of config: %s", err)
				return nil
			}

			if err := prv.Destroy(ctx, cli.Platform); err != nil {
				return err
			}

			c.Status = cli.PlatformCleanDone
			if err := c.Save(); err != nil {
				log.Error().Msgf("error during config update of config: %s", err)
				return nil
			}
		}

		if c.Status > cli.InfraCleanDone {
			c.Status = cli.InfraCleanRunning
			if err := c.Save(); err != nil {
				log.Error().Msgf("error during config update of config: %s", err)
				return nil
			}

			if err := prv.Destroy(ctx, cli.Infrastructure); err != nil {
				return err
			}

			c.Status = cli.InfraCleanDone
			if err := c.Save(); err != nil {
				log.Error().Msgf("error during config update of config: %s", err)
				return nil
			}
		}

		err = prv.CleanProvider(ctx)
		if err != nil {
			log.Error().Msgf("error during clean of cloud resources: %s", err)
			return nil
		}
		log.Info().Msgf("removing %s/%s", c.Workdir, c.Name)

		os.RemoveAll(c.WorkdirProject)
		os.RemoveAll(c.Workdir + "/caravan.state")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.PersistentFlags().BoolVarP(&force, FlagForce, FlagForceShort, false, "force cleanup of S3 bucket")
}
