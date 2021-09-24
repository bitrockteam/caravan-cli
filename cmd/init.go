// Init command.
//
// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"caravan-cli/cli"
	"caravan-cli/git"
	"caravan-cli/provider"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"

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

var (
	// Common.
	prv    = ""
	name   = ""
	region = ""
	branch = ""

	// GCP.
	gcpParentProject = ""
	gcpDnsZone = ""

	// Azure.
	azResourceGroup  = ""
	azSubscriptionID = ""
	azTenantID       = ""
	azUseCLI         = false
)

func init() {
	rootCmd.AddCommand(initCmd)

	// Common
	initCmd.Flags().StringVarP(&name, FlagProject, FlagProjectShort, "", "name of project")
	_ = initCmd.MarkPersistentFlagRequired(FlagProject)
	initCmd.Flags().StringVarP(&prv, FlagProvider, FlagProviderShort, "", "cloud provider")
	_ = initCmd.MarkPersistentFlagRequired(FlagProvider)
	initCmd.Flags().StringVarP(&region, FlagRegion, FlagRegionShort, "", "region for the deployment")
	initCmd.Flags().StringVarP(&branch, FlagBranch, FlagBranchShort, "", "")

	// GCP
	initCmd.Flags().StringVar(&gcpParentProject, FlagGCPParentProject, "", "(GCP only) parent-project")
	initCmd.Flags().StringVar(&gcpDnsZone, FlagGCPDnsZone, "", "(GCP only) cloud dns zone name")

	// Azure
	initCmd.Flags().StringVar(&azResourceGroup, FlagAZResourceGroup, "", "(Azure only) resource group name")
	initCmd.Flags().StringVar(&azSubscriptionID, FlagAZSubscriptionID, "", "(Azure only) subscription ID")
	initCmd.Flags().StringVar(&azTenantID, FlagAZTenantID, "", "(Azure only) tenant ID")
	initCmd.Flags().BoolVar(&azUseCLI, FlagAZLoginViaCLI, false, "(Azure only) login via CLI")
}

func preRunInit(cmd *cobra.Command, args []string) error {
	prv, _ := cmd.Flags().GetString("provider")
	switch prv {
	case "":
		return nil
	case provider.AWS, provider.GCP, provider.Azure:
		break
	default:
		return fmt.Errorf("unsupported provider: %s", prv)
	}
	return nil
}

func executeInit(cmd *cobra.Command, args []string) error {
	c, err := cli.NewConfigFromFile()
	if err != nil {
		// TODO better error checking
		if !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Printf("unable to create config from file: %s\n", err)
			return err
		}
		c, err = cli.NewConfigFromScratch(name, prv, region)
		if err != nil {
			fmt.Printf("unable to create config from scratch: %s\n", err)
			return err
		}
	}

	p, err := getProvider(c)
	if err != nil {
		return err
	}

	if c.Name != name || c.Provider != prv {
		return fmt.Errorf("please run a clean before changing project name or provider")
	}

	if err = processFlags(c); err != nil {
		return err
	}

	if err := initRepos(c, branch); err != nil {
		fmt.Printf("error: %s\n", err)
		return err
	}

	if err := initProvider(c, p); err != nil {
		fmt.Printf("error during init: %s\n", err)
		return err
	}

	if c.Status < cli.InitDone {
		c.Status = cli.InitDone
		if err := c.Save(); err != nil {
			fmt.Printf("error saving state: %s\n", err)
		}
	}

	return nil
}

func initProvider(c *cli.Config, p provider.Provider) error {
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
		if t.Name != "baking-vars" {
			if err := t.Render(c); err != nil {
				return err
			}
		}
	}

	return nil
}

func initRepos(c *cli.Config, b string) (err error) {
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

func processFlags(c *cli.Config) error {
	var err error

	if prv == provider.GCP {
		requiredFlags := map[string]string{
			FlagGCPParentProject: gcpParentProject,
			FlagGCPDnsZone: gcpDnsZone,
		}
		for param, value := range requiredFlags {
			if err2 := mustBeNonEmpty(value, param, provider.GCP); err2 != nil {
				err = multierror.Append(err, err2)
			}
		}

		c.ParentProject = gcpParentProject
	}

	if prv == provider.Azure {
		requiredFlags := map[string]string{
			FlagAZResourceGroup:  azResourceGroup,
			FlagAZSubscriptionID: azSubscriptionID,
			FlagAZTenantID:       azTenantID,
		}
		for param, value := range requiredFlags {
			if err2 := mustBeNonEmpty(value, param, provider.Azure); err2 != nil {
				err = multierror.Append(err, err2)
			}
		}

		c.AzureResourceGroup = azResourceGroup
		c.AzureSubscriptionID = azSubscriptionID
		c.AzureTenantID = azTenantID
		c.AzureUseCLI = azUseCLI
	}

	return err
}

func mustBeNonEmpty(value, paramName, provider string) error {
	if value == "" {
		return fmt.Errorf("--%s flag is needed for %s provider", paramName, provider)
	}
	return nil
}
