package cli

import (
	"caravan-cli/vault"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/asaskevich/govalidator"
)

// Config is the main configuration data structure that is persisted to JSON.
type Config struct {
	Name                      string              `json:",omitempty"`
	Region                    string              `json:",omitempty"`
	Regions                   map[string][]string `json:",omitempty"`
	Profile                   string              `json:",omitempty"`
	Provider                  string              `json:",omitempty"`
	Branch                    string              `json:",omitempty"`
	LockName                  string              `json:",omitempty"`
	StateStoreName            string              `json:",omitempty"`
	Repos                     []string            `json:",omitempty"`
	Domain                    string              `json:",omitempty"`
	Workdir                   string              `json:",omitempty"`
	WorkdirProject            string              `json:",omitempty"`
	WorkdirBaking             string              `json:",omitempty"`
	WorkdirBakingVars         string              `json:",omitempty"`
	WorkdirInfra              string              `json:",omitempty"`
	WorkdirInfraVars          string              `json:",omitempty"`
	WorkdirInfraBackend       string              `json:",omitempty"`
	WorkdirPlatform           string              `json:",omitempty"`
	WorkdirPlatformVars       string              `json:",omitempty"`
	WorkdirPlatformBackend    string              `json:",omitempty"`
	WorkdirApplication        string              `json:",omitempty"`
	WorkdirApplicationVars    string              `json:",omitempty"`
	WorkdirApplicationBackend string              `json:",omitempty"`
	Force                     bool                `json:",omitempty"`
	Status                    Status              `json:",omitempty"`
	VaultRootToken            string              `json:",omitempty"`
	NomadToken                string              `json:",omitempty"`
	VaultURL                  string              `json:",omitempty"`
	CAPath                    string              `json:",omitempty"`
	ServiceAccount            string              `json:",omitempty"`
	Datacenter                string              `json:",omitempty"`
	DeployNomad               bool                `json:",omitempty"`
	LinuxOSFamily             string              `json:",omitempty"`
	LinuxOSVersion            string              `json:",omitempty"`
	LinuxOS                   string              `json:",omitempty"`
	Edition                   string              `json:",omitempty"`

	GCPConfig
	AzureConfig
}

// NewConfigFromScratch is used to construct a minimal configuration when no state
// is yet persisted on a local state file.
func NewConfigFromScratch(name, provider, region string) (c *Config, err error) {
	wd := ".caravan"
	repos := []string{"caravan-platform", "caravan-application-support"}

	c = &Config{
		Name:           name,
		Profile:        "default",
		StateStoreName: name + "-caravan-terraform-state",
		LockName:       name + "-caravan-terraform-state-lock",
		Repos:          repos,
		Workdir:        wd,
		WorkdirProject: wd + "/" + name,
		ServiceAccount: name + "-terraform",
		Datacenter:     provider + "-dc",
		Region:         region,
	}

	c.SetWorkdir(wd, provider)

	return c, err
}

// NewConfigFromFile constructs a configuration from  the content of the state file (caravan.state).
func NewConfigFromFile() (c *Config, err error) {
	wd := ".caravan"
	var b []byte

	if b, err = os.ReadFile(filepath.Join(wd, "caravan.state")); err != nil {
		return c, ConfigFileNotFound{Err: err}
	}

	if err = json.Unmarshal(b, &c); err != nil {
		return c, ConfigFileUnreadable{Err: err}
	}
	return c, nil
}

// SetWorkdir is used for white box testing.
func (c *Config) SetWorkdir(wd, provider string) {
	c.Workdir = wd
	c.Repos = append(c.Repos, "caravan-infra-"+provider)
	c.Provider = provider
	c.WorkdirProject = filepath.Join(wd, c.Name)
	c.WorkdirInfra = filepath.Join(c.WorkdirProject, "caravan-infra-"+c.Provider)
	c.WorkdirInfraVars = filepath.Join(c.WorkdirInfra, c.Name+"-infra.tfvars")
	c.WorkdirInfraBackend = filepath.Join(c.WorkdirInfra, c.Name+"-backend.tf")
	c.WorkdirBakingVars = filepath.Join(c.WorkdirProject, "caravan-baking", "terraform", c.Provider+"-baking.tfvars")
	c.WorkdirBaking = filepath.Join(c.WorkdirProject, "caravan-baking", "terraform")
	c.WorkdirPlatform = filepath.Join(c.WorkdirProject, "caravan-platform")
	c.WorkdirPlatformVars = filepath.Join(c.WorkdirProject, "caravan-platform", c.Name+"-"+c.Provider+"-cli.tfvars")
	c.WorkdirPlatformBackend = filepath.Join(c.WorkdirProject, "caravan-platform", "backend.tf")
	c.WorkdirApplication = filepath.Join(c.WorkdirProject, "caravan-application-support")
	c.WorkdirApplicationVars = filepath.Join(c.WorkdirProject, "caravan-application-support", c.Name+"-"+c.Provider+"-cli.tfvars")
	c.WorkdirApplicationBackend = filepath.Join(c.WorkdirProject, "caravan-application-support", "backend.tf")
	c.CAPath = filepath.Join(c.WorkdirInfra, "ca_certs.pem")
}

func (c *Config) SetDomain(domain string) (err error) {
	if isValidDomain(domain) {
		c.Domain = domain
		c.VaultURL = "https://vault." + c.Name + "." + c.Domain
		return nil
	}
	return fmt.Errorf("please provide a valid domain name")
}

func (c *Config) SetBranch(branch string) {
	c.Branch = branch
}

// SetVaultRootToken reads the content of the token file into config.
func (c *Config) SetVaultRootToken() error {
	// TODO consolidate in constructor
	vrt, err := os.ReadFile(filepath.Join(c.WorkdirInfra, "."+c.Name+"-root_token"))
	if err != nil {
		return err
	}
	// TODO make more robust
	c.VaultRootToken = string(vrt[0 : len(vrt)-1])
	return nil
}

// SetDistro sets the linux ditribution.
func (c *Config) SetDistro(d string) (err error) {
	if len(strings.Split(d, "-")) < 2 {
		return fmt.Errorf("unsupported linux distribution")
	}
	c.LinuxOS = strings.Split(d, "-")[0]
	c.LinuxOSVersion = strings.Split(d, "-")[1]
	switch d {
	case "ubuntu-2204", "ubuntu-2104", "ubuntu-2004":
		c.LinuxOSFamily = "debian"
	case "centos-7", "centos-8":
		c.LinuxOSFamily = "redhat"
	default:
		return fmt.Errorf("unsupported linux distribution")
	}
	return nil
}

// SetNomadToken reads into config the Nomad Token.
func (c *Config) SetNomadToken() error {
	v, err := vault.New(c.VaultURL, c.VaultRootToken, c.CAPath)
	if err != nil {
		return err
	}

	t, err := v.GetToken("nomad/creds/token-manager")
	if err != nil {
		return err
	}
	log.Info().Msgf("setting nomad token: %s", t)
	c.NomadToken = t

	return nil
}

// Save serializes to JSON the configuration and a local state store (caravan.state).
func (c *Config) Save() (err error) {
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(c.Workdir, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(c.Workdir, "caravan.state"), data, 0600)
	if err != nil {
		return err
	}
	return nil
}

// SetEdition sets the edition value in the config.
func (c *Config) SetEdition(edition string) error {
	if edition == "os" || edition == "ent" {
		c.Edition = edition
		return nil
	}
	return fmt.Errorf("only os or ent are allowed: %s", edition)
}

// isValidDomain checks if the provided string is a valid domain name.
func isValidDomain(domain string) bool {
	return govalidator.IsDNSName(domain)
}
