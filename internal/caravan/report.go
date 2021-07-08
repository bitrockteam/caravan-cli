package caravan

import (
	"fmt"
	"html/template"
	"os"
)

func (c *Config) StatusReport() {
	t, err := template.New("status").Parse(`Name:		{{.Name }}@{{or .Branch "default"}}
Status:		{{.Status}}
Provider:	{{.Provider}} 
Region:		{{ or .Region "default"}}
{{- if gt .Status 3 }}
Vault	URL:		https://vault.{{.Name }}.{{.Domain}}
	Status:		{{.VaultCheck}}
	Version:	{{.VaultVersion}}
	Token:		{{.VaultRootToken}}
Consul	URL:		https://consul.{{.Name }}.{{.Domain}} 
	Status:		{{.ConsulCheck }}
	Version:	{{.ConsulVersion}}
Nomad	URL: 		https://nomad.{{.Name }}.{{.Domain}} 
	Status:		{{.NomadCheck}} 
	Version:	{{.NomadVersion}}
	Token:		{{.NomadToken}}
{{- end }}
`)

	if err != nil {
		fmt.Printf("error parsing report: %s\n", err)
	}

	if err := t.Execute(os.Stdout, c); err != nil {
		fmt.Printf("error executing report: %s\n", err)
	}
}

func (c *Config) VaultCheck() string {
	v := NewVaultHealth("https://vault."+c.Name+"."+c.Domain+"/", c.CApath)
	return v.Check()
}

func (c *Config) VaultVersion() string {
	v := NewVaultHealth("https://vault."+c.Name+"."+c.Domain+"/", c.CApath)
	return v.Version()
}

func (c *Config) ConsulCheck() bool {
	co := NewConsulHealth("https://consul."+c.Name+"."+c.Domain+"/", c.CApath)
	return co.Check()
}

func (c *Config) ConsulVersion() string {
	co := NewConsulHealth("https://consul."+c.Name+"."+c.Domain+"/", c.CApath)
	return co.Version()
}
func (c *Config) NomadCheck() bool {
	n := NewNomadHealth("https://nomad."+c.Name+"."+c.Domain+"/", c.CApath)
	return n.Check()
}

func (c *Config) NomadVersion() string {
	n := NewNomadHealth("https://nomad."+c.Name+"."+c.Domain+"/", c.CApath)
	return n.Version()
}
