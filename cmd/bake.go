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
		_, err := cli.NewConfigFromScratch(name, prv, region)
		if err != nil {
			return fmt.Errorf("error generating config: %w", err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		log.Debug().Msgf("bake called")

		c, err := cli.NewConfigFromScratch(name, prv, region)
		if err != nil {
			return err
		}
		c.LogLevel = logLevel
		p, err := getProvider(ctx, c)
		if err != nil {
			return err
		}

		if err := c.SetDistro(distro); err != nil {
			return err
		}

		git := git.NewGit("bitrockteam", logLevel)
		if err := git.Clone("caravan-baking", filepath.Join(c.WorkdirProject, "caravan-baking"), branch); err != nil {
			return err
		}

		templates, err := p.GetTemplates(ctx)
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
		if err := p.Bake(ctx); err != nil {
			return err
		}
		os.RemoveAll(filepath.Join(c.WorkdirProject, "caravan-baking"))
		return nil

	},
}

func init() {
	rootCmd.AddCommand(bakeCmd)

	bakeCmd.Flags().StringVarP(&name, FlagProject, FlagProjectShort, "", "name of project")
	bakeCmd.Flags().StringVarP(&prv, FlagProvider, FlagProviderShort, "", "cloud provider")
	bakeCmd.Flags().StringVarP(&distro, FlagLinuxDistro, FlagLinuxDistroShort, "centos7", "linux distribution")
	bakeCmd.Flags().StringVarP(&region, FlagRegion, FlagRegionShort, "", "optional: override default profile region")
	bakeCmd.Flags().StringVarP(&branch, FlagBranch, FlagBranchShort, "main", "optional: define a branch to checkout instead of default")

	_ = bakeCmd.MarkFlagRequired(FlagProject)
	_ = bakeCmd.MarkFlagRequired(FlagProvider)
}
