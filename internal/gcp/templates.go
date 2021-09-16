package gcp

const (
	bakingTfVarsTmpl = `
build_on_google        = true
build_image_name       = "caravan-centos-image"
google_project_id      = "{{ .Caravan.Name }}"
google_account_file    = "YOUR-JSON-KEY"
google_network_name    = "caravan-gcp-vpc"
google_subnetwork_name = "caravan-gcp-subnet"
`

	infraTfVarsTmpl = `
region                = "{{ .Caravan.Region }}"
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
`

	platformTfVarsTmpl = `
vault_endpoint  = "https://vault.{{ .Caravan.Name }}.{{ .Caravan.Domain }}"
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
`

	applicationTfVarsTmpl = `
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
`

	infraBackendTmpl = `
  backend "gcs" {
     bucket = "{{ .Caravan.StateStoreName }}"
     prefix = "infraboot/terraform/state"
     credentials = ".{{ .Caravan.Name }}-terraform-sa-key.json"
   }
}
`

	platformBackendTmpl = `
  backend "s3" {
     bucket = "{{ .Caravan.StateStoreName }}"
     prefix = "platform/terraform/state"
     credentials = ".{{ .Caravan.Name }}-terraform-sa-key.json"
  }
}
`
	applicationSupportBackendTmpl = `
terraform {
  backend "gcs" {
     bucket = "{{ .Caravan.StateStoreName }}"
     prefix = "appsupport/terraform/state"
     credentials = ".{{ .Caravan.Name }}-terraform-sa-key.json"
   }
}
`
)
