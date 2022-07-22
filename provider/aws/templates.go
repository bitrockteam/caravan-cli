package aws

const (
	bakingTfVarsTmpl = `
build_on_aws      = true
aws_region        = "{{ .Region }}"
aws_instance_type = "t3.small"
linux_os          = "{{ .LinuxOS }}"
linux_os_version  = "{{ .LinuxOSVersion }}"
linux_os_family   = "{{ .LinuxOSFamily }}"
ssh_username      = "{{ .LinuxOS }}"
`

	infraTfVarsTmpl = `
region                  = "{{ .Region }}"
awsprofile              = "{{ .Profile }}"
shared_credentials_file = "~/.aws/credentials"
prefix                  = "{{ .Name }}"
personal_ip_list        = ["0.0.0.0/0"]
use_le_staging          = true
external_domain         = "{{ .Domain }}"
tfstate_bucket_name     = "{{ .StateStoreName }}"
tfstate_table_name      = "{{ .LockName }}"
tfstate_region          = "{{ .Region }}"
ami_filter_name         = "caravan-{{ .Edition }}-{{ .LinuxOS }}-{{ .LinuxOSVersion }}-*"
ssh_username            = "{{ .LinuxOS }}"
`

	platformTfVarsTmpl = `
vault_endpoint  = "https://vault.{{.Name}}.{{.Domain}}"
consul_endpoint = "https://consul.{{.Name}}.{{.Domain}}"
nomad_endpoint  = "https://nomad.{{.Name}}.{{.Domain}}"

vault_skip_tls_verify = true
consul_insecure_https = true
ca_cert_file          = "../caravan-infra-aws/ca_certs.pem"

auth_providers = ["{{.Provider}}"]

aws_region                  = "{{ .Region }}"
aws_shared_credentials_file = "~/.{{.Provider}}/credentials"
aws_profile                 = "default"

bootstrap_state_backend_provider   = "{{ .Provider }}"
bootstrap_state_bucket_name_prefix = "{{ .StateStoreName }}"
bootstrap_state_object_name_prefix = "infraboot/terraform/state"
s3_bootstrap_region                = "{{ .Region }}"
`

	applicationTfVarsTmpl = `
vault_endpoint  = "https://vault.{{.Name}}.{{.Domain}}"
consul_endpoint = "https://consul.{{.Name}}.{{.Domain}}"
nomad_endpoint  = "https://nomad.{{.Name}}.{{.Domain}}"
domain = "{{.Name}}.{{.Domain}}"

artifacts_source_prefix    = ""
container_registry         = ""
services_domain            = "service.consul"
dc_names                   = ["{{.Provider}}-dc"]
cloud                      = "{{.Provider}}"
jenkins_volume_external_id = ""


vault_skip_tls_verify = true
consul_insecure_https = true
ca_cert_file          = "../caravan-infra-aws/ca_certs.pem"
`

	infraBackendTmpl = `
terraform {
  backend "s3" {
    bucket         = "{{ .StateStoreName }}"
    key            = "infraboot/terraform/state/terraform.tfstate"
    region         = "{{ .Region }}"
    dynamodb_table = "{{ .LockName }}"
  }
}
`
	platformBackendTmpl = `
terraform {
  backend "s3" {
    bucket         = "{{ .StateStoreName }}"
    key            = "platform/terraform/state/terraform.tfstate"
    region         = "{{ .Region }}"
    dynamodb_table = "{{ .LockName }}"
  }
}
`

	applicationSupportBackendTmpl = `
terraform {
  backend "s3" {
    bucket         = "{{ .StateStoreName }}"
    key            = "appsupport/terraform/state/terraform.tfstate"
    region         = "{{ .Region }}"
    dynamodb_table = "{{ .LockName }}"
  }
}
`
)
