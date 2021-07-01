package caravan

import (
	"fmt"
	"html/template"
	"os"
)

func (c *Config) StatusReport() {
	t, err := template.New("status").Parse(`====== {{.Name }}@{{or .Branch "default"}}
status: {{ printf "%d" .Status}}-{{.Status}}
provider: {{.Provider}} region: {{ or .Region "default"}}
{{- if gt .Status 3 }}
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
{{- end }}
======`)

	if err != nil {
		fmt.Printf("error parsing report: %s\n", err)
	}

	if err := t.Execute(os.Stdout, c); err != nil {
		fmt.Printf("error executing report: %s\n", err)
	}
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
