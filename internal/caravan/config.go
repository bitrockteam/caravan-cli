package caravan

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/asaskevich/govalidator"
)

type Config struct {
	Name                   string              `json:",omitempty"`
	Region                 string              `json:",omitempty"`
	Regions                map[string][]string `json:",omitempty"`
	Profile                string              `json:",omitempty"`
	Provider               string              `json:",omitempty"`
	Providers              []string            `json:",omitempty"`
	Branch                 string              `json:",omitempty"`
	TableName              string              `json:",omitempty"`
	BucketName             string              `json:",omitempty"`
	Repos                  []string            `json:",omitempty"`
	Domain                 string              `json:",omitempty"`
	Workdir                string              `json:",omitempty"`
	WorkdirProject         string              `json:",omitempty"`
	WorkdirBaking          string              `json:",omitempty"`
	WorkdirBakingVars      string              `json:",omitempty"`
	WorkdirInfra           string              `json:",omitempty"`
	WorkdirInfraVars       string              `json:",omitempty"`
	WorkdirInfraBackend    string              `json:",omitempty"`
	WorkdirPlatform        string              `json:",omitempty"`
	WorkdirPlatformVars    string              `json:",omitempty"`
	WorkdirPlatformBackend string              `json:",omitempty"`
	Force                  bool                `json:",omitempty"`
	Status                 Status              `json:",omitempty"`
}

func NewConfigFromScratch(name, provider, region string) (c *Config, err error) {
	wd := ".caravan"
	repos := []string{"caravan", "caravan-baking", "caravan-platform", "caravan-application-support"}

	providers := []string{"aws"}

	if len(name) > 12 {
		return c, fmt.Errorf("name too long %d: max length is 12", len(name))
	}

	c = &Config{
		Name:           name,
		Profile:        "default",
		BucketName:     name + "-caravan-terraform-state",
		TableName:      name + "-caravan-terraform-state-lock",
		Repos:          repos,
		Providers:      providers,
		Domain:         "reactive-labs.io",
		Workdir:        wd,
		WorkdirProject: wd + "/" + name,
	}
	if provider != "" {
		err = c.setProvider(provider)
	}
	if region != "" {
		err = c.setRegion(region)
	}
	return c, err
}

func NewConfigFromFile() (c *Config, err error) {
	wd := ".caravan"
	b, err := ioutil.ReadFile(filepath.Join(wd, "caravan.state"))
	if err != nil {
		return c, err
	}

	if err := json.Unmarshal(b, &c); err != nil {
		return c, err
	}
	return c, nil
}

func (c *Config) SetWorkdir(wd string) {
	c.Workdir = wd
	c.WorkdirProject = filepath.Join(wd, c.Name)
	c.WorkdirInfra = filepath.Join(c.WorkdirProject, "caravan-infra-"+c.Provider)
	c.WorkdirInfraVars = filepath.Join(c.WorkdirInfra, c.Name+"-infra.tfvars")
	c.WorkdirInfraBackend = filepath.Join(c.WorkdirInfra, c.Name+"-backend.tf")
	c.WorkdirBakingVars = filepath.Join(c.WorkdirProject, "caravan-baking", "terraform", c.Provider+"-baking.tfvars")
	c.WorkdirBaking = filepath.Join(c.WorkdirProject, "caravan-baking", "terraform")
	c.WorkdirPlatform = filepath.Join(c.WorkdirProject, "caravan-platform")
	c.WorkdirPlatformVars = filepath.Join(c.WorkdirProject, "caravan-platform", c.Name+"-"+c.Provider+".tfvars")
	c.WorkdirPlatformBackend = filepath.Join(c.WorkdirProject, "caravan-platform", "backend.tf")
}

func (c *Config) setProvider(provider string) (err error) {
	for _, v := range c.Providers {
		if v == provider {
			c.Repos = append(c.Repos, "caravan-infra-"+v)
			c.Provider = provider
			c.WorkdirInfra = filepath.Join(c.WorkdirProject, "caravan-infra-"+provider)
			c.WorkdirInfraVars = filepath.Join(c.WorkdirInfra, c.Name+"-infra.tfvars")
			c.WorkdirInfraBackend = filepath.Join(c.WorkdirInfra, c.Name+"-backend.tf")
			c.WorkdirBaking = filepath.Join(c.WorkdirProject, "caravan-baking", "terraform")
			c.WorkdirBakingVars = filepath.Join(c.WorkdirProject, "caravan-baking", "terraform", c.Provider+"-baking.tfvars")
			c.WorkdirPlatform = filepath.Join(c.WorkdirProject, "caravan-platform")
			c.WorkdirPlatformVars = filepath.Join(c.WorkdirProject, "caravan-platform", c.Name+"-"+c.Provider+".tfvars")
			c.WorkdirPlatformBackend = filepath.Join(c.WorkdirProject, "caravan-platform", "backend.tf")
			return nil
		}
	}
	return fmt.Errorf("provider not supported: %s - %v", provider, c.Providers)
}

func (c *Config) setRegion(region string) (err error) {
	if isValidRegion(c.Provider, region) {
		c.Region = region
		return nil
	}
	return fmt.Errorf("please provide a valid region")
}

func (c *Config) SetDomain(domain string) (err error) {
	if isValidDomain(domain) {
		c.Domain = domain
		return nil
	}
	return fmt.Errorf("please provide a valid domain name")
}

func (c *Config) SetBranch(branch string) {
	c.Branch = branch
}

// check the name of the region for the given provider.
func isValidRegion(provider, region string) bool {
	// TODO temp method, use SDK resources to validate
	if provider == "aws" {
		if region == "eu-south-1" {
			return true
		}
	}
	return false
}

func (c *Config) SaveConfig() (err error) {
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(c.Workdir, os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(c.Workdir, "caravan.state"), data, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) StatusReport() {
	t, err := template.New("status").Parse(`======
PROJECT: {{.Name }}@{{.Branch }}
CARAVAN STATUS: {{ printf "%d" .Status}}-{{.Status}}

provider: {{.Provider}}
region: {{ or .Region "-- provider default --"}}
{{ if gt .Status 3 }}
 VAULT
URL: https://vault.{{.Name }}.{{.Domain}}
status: {{.VaultCheck}}
version: {{.VaultVersion}} 
 CONSUL
URL: https://consul.{{.Name }}.{{.Domain}} 
status: {{.ConsulCheck }}
version: {{.ConsulVersion}}
 NOMAD
URL: https://nomad.{{.Name }}.{{.Domain}} 
status: {{.NomadCheck}} 
version: {{.NomadVersion}}
{{ end }}
======`)
	if err != nil {
		fmt.Printf("error parsing report: %s\n", err)
	}

	if err := t.Execute(os.Stdout, c); err != nil {
		fmt.Printf("error executing report: %s\n", err)
	}
}

// check if the provided string is a valid domain name.
func isValidDomain(domain string) bool {
	return govalidator.IsDNSName(domain)
}

func (c *Config) VaultCheck() string {
	v := NewVaultHealth("https://vault."+c.Name+"."+c.Domain+"/", c.WorkdirInfra+"/ca_certs.pem")
	return v.Check()
}

func (c *Config) VaultVersion() string {
	v := NewVaultHealth("https://vault."+c.Name+"."+c.Domain+"/", c.WorkdirInfra+"/ca_certs.pem")
	return v.Version()
}

func (c *Config) ConsulCheck() bool {
	co := NewConsulHealth("https://consul."+c.Name+"."+c.Domain+"/", c.WorkdirInfra+"/ca_certs.pem")
	return co.Check()
}

func (c *Config) ConsulVersion() string {
	co := NewConsulHealth("https://consul."+c.Name+"."+c.Domain+"/", c.WorkdirInfra+"/ca_certs.pem")
	return co.Version()
}
func (c *Config) NomadCheck() bool {
	n := NewNomadHealth("https://nomad."+c.Name+"."+c.Domain+"/", c.WorkdirInfra+"/ca_certs.pem")
	return n.Check()
}

func (c *Config) NomadVersion() string {
	n := NewNomadHealth("https://nomad."+c.Name+"."+c.Domain+"/", c.WorkdirInfra+"/ca_certs.pem")
	return n.Version()
}
