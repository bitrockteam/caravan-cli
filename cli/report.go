package cli

import (
	"caravan-cli/cli/checker"
	"context"
	"fmt"
	"html/template"
	"os"

	"github.com/rs/zerolog/log"
)

type Target struct {
	Name string
}

const (
	Vault  = "vault"
	Nomad  = "nomad"
	Consul = "consul"
)

type Report struct {
	Caravan *Config
	Tools   map[string]Tool
	Targets []string
}

type Tool struct {
	Status  bool
	Version string
}

func NewReport(c *Config) (r *Report) {
	targets := []string{Vault, Consul}
	if c.DeployNomad {
		targets = append(targets, Nomad)
	}

	r = &Report{
		Targets: targets,
		Caravan: c,
		Tools:   map[string]Tool{},
	}
	return r
}

func (r *Report) CheckStatus(ctx context.Context) (err error) {
	// check CA for https
	if _, err := os.Stat(r.Caravan.CAPath); os.IsNotExist(err) {
		return nil
	}
	dc := func(gc *checker.GenericChecker) {
		gc.Datacenter = r.Caravan.Datacenter
	}
	for _, t := range r.Targets {
		var h checker.Checker
		switch t {
		case Nomad:
			h, err = checker.NewNomadChecker(fmt.Sprintf("https://%s.%s.%s", t, r.Caravan.Name, r.Caravan.Domain), r.Caravan.CAPath)
			if err != nil {
				return err
			}
		case Consul:
			h, err = checker.NewConsulChecker(fmt.Sprintf("https://%s.%s.%s", t, r.Caravan.Name, r.Caravan.Domain), r.Caravan.CAPath, dc)
			if err != nil {
				return err
			}
		case Vault:
			h, err = checker.NewVaultChecker(fmt.Sprintf("https://%s.%s.%s", t, r.Caravan.Name, r.Caravan.Domain), r.Caravan.CAPath)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported target")
		}
		status := h.Status(ctx)
		version := h.Version(ctx)
		c := Tool{Status: status, Version: version}
		r.Tools[t] = c
	}
	return nil
}

func (r *Report) PrintReport() {
	t, err := template.New("status").Parse(`
Name:		{{.Caravan.Name }}@{{or .Caravan.Branch "default"}}
Status:		{{.Caravan.Status}}
Provider:	{{.Caravan.Provider}} 
Region:		{{ or .Caravan.Region "default"}}
Domain:		{{ .Caravan.Domain }}
DeployNomad:    {{ .Caravan.DeployNomad }}
Linux Distro:   {{ .Caravan.LinuxOS }}-{{ .Caravan.LinuxOSVersion }}
{{- if gt .Caravan.Status 3 }}
{{ range $k,$v:= .Tools }}
{{ $k }}
	URL:		https://{{ $k }}.{{ $.Caravan.Name}}.{{ $.Caravan.Domain }}
	Status:		{{ $v.Status}}
	Version:	{{ $v.Version}}
{{- end }}
{{- end }}
`)

	if err != nil {
		log.Error().Msgf("error parsing report: %s\n", err)
	}

	if err := t.Execute(os.Stdout, r); err != nil {
		log.Error().Msgf("error executing report: %s\n", err)
	}
}
