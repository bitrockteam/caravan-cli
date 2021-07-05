/*
Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>

*/
package cmd

import (
	"caravan/internal/aws"
	"caravan/internal/caravan"
	"caravan/internal/terraform"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command.
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleanup the needed config and terraform state store",
	Long: `Deletion of the config files and supporting state stores/locking for terraform created during either up or init. 

The following optional parameters can be specified:

	--force: set to true deletes all the objects from the cloud store.`,
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

		if c.Status >= caravan.PlatformDeployRunning {
			tf := terraform.Terraform{}
			if tf.Init(c.WorkdirPlatform); err != nil {
				return err
			}
			env := map[string]string{
				"VAULT_TOKEN": c.VaultRootToken,
				"NOMAD_TOKEN": c.NomadToken,
			}
			if err := tf.Destroy(filepath.Base(c.WorkdirPlatformVars), env); err != nil {
				fmt.Printf("error during destroy of cloud resources: %s\n", err)
				if !force {
					return nil
				}
			}
			c.Status = caravan.InfraDeployDone
			if err := c.SaveConfig(); err != nil {
				fmt.Printf("error during config update of config: %s\n", err)
				return nil
			}
		}

		if c.Status >= caravan.InfraDeployRunning {
			tf := terraform.Terraform{}
			if tf.Init(c.WorkdirInfra); err != nil {
				return err
			}
			env := map[string]string{}
			if err := tf.Destroy(filepath.Base(c.WorkdirInfraVars), env); err != nil {
				fmt.Printf("error during destroy of cloud resources: %s\n", err)
				if !force {
					return nil
				}
			}
			c.Status = caravan.InitDone
			if err := c.SaveConfig(); err != nil {
				fmt.Printf("error during config update of config: %s\n", err)
				return nil
			}
		}

		err = cleanCloud(c)
		if err != nil {
			fmt.Printf("error during clean of cloud resources: %s\n", err)
			return nil
		}
		os.RemoveAll(c.Workdir + "/" + c.Name)
		os.RemoveAll(c.Workdir + "/caravan.state")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.PersistentFlags().Bool("force", false, "force cleanup of S3 bucket")
}

func cleanCloud(cfg *caravan.Config) (err error) {
	// generate configs and supporting items (bucket and locktable)
	fmt.Printf("removing terraform state and locking structures\n")

	cloud, err := aws.NewAWS(*cfg)
	if err != nil {
		return err
	}

	if cfg.Force {
		fmt.Printf("emptying bucket %s\n", cfg.Name+"-caravan-terraform-state")
		err = cloud.EmptyBucket(cfg.Name + "-caravan-terraform-state")
		if err != nil {
			return fmt.Errorf("error emptying: %w", err)
		}
	}

	// TODO cleanup before delete with force option
	err = cloud.DeleteBucket(cfg.Name + "-caravan-terraform-state")
	if err != nil {
		return err
	}

	err = cloud.DeleteLockTable(cfg.Name + "-caravan-terraform-state-lock")
	if err != nil {
		return err
	}

	return nil
}
