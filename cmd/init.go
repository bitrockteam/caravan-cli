/*
Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>

*/
package cmd

import (
	"caravan/internal/aws"
	"caravan/internal/caravan"
	"caravan/internal/git"
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command.
var initCmd = &cobra.Command{
	Use:   "init --project=<project name> --provider=<provider>",
	Short: "Create the needed config and terraform state store",
	Long: `Initialization of the needed config files and supporting state stores/locking 
	for terraform. The following optional parameters can be specified:
	--region
	--domain
	optional parameters default respectively to the value defined in the default profile and <project>.com.`,
	Args: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("project")
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")

		_, err := caravan.NewConfigFromScratch(name, provider, region)
		if err != nil {
			return fmt.Errorf("error generating config: %w", err)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		fmt.Println("init called")

		name, _ := cmd.Flags().GetString("project")
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")
		c, err := caravan.NewConfigFromScratch(name, provider, region)
		if err != nil {
			return err
		}

		b, _ := cmd.Flags().GetString("branch")

		// checkout repos
		git := git.NewGit("bitrockteam")
		for _, repo := range c.Repos {
			err := git.Clone(repo, ".caravan/"+c.Name+"/"+repo, b)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return err
			}
		}

		// init AWS
		err = initCloud(c)
		if err != nil {
			fmt.Printf("error during init: %s\n", err)
			return err
		}
		c.Status = "INIT_COMPLETE"
		if err := c.SaveConfig(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	initCmd.PersistentFlags().String("project", "", "name of project")
	_ = initCmd.MarkPersistentFlagRequired("project")

	initCmd.PersistentFlags().String("provider", "", "name of cloud provider: aws, gcp, az, oci,..")
	_ = initCmd.MarkPersistentFlagRequired("provider")

	initCmd.PersistentFlags().String("region", "", "provider target region")
	initCmd.PersistentFlags().String("branch", "", "branch to checkout on repos")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initCloud(c *caravan.Config) (err error) {
	// generate configs and supporting items (bucket and locktable)
	fmt.Printf("initializing cloud resources\n")
	cloud, err := aws.NewAWS(*c)
	if err != nil {
		return err
	}

	err = cloud.GenerateConfig()
	if err != nil {
		return fmt.Errorf("error generating config files: %w", err)
	}

	fmt.Printf("creating bucket: %s\n", c.BucketName)
	err = cloud.CreateBucket(c.BucketName)
	if err != nil {
		return err
	}

	fmt.Printf("creating lock table: %s\n", c.TableName)
	err = cloud.CreateLockTable(c.TableName)
	if err != nil {
		return err
	}

	return nil
}
