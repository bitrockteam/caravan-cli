/*
Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>

*/
package cmd

import (
	"caravan/internal/caravan"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Reports the current caravan status",
	Long: `Gets and diplay the current status for caravan both locally and remote:

	--project: project name to get the status for
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("project")
		c, err := caravan.NewConfigFromFile(name)
		if err != nil {
			if strings.Contains(err.Error(), "no such file or directory") {
				fmt.Printf("project %s status %s\n", name, caravan.InitMissing)
				return nil
			}
			return err
		}
		c.StatusReport()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.PersistentFlags().String("project", "", "name of project")
	_ = statusCmd.MarkPersistentFlagRequired("project")
}
