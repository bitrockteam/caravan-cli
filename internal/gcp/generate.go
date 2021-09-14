package gcp

import (
	"caravan/internal/caravan"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// TODO refactor with common generate code.
func (g GCP) GenerateConfig() (err error) {
	fmt.Printf("generating config files on: %s\n", g.Caravan.WorkdirProject)
	if err := os.MkdirAll(g.Caravan.WorkdirProject, 0777); err != nil {
		return err
	}

	for _, t := range g.Templates {
		fmt.Printf("generating %v:%s \n", t.Name, t.Path)
		if err := g.Generate(t); err != nil {
			return err
		}
	}

	return nil
}

func loadTemplates(g GCP) []caravan.Template {
	return []caravan.Template{
		{
			Name: "baking-vars",
			Text: `build_on_google        = true
build_image_name       = "caravan-centos-image"
google_project_id      = "{{ .Caravan.Name }}"
google_account_file    = "YOUR-JSON-KEY"
google_network_name    = "caravan-gcp-vpc"
google_subnetwork_name = "caravan-gcp-subnet"
`,
			Path: g.Caravan.WorkdirBakingVars,
		},
		{
			Name: "infra-vars",
			Text: `region                = "{{ .Caravan.Region }}"
zone                  = "{{ .Caravan.Region }}-a"
project_id            = "{{ .Caravan.Name }}"
external_domain       = "{{ .Caravan.Domain }}"
use_le_staging        = true
dc_name               = "gcp-dc"
control_plane_sa_name = "control-plane"
worker_plane_sa_name  = "worker-plane"
image                 = "projects/{{ .Caravan.ParentProject }}/global/images/family/caravan-centos-image-os"
parent_dns_project_id = "{{ .Caravan.ParentProject }}"
parent_dns_zone_name  = "dns-example-zone"
`,
			Path: g.Caravan.WorkdirInfraVars,
		},
		{
			Name: "infra-backend",
			Text: `terraform {
  backend "gcs" {
     bucket = "{{ .Caravan.StateStoreName }}"
     prefix = "infraboot/terraform/state"
     credentials = ".{{ .Caravan.Name }}-terraform-sa-key.json"
   }
}
`,
			Path: g.Caravan.WorkdirInfraBackend,
		},
		{
			Name: "platform-backend",
			Text: `terraform {
  backend "s3" {
     bucket = "{{ .Caravan.StateStoreName }}"
     prefix = "platform/terraform/state"
     credentials = ".{{ .Caravan.Name }}-terraform-sa-key.json"
  }
}
`,
			Path: g.Caravan.WorkdirPlatformBackend,
		},
		{
			Name: "platform-vars",
			Text: `vault_endpoint  = "https://vault.{{ .Caravan.Name }}.{{ .Caravan.Domain }}"
consul_endpoint = "https://consul.{{ .Caravan.Name }}.{{ .Caravan.Domain }}"
nomad_endpoint  = "https://nomad.{{ .Caravan.Name }}.{{ .Caravan.Domain }}"

bootstrap_state_backend_provider = "gcp"
auth_providers                   = ["gcp", "gsuite"]
gcp_project_id                   = "{{ .Caravan.Name }}"
gcp_csi                          = true
gcp_region                       = "{{ .Caravan.Region }}"
google_account_file              = "../caravan-infra-gcp/.{{ .Caravan.Name }}-terraform-sa-key.json"

gsuite_domain                = ""
gsuite_client_id             = ""
gsuite_client_secret         = ""
gsuite_default_role          = "bitrock"
gsuite_default_role_policies = [ "default", "bitrock", "vault-admin-role" ]
gsuite_allowed_redirect_uris = [ "https://vault.{{ .Caravan.Name }}.{{ .Caravan.Domain }}/ui/vault/auth/gsuite/oidc/callback", "https://vault.{{ .Caravan.Name }}.{{ .Caravan.Domain }}/ui/vault/auth/oidc/oidc/callback"]

bootstrap_state_bucket_name_prefix = "states-bucket"
bootstrap_state_object_name_prefix = "infraboot/terraform/state"
control_plane_role_name            = "control-plane"

vault_skip_tls_verify = true
consul_insecure_https = true
ca_cert_file          = "../caravan-infra-{{ .Caravan.Provider }}/ca_certs.pem"
`,
			Path: g.Caravan.WorkdirPlatformVars,
		},
		{
			Name: "application-vars",
			Text: `
vault_endpoint  = "https://vault.{{ .Caravan.Name }}.{{ .Caravan.Domain }}"
consul_endpoint = "https://consul.{{ .Caravan.Name }}.{{ .Caravan.Domain }}"
nomad_endpoint  = "https://nomad.{{ .Caravan.Name }}.{{ .Caravan.Domain }}"
domain = "{{ .Caravan.Name }}.{{ .Caravan.Domain }}"

artifacts_source_prefix    = ""
container_registry         = ""
services_domain            = "service.consul"
dc_names                   = ["{{ .Caravan.Provider }}-dc"]
cloud                      = "{{ .Caravan.Provider }}"
jenkins_volume_external_id = ""


vault_skip_tls_verify = true
consul_insecure_https = true
ca_cert_file          = "../caravan-infra-{{ .Caravan.Provider }}/ca_certs.pem"
`,
			Path: g.Caravan.WorkdirApplicationVars,
		},
		{
			Name: "application-backend",
			Text: `terraform {
  backend "gcs" {
     bucket = "{{ .Caravan.StateStoreName }}"
     prefix = "appsupport/terraform/state"
     credentials = ".{{ .Caravan.Name }}-terraform-sa-key.json"
   }
}
`,
			Path: g.Caravan.WorkdirApplicationBackend,
		},
	}
}

func (g GCP) Generate(t caravan.Template) (err error) {
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

	if err := temp.Execute(f, g); err != nil {
		return err
	}
	return nil
}
