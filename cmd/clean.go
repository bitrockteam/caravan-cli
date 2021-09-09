// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"caravan/internal/caravan"
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

		provider, err := getProvider(c.Provider, c)
		if err != nil {
			return err
		}

		if c.Status >= caravan.ApplicationDeployRunning {
			if err := provider.Destroy(caravan.ApplicationSupport); err != nil {
				return err
			}
		}

		if c.Status >= caravan.PlatformDeployRunning {
			if err := provider.Destroy(caravan.Platform); err != nil {
				return err
			}
		}

		if c.Status >= caravan.InfraDeployRunning {
			if err := provider.Destroy(caravan.Infrastructure); err != nil {
				return err
			}
		}

		err = provider.Clean()
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
