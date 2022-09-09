// Status command.
//
// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"caravan-cli/cli"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Reports the current caravan status",
	Long: `Gets and diplay the current status for caravan both locally and remotely
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cli.NewConfigFromFile()
		if err != nil {
			if strings.Contains(err.Error(), "no such file or directory") {
				fmt.Printf("project status is missing: %s\n", cli.InitMissing)
				return nil
			}
			return err
		}

		log.Info().Msgf("[%s] running status on project %s", c.Status, c.Name)

		r := cli.NewReport(c)
		if c.Status >= cli.InfraDeployDone {
			if err := r.CheckStatus(ctx); err != nil {
				return err
			}
		}
		r.PrintReport()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
