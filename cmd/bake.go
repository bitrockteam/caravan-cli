// Bake command.
//
// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"caravan-cli/cli"
	"caravan-cli/git"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

// bakeCmd represents the bake command.
var bakeCmd = &cobra.Command{
	Use:   "bake",
	Short: "Generate (bake) up to date VM images for caravan",
	Long:  `Baked images are available for usage in the selected provider's registry provided region.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("project")
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")

		_, err := cli.NewConfigFromScratch(name, provider, region)
		if err != nil {
			return fmt.Errorf("error generating config: %w", err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		log.Debug().Msgf("bake called")

		name, _ := cmd.Flags().GetString("project")
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")
		branch, _ := cmd.Flags().GetString("branch")

		c, err := cli.NewConfigFromScratch(name, provider, region)
		if err != nil {
			return err
		}

		p, err := getProvider(c)
		if err != nil {
			return err
		}

		git := git.NewGit("bitrockteam")
		if err := git.Clone("caravan-baking", filepath.Join(c.WorkdirProject, "caravan-baking"), branch); err != nil {
			return err
		}

		templates, err := p.GetTemplates()
		if err != nil {
			return err
		}
		for _, t := range templates {
			if t.Name == "baking-vars" {
				if err := t.Render(c); err != nil {
					return err
				}
				break
			}
		}
		if err := p.Bake(); err != nil {
			return err
		}
		os.RemoveAll(filepath.Join(c.WorkdirProject, "caravan-baking"))
		return nil

	},
}

func init() {
	rootCmd.AddCommand(bakeCmd)

	bakeCmd.PersistentFlags().String("project", "", "Project name, used for tagging and namespacing")
	bakeCmd.PersistentFlags().String("provider", "", "Cloud provider name. Can be on of aws,gcp, ...")
	bakeCmd.PersistentFlags().String("region", "", "Optional: override default profile region")
	bakeCmd.PersistentFlags().String("branch", "main", "Optional: define a branch to checkout instead of default")

	_ = bakeCmd.MarkPersistentFlagRequired("project")
	_ = bakeCmd.MarkPersistentFlagRequired("provider")
}
