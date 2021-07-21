// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"caravan/internal/aws"
	"caravan/internal/caravan"
	"caravan/internal/gcp"
	"caravan/internal/git"
	"fmt"
	"strings"

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

		name, _ := cmd.Flags().GetString("project")
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")

		c, err := caravan.NewConfigFromFile()
		if err != nil {
			// TODO better error checking
			if !strings.Contains(err.Error(), "no such file or directory") {
				return err
			}
			c, err = caravan.NewConfigFromScratch(name, provider, region)
			if err != nil {
				return err
			}
		}

		if name != c.Name {
			fmt.Printf("please run: \"caravan clean --force\" before init a new project")
			return nil
		}

		b, _ := cmd.Flags().GetString("branch")

		c.SetBranch(b)
		if err := c.Save(); err != nil {
			return err
		}

		// checkout repos
		git := git.NewGit("bitrockteam")
		for _, repo := range c.Repos {
			err := git.Clone(repo, ".caravan/"+c.Name+"/"+repo, b)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return err
			}
		}

		// init Provider
		if err := initCloud(c); err != nil {
			fmt.Printf("error during init: %s\n", err)
			return err
		}
		if c.Status < caravan.InitDone {
			c.Status = caravan.InitDone
			if err := c.Save(); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.PersistentFlags().String("project", "", "name of project")
	_ = initCmd.MarkPersistentFlagRequired("project")

	initCmd.PersistentFlags().String("provider", "", "name of cloud provider: aws, gcp, az, oci,..")
	_ = initCmd.MarkPersistentFlagRequired("provider")

	initCmd.PersistentFlags().String("region", "", "provider target region")
	initCmd.PersistentFlags().String("branch", "", "branch to checkout on repos")
}

func initCloud(c *caravan.Config) (err error) {
	// generate configs and supporting items (bucket and locktable)
	fmt.Printf("initializing cloud resources\n")
	var p caravan.Provider
	switch c.Provider {
	case "aws":
		p, err = aws.New(*c)
		if err != nil {
			return err
		}
	case "gcp":
		p, err = gcp.New(*c)
		if err != nil {
			return err
		}
	default:
		fmt.Printf("impl not found")
		return err
	}

	if err := p.GenerateConfig(); err != nil {
		return fmt.Errorf("error generating config files: %w", err)
	}

	fmt.Printf("creating bucket: %s\n", c.BucketName)
	if err := p.CreateBucket(c.BucketName); err != nil {
		return err
	}

	fmt.Printf("creating lock table: %s\n", c.TableName)
	if err := p.CreateLockTable(c.TableName); err != nil {
		return err
	}

	return nil
}
