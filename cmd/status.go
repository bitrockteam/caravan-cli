// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"caravan-cli/cli"
	caravan "caravan-cli/config"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Reports the current caravan status",
	Long: `Gets and diplay the current status for caravan both locally and remotely
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := caravan.NewConfigFromFile()
		if err != nil {
			if strings.Contains(err.Error(), "no such file or directory") {
				fmt.Printf("project status is missing: %s\n", caravan.InitMissing)
				return nil
			}
			return err
		}
		r := cli.Report{Caravan: c}
		r.StatusReport()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
