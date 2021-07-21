// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
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
		if c.Status >= caravan.ApplicationDeployRunning {
			if err := cleanApplication(c); err != nil {
				return err
			}
		}

		if c.Status >= caravan.PlatformDeployRunning {
			if err := cleanPlatform(c); err != nil {
				return err
			}
		}

		if c.Status >= caravan.InfraDeployRunning {
			if err := cleanInfra(c); err != nil {
				return err
			}
		}

		err = cleanCloud(c)
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

func cleanCloud(c *caravan.Config) (err error) {
	fmt.Printf("removing terraform state and locking structures\n")

	cloud, err := aws.New(*c)
	if err != nil {
		return err
	}

	if c.Force {
		fmt.Printf("emptying bucket %s\n", c.Name+"-caravan-terraform-state")
		err = cloud.EmptyBucket(c.Name + "-caravan-terraform-state")
		if err != nil {
			return fmt.Errorf("error emptying: %w", err)
		}
	}

	// TODO cleanup before delete with force option
	err = cloud.DeleteBucket(c.Name + "-caravan-terraform-state")
	if err != nil {
		return err
	}

	err = cloud.DeleteLockTable(c.Name + "-caravan-terraform-state-lock")
	if err != nil {
		return err
	}

	return nil
}

func cleanInfra(c *caravan.Config) (err error) {
	fmt.Printf("removing terraform infrastructure\n")
	tf := terraform.Terraform{}
	err = tf.Init(c.WorkdirInfra)
	if err != nil {
		return err
	}
	c.Status = caravan.InfraCleanRunning
	if err := c.Save(); err != nil {
		fmt.Printf("error during config update of config: %s\n", err)
		return nil
	}
	env := map[string]string{}
	if err := tf.Destroy(filepath.Base(c.WorkdirInfraVars), env); err != nil {
		fmt.Printf("error during destroy of cloud resources: %s\n", err)
		if !c.Force {
			return nil
		}
	}
	c.Status = caravan.InfraCleanDone
	if err := c.Save(); err != nil {
		fmt.Printf("error during config update of config: %s\n", err)
		return nil
	}
	return nil
}

func cleanPlatform(c *caravan.Config) (err error) {
	fmt.Printf("removing terraform platform\n")
	tf := terraform.Terraform{}
	err = tf.Init(c.WorkdirPlatform)
	if err != nil {
		return err
	}
	env := map[string]string{
		"VAULT_TOKEN": c.VaultRootToken,
		"NOMAD_TOKEN": c.NomadToken,
	}
	c.Status = caravan.PlatformCleanRunning
	if err := c.Save(); err != nil {
		fmt.Printf("error during config update of config: %s\n", err)
		return nil
	}
	if err := tf.Destroy(filepath.Base(c.WorkdirPlatformVars), env); err != nil {
		fmt.Printf("error during destroy of cloud resources: %s\n", err)
		if !c.Force {
			return nil
		}
	}
	c.Status = caravan.PlatformCleanDone
	if err := c.Save(); err != nil {
		fmt.Printf("error during config update of config: %s\n", err)
		return nil
	}
	return nil
}

func cleanApplication(c *caravan.Config) (err error) {
	fmt.Printf("removing terraform application\n")
	tf := terraform.Terraform{}
	err = tf.Init(c.WorkdirApplicationVars)
	if err != nil {
		return err
	}
	c.Status = caravan.ApplicationCleanRunning
	if err := c.Save(); err != nil {
		fmt.Printf("error during config update of config: %s\n", err)
		return nil
	}
	env := map[string]string{
		"VAULT_TOKEN": c.VaultRootToken,
		"NOMAD_TOKEN": c.NomadToken,
	}
	if err := tf.Destroy(filepath.Base(c.WorkdirApplicationVars), env); err != nil {
		fmt.Printf("error during destroy of cloud resources: %s\n", err)
		if !c.Force {
			return nil
		}
	}
	c.Status = caravan.ApplicationCleanDone
	if err := c.Save(); err != nil {
		fmt.Printf("error during config update of config: %s\n", err)
		return nil
	}
	return nil
}
