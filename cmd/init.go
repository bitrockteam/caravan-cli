// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"caravan/internal/aws"
	"caravan/internal/caravan"
	"caravan/internal/gcp"
	"caravan/internal/git"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

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

	initCmd.PersistentFlags().String("region", "", "region for the deployment")

	initCmd.PersistentFlags().String("provider", "", "cloud provider")
	_ = initCmd.MarkPersistentFlagRequired("provider")
}

func preRunInit(cmd *cobra.Command, args []string) error {
	provider, _ := cmd.Flags().GetString("provider")
	switch provider {
	case caravan.AWS, caravan.GCP:
		break
	default:
		return fmt.Errorf("unsupported %s provider", provider)
	}

	return nil
}

func executeInit(cmd *cobra.Command, args []string) error {
	provider, _ := cmd.Flags().GetString("provider")
	name, _ := cmd.Flags().GetString("project")
	region, _ := cmd.Flags().GetString("region")
	branch, _ := cmd.Flags().GetString("branch")

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

func getProvider(c *caravan.Config) (caravan.Provider, error) {
	var p caravan.Provider
	var err error
	switch c.Provider {
	case caravan.AWS:
		p, err = aws.New(c)
	case caravan.GCP:
		p, err = gcp.New(c)
	// case caravan.Azure:
	//	p, err = azure.New(c)
	default:
		p, err = nil, fmt.Errorf("unknown provider")
	}
	if err != nil {
		return nil, err
	}
	return p, nil
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

	fmt.Printf("generating terraform config\n")
	if err := GenerateConfig(c, templates); err != nil {
		return fmt.Errorf("error generating config files: %w", err)
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

func GenerateConfig(c *caravan.Config, t []caravan.Template) (err error) {
	// FIXME: rename to process templates
	fmt.Printf("generating config files on: %s\n", c.WorkdirProject)
	if err := os.MkdirAll(c.WorkdirProject, 0777); err != nil {
		return err
	}

	for _, t := range t {
		fmt.Printf("generating %v:%s \n", t.Name, t.Path)
		if err := Generate(t, c); err != nil {
			return err
		}
	}

	return nil
}

func Generate(t caravan.Template, c *caravan.Config) (err error) {
	// FIXME: rename to render template
	temp, err := template.New(t.Name).Parse(t.Text)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(t.Path), 0777); err != nil {
		return err
	}
	f, err := os.Create(t.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := temp.Execute(f, c); err != nil {
		return err
	}
	return nil
}
