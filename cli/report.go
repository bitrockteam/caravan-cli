package cli

import (
	"caravan-cli/health"
	"context"
	"html/template"
	"os"

	"github.com/rs/zerolog/log"
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
		log.Error().Msgf("error parsing report: %s\n", err)
	}

	if err := t.Execute(os.Stdout, r); err != nil {
		log.Error().Msgf("error executing report: %s\n", err)
	}
}

func (r *Report) VaultCheck(ctx context.Context) string {
	v := health.NewVaultHealth("https://vault."+r.Caravan.Name+"."+r.Caravan.Domain+"/", r.Caravan.CAPath)
	return v.Check(ctx)
}

func (r *Report) VaultVersion(ctx context.Context) string {
	v := health.NewVaultHealth("https://vault."+r.Caravan.Name+"."+r.Caravan.Domain+"/", r.Caravan.CAPath)
	return v.Version(ctx)
}

func (r *Report) ConsulCheck(ctx context.Context) bool {
	co := health.NewConsulHealth("https://consul."+r.Caravan.Name+"."+r.Caravan.Domain+"/", r.Caravan.CAPath, r.Caravan.Datacenter)
	return co.Check(ctx)
}

func (r *Report) ConsulVersion(ctx context.Context) string {
	co := health.NewConsulHealth("https://consul."+r.Caravan.Name+"."+r.Caravan.Domain+"/", r.Caravan.CAPath, r.Caravan.Datacenter)
	return co.Version(ctx)
}
func (r *Report) NomadCheck(ctx context.Context) bool {
	n := health.NewNomadHealth("https://nomad."+r.Caravan.Name+"."+r.Caravan.Domain+"/", r.Caravan.CAPath)
	return n.Check(ctx)
}

func (r *Report) NomadVersion(ctx context.Context) string {
	n := health.NewNomadHealth("https://nomad."+r.Caravan.Name+"."+r.Caravan.Domain+"/", r.Caravan.CAPath)
	return n.Version(ctx)
}
