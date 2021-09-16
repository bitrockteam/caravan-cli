/*
Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>

*/
package cmd

import (
	"caravan/internal/caravan"
	"caravan/internal/gcp"
	"caravan/internal/git"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// gcpCmd represents the gcp command.
var gcpCmd = &cobra.Command{
	Use:   "gcp",
	Short: "Initialize GCP provider",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gcp called")
		name, _ := cmd.Flags().GetString("project")
		region, _ := cmd.Flags().GetString("region")
		parentProject, _ := cmd.Flags().GetString("parent-project")

		c, err := caravan.NewConfigFromFile()
		if err != nil {
			// TODO better error checking
			if !strings.Contains(err.Error(), "no such file or directory") {
				fmt.Printf("unable to create config from file: %s\n", err)
				return
			}
			c, err = caravan.NewConfigFromScratch(name, caravan.GCP, region)
			if err != nil {
				fmt.Printf("error creating config: %s\n", err)
				return
			}
		}

		if name != c.Name || c.Provider != caravan.GCP {
			fmt.Printf("please run: \"caravan clean --force\" before init a new project")
			return
		}

		b, _ := cmd.Flags().GetString("branch")

		p, err := gcp.New(c)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		c.ParentProject = parentProject
		c.SetBranch(b)
		if err := c.Save(); err != nil {
			fmt.Printf("error saving state: %s\n", err)
			return
		}

		// checkout repos
		git := git.NewGit("bitrockteam")
		for _, repo := range c.Repos {
			err := git.Clone(repo, ".caravan/"+c.Name+"/"+repo, b)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}
		}

		// init Provider
		if err := initProvider(c, p); err != nil {
			fmt.Printf("error during init: %s\n", err)
			return
		}
		if c.Status < caravan.InitDone {
			c.Status = caravan.InitDone
			if err := c.Save(); err != nil {
				fmt.Printf("error saving state: %s\n", err)
			}
		}

	},
}

func init() {
	initCmd.AddCommand(gcpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gcpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gcpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	gcpCmd.PersistentFlags().String("project", "", "GCP project name")
	_ = gcpCmd.MarkPersistentFlagRequired("project")

	gcpCmd.PersistentFlags().String("parent-project", "", "GCP parent project name")
	_ = gcpCmd.MarkPersistentFlagRequired("parent-project")

	gcpCmd.PersistentFlags().String("region", "europe-west6", "GCP deployment region")
	// assume project already created
	/*
		gcpCmd.PersistentFlags().String("orgID", "", "GCP organization ID")
		_ = gcpCmd.MarkPersistentFlagRequired("orgID")

		gcpCmd.PersistentFlags().String("billingAccountID", "", "GCP billing account  ID")
		_ = gcpCmd.MarkPersistentFlagRequired("billingAccountID")
	*/
}
