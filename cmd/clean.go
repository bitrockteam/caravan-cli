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
	Use:   "clean --project=<project name> --provider=<provider>",
	Short: "Cleanup the needed config and terraform state store",
	Long: `Deletion of the config files and supporting state stores/locking for terraform created during either up or init. 

The following optional parameters can be specified:

	--region: override of the region as specified in the cloud provider config
	--force: set to true deletes all the objects from the cloud store.`,
	Args: func(cmd *cobra.Command, args []string) error {
		n, _ := cmd.Flags().GetString("project")
		p, _ := cmd.Flags().GetString("provider")
		r := ""

		_, err := caravan.NewConfigFromScratch(n, p, r)
		if err != nil {
			return fmt.Errorf("error generating config: %w", err)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		fmt.Println("clean called")

		name, _ := cmd.Flags().GetString("project")
		project, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")
		force, _ := cmd.Flags().GetBool("force")

		var c *caravan.Config
		c, err = caravan.NewConfigFromFile(name)
		if err != nil {
			// TODO better error handling
			if !strings.Contains(err.Error(), "no such file or directory") {
				return err
			}
			fmt.Printf("getting config from scratch: %s\n", name)
			c, err = caravan.NewConfigFromScratch(name, project, region)
			if err != nil {
				return err
			}
		}
		fmt.Printf("Config: %v\n", c)

		if force {
			c.Force = true
		}

		if c.Destroy {
			tf := terraform.NewTerraform(c.WorkdirInfra)
			err := tf.Destroy(filepath.Base(c.WorkdirInfraVars))
			if err != nil {
				fmt.Printf("error during destroy of cloud resources: %s\n", err)
				if !force {
					return nil
				}
			}
			c.Destroy = false
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
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.PersistentFlags().String("project", "", "name of project")
	_ = cleanCmd.MarkPersistentFlagRequired("project")

	cleanCmd.PersistentFlags().String("provider", "", "name of cloud provider: aws, gcp, az, oci,..")
	_ = cleanCmd.MarkPersistentFlagRequired("provider")

	cleanCmd.PersistentFlags().String("region", "", "provider target region")
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
