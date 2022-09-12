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
			return fmt.Errorf("error reading config: %w", err)
		}
		if force {
			c.Force = true
		}
		log.Info().Msgf("running clean on project %s", c.Name)

		if c.Status == cli.InitMissing {
			os.RemoveAll(c.WorkdirProject)
			os.RemoveAll(c.Workdir + "/caravan.state")
			return nil
		}

		prv, err := getProvider(ctx, c)
		if err != nil {
			return fmt.Errorf("error getting provider: %w", err)
		}

		if c.Status > cli.ApplicationCleanDone {
			log.Info().Msgf("[%s] removing application layer", c.Status)
			c.SaveStatus(cli.ApplicationCleanRunning)

			if err := prv.Destroy(ctx, cli.ApplicationSupport); err != nil {
				return err
			}

			c.SaveStatus(cli.ApplicationCleanDone)
			log.Info().Msgf("[%s] application layer removed", c.Status)
		}

		if c.Status > cli.PlatformCleanDone {
			log.Info().Msgf("[%s] removing platform layer", c.Status)
			c.SaveStatus(cli.PlatformCleanRunning)

			if err := prv.Destroy(ctx, cli.Platform); err != nil {
				return err
			}

			c.SaveStatus(cli.PlatformCleanDone)
			log.Info().Msgf("[%s] platform layer removed", c.Status)
		}

		if c.Status > cli.InfraCleanDone {
			log.Info().Msgf("[%s] removing infra layer", c.Status)
			c.SaveStatus(cli.InfraCleanRunning)

			if err := prv.Destroy(ctx, cli.Infrastructure); err != nil {
				return err
			}

			c.SaveStatus(cli.InfraCleanDone)
			log.Info().Msgf("[%s] infra layer removed", c.Status)
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
