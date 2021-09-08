// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"caravan/internal/caravan"
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command.
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Select a provider to initialize",
	Long:  `Initialization of the needed config files and supporting config for a given provider (project, state stores/locking ...)`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initProvider(c *caravan.Config, p caravan.Provider) (err error) {
	fmt.Printf("initializing cloud resources\n")
	if err := p.Init(); err != nil {
		return fmt.Errorf("error initing provider: %w", err)
	}
	fmt.Printf("generating terraform config\n")
	if err := p.GenerateConfig(); err != nil {
		return fmt.Errorf("error generating config files: %w", err)
	}

	fmt.Printf("creating statestore: %s\n", c.StateStoreName)
	if err := p.CreateStateStore(c.StateStoreName); err != nil {
		return err
	}

	fmt.Printf("creating lock: %s\n", c.LockName)
	if err := p.CreateLock(c.LockName); err != nil {
		return err
	}

	return nil
}
