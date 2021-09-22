// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"caravan-cli/cli"
	caravan "caravan-cli/config"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command.
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleanup the needed config and terraform state store",
	Long:  `Deletion of the config files and supporting state stores/locking for terraform created during either up or init.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		force, _ := cmd.Flags().GetBool("force")

		var c *caravan.Config
		c, err = caravan.NewConfigFromFile()
		if err != nil {
			// TODO better error handling
			if !strings.Contains(err.Error(), "no such file or directory") {
				return err
			}
			fmt.Printf("all clean\n")
			return nil
		}
		if force {
			c.Force = true
		}

		prv, err := getProvider(c)
		if err != nil {
			return err
		}

		if c.Status >= caravan.ApplicationDeployRunning {
			c.Status = caravan.ApplicationCleanRunning
			if err := c.Save(); err != nil {
				fmt.Printf("error during config update of config: %s\n", err)
				return nil
			}

			if err := prv.Destroy(cli.ApplicationSupport); err != nil {
				return err
			}

			c.Status = caravan.ApplicationCleanDone
			if err := c.Save(); err != nil {
				fmt.Printf("error during config update of config: %s\n", err)
				return nil
			}
		}

		if c.Status >= caravan.PlatformDeployRunning {
			c.Status = caravan.PlatformCleanRunning
			if err := c.Save(); err != nil {
				fmt.Printf("error during config update of config: %s\n", err)
				return nil
			}

			if err := prv.Destroy(cli.Platform); err != nil {
				return err
			}

			c.Status = caravan.PlatformCleanDone
			if err := c.Save(); err != nil {
				fmt.Printf("error during config update of config: %s\n", err)
				return nil
			}
		}

		if c.Status >= caravan.InfraDeployRunning {
			c.Status = caravan.InfraCleanRunning
			if err := c.Save(); err != nil {
				fmt.Printf("error during config update of config: %s\n", err)
				return nil
			}

			if err := prv.Destroy(cli.Infrastructure); err != nil {
				return err
			}

			c.Status = caravan.InfraCleanDone
			if err := c.Save(); err != nil {
				fmt.Printf("error during config update of config: %s\n", err)
				return nil
			}
		}

		err = prv.CleanProvider()
		if err != nil {
			fmt.Printf("error during clean of cloud resources: %s\n", err)
			return nil
		}
		fmt.Printf("removing %s/%s\n", c.Workdir, c.Name)

		os.RemoveAll(c.Workdir + "/" + c.Name)
		os.RemoveAll(c.Workdir + "/caravan.state")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.PersistentFlags().Bool("force", false, "force cleanup of S3 bucket")
}
