package cli

import (
	"caravan-cli/health"
	"fmt"
	"html/template"
	"os"
)

type Report struct {
	Caravan *Config
}

func (r *Report) StatusReport() {
	t, err := template.New("status").Parse(`
Name:		{{.Caravan.Name }}@{{or .Caravan.Branch "default"}}
Status:		{{.Caravan.Status}}
Provider:	{{.Caravan.Provider}} 
Region:		{{ or .Caravan.Region "default"}}
{{- if gt .Caravan.Status 3 }}
Vault	URL:		https://vault.{{.Caravan.Name }}.{{.Caravan.Domain}}
	Status:		{{.VaultCheck}}
	Version:	{{.VaultVersion}}
	Token:		{{.Caravan.VaultRootToken}}
Consul	URL:		https://consul.{{.Caravan.Name }}.{{.Caravan.Domain}} 
	Status:		{{.ConsulCheck }}
	Version:	{{.ConsulVersion}}
Nomad	URL: 		https://nomad.{{.Caravan.Name }}.{{.Caravan.Domain}} 
	Status:		{{.NomadCheck}} 
	Version:	{{.NomadVersion}}
	Token:		{{.Caravan.NomadToken}}
{{- end }}
`)

	if err != nil {
		fmt.Printf("error parsing report: %s\n", err)
	}

	if err := t.Execute(os.Stdout, r); err != nil {
		fmt.Printf("error executing report: %s\n", err)
	}
}

func (r *Report) VaultCheck() string {
	v := health.NewVaultHealth("https://vault."+r.Caravan.Name+"."+r.Caravan.Domain+"/", r.Caravan.CAPath)
	return v.Check()
}

func (r *Report) VaultVersion() string {
	v := health.NewVaultHealth("https://vault."+r.Caravan.Name+"."+r.Caravan.Domain+"/", r.Caravan.CAPath)
	return v.Version()
}

func (r *Report) ConsulCheck() bool {
	co := health.NewConsulHealth("https://consul."+r.Caravan.Name+"."+r.Caravan.Domain+"/", r.Caravan.CAPath, r.Caravan.Datacenter)
	return co.Check()
}

func (r *Report) ConsulVersion() string {
	co := health.NewConsulHealth("https://consul."+r.Caravan.Name+"."+r.Caravan.Domain+"/", r.Caravan.CAPath, r.Caravan.Datacenter)
	return co.Version()
}
func (r *Report) NomadCheck() bool {
	n := health.NewNomadHealth("https://nomad."+r.Caravan.Name+"."+r.Caravan.Domain+"/", r.Caravan.CAPath)
	return n.Check()
}

func (r *Report) NomadVersion() string {
	n := health.NewNomadHealth("https://nomad."+r.Caravan.Name+"."+r.Caravan.Domain+"/", r.Caravan.CAPath)
	return n.Version()
}
