// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"caravan/internal/caravan"
	"caravan/internal/git"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// initCmd represents the init command.
var initCmd = &cobra.Command{
	Use:               "init",
	Short:             "Select a provider to initialize",
	Long:              `Initialization of the needed config files and supporting config for a given provider (project, state stores/locking ...)`,
	RunE:              executeInit,
	PersistentPreRunE: preRunInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.PersistentFlags().String("project", "", "name of project")
	_ = initCmd.MarkPersistentFlagRequired("project")

	initCmd.PersistentFlags().String("provider", "", "cloud provider")
	_ = initCmd.MarkPersistentFlagRequired("provider")

	initCmd.PersistentFlags().String("region", "", "region for the deployment")
	initCmd.PersistentFlags().String("parent-project", "", "(GCP only) parent-project")
}

func preRunInit(cmd *cobra.Command, args []string) error {
	provider, _ := cmd.Flags().GetString("provider")
	switch provider {
	case "":
		return nil
	case caravan.AWS, caravan.GCP:
		break
	default:
		return fmt.Errorf("unsupported provider: %s", provider)
	}
	return nil
}

func executeInit(cmd *cobra.Command, args []string) error {
	provider, _ := cmd.Flags().GetString("provider")
	name, _ := cmd.Flags().GetString("project")
	region, _ := cmd.Flags().GetString("region")
	branch, _ := cmd.Flags().GetString("branch")
	parentProject, _ := cmd.Flags().GetString("parent-project")

	c, err := caravan.NewConfigFromFile()
	if err != nil {
		// TODO better error checking
		if !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Printf("unable to create config from file: %s\n", err)
			return err
		}
		c, err = caravan.NewConfigFromScratch(name, provider, region)
		if err != nil {
			fmt.Printf("unable to create config from scratch: %s\n", err)
			return err
		}
	}

	p, err := getProvider(c)
	if err != nil {
		return err
	}

	if c.Name != name || c.Provider != provider {
		return fmt.Errorf("please run a clean before changing project name or provider")
	}

	if provider == caravan.GCP {
		if parentProject == "" {
			return fmt.Errorf("parent-project parameter is needed for GCP provider")
		}
		c.ParentProject = parentProject
	}

	if err := initRepos(c, branch); err != nil {
		fmt.Printf("error: %s\n", err)
		return err
	}

	if err := initProvider(c, p); err != nil {
		fmt.Printf("error during init: %s\n", err)
		return err
	}

	if c.Status < caravan.InitDone {
		c.Status = caravan.InitDone
		if err := c.Save(); err != nil {
			fmt.Printf("error saving state: %s\n", err)
		}
	}

	return nil
}

func initProvider(c *caravan.Config, p caravan.Provider) error {
	fmt.Printf("initializing cloud resources\n")
	if err := p.InitProvider(); err != nil {
		return fmt.Errorf("error initing provider: %w", err)
	}

	templates, err := p.GetTemplates()
	if err != nil {
		return fmt.Errorf("failed to get templates: %w", err)
	}

	fmt.Printf("generating terraform config files on: %s\n", c.WorkdirProject)
	if err := os.MkdirAll(c.WorkdirProject, 0777); err != nil {
		return err
	}
	for _, t := range templates {
		fmt.Printf("generating %v: %s \n", t.Name, t.Path)
		if err := t.Render(c); err != nil {
			return err
		}
	}

	return nil
}

func initRepos(c *caravan.Config, b string) (err error) {
	c.SetBranch(b)
	if err := c.Save(); err != nil {
		return fmt.Errorf("unable to save config after setting branch %s: %w", b, err)
	}
	// checkout repos
	git := git.NewGit("bitrockteam")
	for _, repo := range c.Repos {
		err := git.Clone(repo, ".caravan/"+c.Name+"/"+repo, b)
		if err != nil {
			return fmt.Errorf("unable to clone repo %s: %w", repo, err)
		}
	}
	return nil
}
