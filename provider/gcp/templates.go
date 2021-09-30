package gcp

const (
	bakingTfVarsTmpl = `
build_on_google        = true
build_image_name       = "caravan-centos-image"
google_project_id      = "{{ .Name }}"
google_network_name    = "caravan-gcp-vpc"
google_subnetwork_name = "caravan-gcp-subnet"
`

	infraTfVarsTmpl = `
region                = "{{ .Region }}"
zone                  = "{{ .Region }}-a"
project_id            = "{{ .Name }}"
prefix                = "{{ .Name }}"
external_domain       = "{{ .Domain }}"
use_le_staging        = true
dc_name               = "gcp-dc"
control_plane_sa_name = "control-plane"
worker_plane_sa_name  = "worker-plane"
image                 = "projects/{{ .GCPParentProject }}/global/images/family/caravan-centos-image-os"
parent_dns_project_id = "{{ .GCPParentProject }}"
parent_dns_zone_name  = "{{ .GCPDNSZone }}"
google_account_file   = ".{{ .Name }}-terraform-sa-key.json"
`

	platformTfVarsTmpl = `
vault_endpoint  = "https://vault.{{ .Name }}.{{ .Domain }}"
consul_endpoint = "https://consul.{{ .Name }}.{{ .Domain }}"
nomad_endpoint  = "https://nomad.{{ .Name }}.{{ .Domain }}"

bootstrap_state_backend_provider = "gcp"
auth_providers                   = ["gcp", "gsuite"]
gcp_project_id                   = "{{ .Name }}"
gcp_csi                          = true
gcp_region                       = "{{ .Region }}"
google_account_file              = "../caravan-infra-gcp/.{{ .Name }}-terraform-sa-key.json"

gsuite_domain                = ""
gsuite_client_id             = ""
gsuite_client_secret         = ""
gsuite_default_role          = "bitrock"
gsuite_default_role_policies = [ "default", "bitrock", "vault-admin-role" ]
gsuite_allowed_redirect_uris = [ "https://vault.{{ .Name }}.{{ .Domain }}/ui/vault/auth/gsuite/oidc/callback", "https://vault.{{ .Name }}.{{ .Domain }}/ui/vault/auth/oidc/oidc/callback"]

bootstrap_state_bucket_name_prefix = "states-bucket"
bootstrap_state_object_name_prefix = "infraboot/terraform/state"
control_plane_role_name            = "control-plane"

vault_skip_tls_verify = true
consul_insecure_https = true
ca_cert_file          = "../caravan-infra-{{ .Provider }}/ca_certs.pem"
`

	applicationTfVarsTmpl = `
vault_endpoint  = "https://vault.{{ .Name }}.{{ .Domain }}"
consul_endpoint = "https://consul.{{ .Name }}.{{ .Domain }}"
nomad_endpoint  = "https://nomad.{{ .Name }}.{{ .Domain }}"
domain = "{{ .Name }}.{{ .Domain }}"

artifacts_source_prefix    = ""
container_registry         = ""
services_domain            = "service.consul"
dc_names                   = ["{{ .Provider }}-dc"]
cloud                      = "{{ .Provider }}"
jenkins_volume_external_id = ""


vault_skip_tls_verify = true
consul_insecure_https = true
ca_cert_file          = "../caravan-infra-{{ .Provider }}/ca_certs.pem"
`

	infraBackendTmpl = `
terraform {
  backend "gcs" {
     bucket = "{{ .StateStoreName }}"
     prefix = "infraboot/terraform/state"
     credentials = ".{{ .Name }}-terraform-sa-key.json"
   }
}
`

	platformBackendTmpl = `
terraform {
  backend "s3" {
     bucket = "{{ .StateStoreName }}"
     prefix = "platform/terraform/state"
     credentials = ".{{ .Name }}-terraform-sa-key.json"
  }
}
`
	applicationSupportBackendTmpl = `
terraform {
  backend "gcs" {
     bucket = "{{ .StateStoreName }}"
     prefix = "appsupport/terraform/state"
     credentials = ".{{ .Name }}-terraform-sa-key.json"
   }
}
`
)
