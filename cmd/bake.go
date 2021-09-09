// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"caravan/internal/caravan"
	"fmt"

	"github.com/spf13/cobra"
)

// bakeCmd represents the bake command.
var bakeCmd = &cobra.Command{
	Use:   "bake",
	Short: "Generate (bake) up to date VM images for caravan",
	Long:  `Baked images are available for usage in the selected provider's registry provided region.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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
		fmt.Println("bake called")

		name, _ := cmd.Flags().GetString("project")
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")

		c, err := caravan.NewConfigFromScratch(name, provider, region)
		if err != nil {
			return err
		}

		p, err := getProvider(provider, c)
		if err != nil {
			return err
		}

		return p.Bake()
	},
}

func init() {
	rootCmd.AddCommand(bakeCmd)

	bakeCmd.PersistentFlags().String("project", "", "Project name, used for tagging and namespacing")
	bakeCmd.PersistentFlags().String("provider", "", "Cloud provider name. Can be on of aws,gcp, ...")
	bakeCmd.PersistentFlags().String("region", "", "Optional: override default profile region")

	_ = bakeCmd.MarkPersistentFlagRequired("project")
	_ = bakeCmd.MarkPersistentFlagRequired("provider")
}
