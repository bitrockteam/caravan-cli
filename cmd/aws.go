/*
Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>

*/
package cmd

import (
	"caravan/internal/aws"
	"caravan/internal/caravan"
	"caravan/internal/git"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// awsCmd represents the aws command.
var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "Initializes aws provider",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("aws called")
		name, _ := cmd.Flags().GetString("project")
		region, _ := cmd.Flags().GetString("region")

		c, err := caravan.NewConfigFromFile()
		if err != nil {
			// TODO better error checking
			if !strings.Contains(err.Error(), "no such file or directory") {
				fmt.Printf("unable to create config from file: %s\n", err)
				return
			}
			c, err = caravan.NewConfigFromScratch(name, caravan.AWS, region)
			if err != nil {
				fmt.Printf("unable to create config from scratch: %s\n", err)
				return
			}
		}

		if name != c.Name || c.Provider != caravan.AWS {
			fmt.Printf("please run: \"caravan clean --force\" before init a new project")
			return
		}

		b, _ := cmd.Flags().GetString("branch")

		p, err := aws.New(*c)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}

		c.SetBranch(b)
		if err := c.Save(); err != nil {
			fmt.Printf("unable to set branch: %s\n", err)
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
	initCmd.AddCommand(awsCmd)

	// awsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// awsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	awsCmd.PersistentFlags().String("project", "", "name of project")
	_ = awsCmd.MarkPersistentFlagRequired("project")

	awsCmd.PersistentFlags().String("region", "", "region for the deployment")
}
