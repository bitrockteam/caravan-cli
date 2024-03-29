// Init command.
//
// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"caravan-cli/cli"
	"caravan-cli/git"
	"caravan-cli/provider"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// initCmd represents the init command.
var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "Select a provider to initialize",
	Long:    `Initialization of the needed config files and supporting config for a given provider (project, state stores/locking ...)`,
	RunE:    executeInit,
	PreRunE: preRunInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Common
	initCmd.Flags().StringVarP(&name, FlagProject, FlagProjectShort, "", "name of project")
	initCmd.Flags().StringVarP(&prv, FlagProvider, FlagProviderShort, "", "cloud provider")
	initCmd.Flags().StringVarP(&domain, FlagDomain, FlagDomainShort, "", "")
	initCmd.Flags().StringVarP(&distro, FlagLinuxDistro, FlagLinuxDistroShort, "ubuntu-2204", "linux distribution for image")
	initCmd.Flags().StringVarP(&edition, FlagEdition, FlagEditionShort, "os", "Hashicorp tools edition (os: open source/ent: enterprise")

	_ = initCmd.MarkFlagRequired(FlagProject)
	_ = initCmd.MarkFlagRequired(FlagProvider)
	_ = initCmd.MarkFlagRequired(FlagDomain)

	initCmd.Flags().StringVarP(&region, FlagRegion, FlagRegionShort, "", "region for the deployment")
	initCmd.Flags().StringVarP(&branch, FlagBranch, FlagBranchShort, "", "")
	initCmd.Flags().BoolVar(&deployNomad, FlagDeployNomad, true, "deploy Nomad")

	// GCP
	initCmd.Flags().StringVar(&gcpParentProject, FlagGCPParentProject, "", "(GCP only) parent-project")
	initCmd.Flags().StringVar(&gcpDNSZone, FlagGCPDnsZone, "", "(GCP only) cloud dns zone name")
	initCmd.Flags().StringVar(&gcpOrgID, FlagGCPOrgID, "", "(GCP only) project organization ID")
	initCmd.Flags().StringVar(&gcpBillingID, FlagGCPBillingID, "", "(GCP only) project organization ID")

	// Azure
	initCmd.Flags().StringVar(&azResourceGroup, FlagAZResourceGroup, "", "(Azure only) resource group name")
	initCmd.Flags().StringVar(&azSubscriptionID, FlagAZSubscriptionID, "", "(Azure only) subscription ID")
	initCmd.Flags().StringVar(&azTenantID, FlagAZTenantID, "", "(Azure only) tenant ID")
	initCmd.Flags().BoolVar(&azUseCLI, FlagAZLoginViaCLI, false, "(Azure only) login via CLI")
}

func preRunInit(cmd *cobra.Command, args []string) error {
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
		if errors.As(err, &cli.ConfigFileNotFound{}) {
			c, err = cli.NewConfigFromScratch(name, prv, region)
			if err != nil {
				log.Error().Msgf("unable to create config from scratch: %s", err)
				return err
			}
		} else {
			log.Error().Msgf("unable to create config from file: %s", err)
			return err
		}
	}
	c.LogLevel = logLevel
	target := cli.InitDone
	log.Info().Msgf("[%s->%s] running init on project %s", c.Status, target, c.Name)
	if c.Name != name || c.Provider != prv {
		return fmt.Errorf("please run a clean before changing project name or provider")
	}

	if err := c.SetDistro(distro); err != nil {
		return err
	}
	if err := c.SetEdition(edition); err != nil {
		return err
	}

	log.Debug().Msgf("input: %t - deploy nomad: %t", deployNomad, c.DeployNomad)
	c.DeployNomad = deployNomad
	c.Save()

	if err := c.SetDomain(domain); err != nil {
		return fmt.Errorf("error setting domain: %w", err)
	}

	if err = processFlags(c); err != nil {
		return err
	}

	p, err := getProvider(ctx, c)
	if err != nil {
		return err
	}

	c.SaveStatus(cli.InitRunning)

	if err := initRepos(c, branch); err != nil {
		log.Error().Msgf("error: %s", err)
		return err
	}

	if err := initProvider(c, p); err != nil {
		log.Error().Msgf("error during init: %s", err)
		return err
	}

	if c.Status < cli.InitDone {
		c.SaveStatus(cli.InitDone)
	}

	log.Info().Msgf("[%s->%s] completed init on project %s", c.Status, target, c.Name)
	return nil
}

func initProvider(c *cli.Config, p provider.Provider) error {
	log.Info().Msgf("initializing cloud resources")
	if err := p.InitProvider(ctx); err != nil {
		return fmt.Errorf("error initing provider: %w", err)
	}

	templates, err := p.GetTemplates(ctx)
	if err != nil {
		return fmt.Errorf("failed to get templates: %w", err)
	}

	log.Info().Msgf("generating terraform config files on: %s", c.WorkdirProject)
	if err := os.MkdirAll(c.WorkdirProject, 0777); err != nil {
		return err
	}
	for _, t := range templates {
		log.Info().Msgf("generating %v: %s", t.Name, t.Path)
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
	c.Save()
	// checkout repos
	git := git.NewGit("bitrockteam", logLevel)
	for _, repo := range c.Repos {
		err := git.Clone(repo, filepath.Join(".caravan", c.Name, repo), b)
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
			FlagGCPDnsZone:       gcpDNSZone,
			FlagGCPOrgID:         gcpOrgID,
			FlagGCPBillingID:     gcpBillingID,
		}
		for param, value := range requiredFlags {
			if err2 := mustBeNonEmpty(value, param, provider.GCP); err2 != nil {
				err = multierror.Append(err, err2)
			}
		}

		c.GCPParentProject = gcpParentProject
		c.GCPDNSZone = gcpDNSZone
		c.GCPBillingID = gcpBillingID
		c.SetGCPOrgID(gcpOrgID)
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
