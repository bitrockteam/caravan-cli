package azure

const (
	bakingTfVarsTmpl = `
build_on_azure              = true
build_image_name            = "caravan-centos-image"
azure_subscription_id       = "{{ or .AzureBakingSubscriptionID .AzureSubscriptionID }}"
azure_target_resource_group = "{{ or .AzureBakingResourceGroup .AzureResourceGroup }}"
azure_client_id             = "{{ or .AzureBakingClientID .AzureClientID }}"
azure_client_secret         = "{{ or .AzureBakingClientSecret .AzureClientSecret }}"
`
	infraTfVarsTmpl = `
resource_group_name        = "{{ .AzureResourceGroup }}"
image_resource_group_name  = "{{ or .AzureBakingResourceGroup .AzureResourceGroup }}"
parent_resource_group_name = "{{ or .AzureDNSResourceGroup .AzureResourceGroup }}"
storage_account_name       = "{{ .AzureStorageAccount }}"
prefix                     = "{{ .Name }}"
location                   = "{{ .Region }}"
external_domain            = "{{ .Domain }}"
client_id       = "{{ .AzureClientID }}"
client_secret   = "{{ .AzureClientSecret }}"
tenant_id       = "{{ .AzureTenantID }}"
subscription_id = "{{ .AzureSubscriptionID }}"
tags = {
  project   = "caravan-{{ .Name }}"
  managedBy = "terraform"
  repo      = "github.com/bitrockteam/caravan-infra-azure"
}
use_le_staging = true
`
	platformTfVarsTmpl = `
vault_endpoint  = "https://vault.{{.Name}}.{{.Domain}}"
consul_endpoint = "https://consul.{{.Name}}.{{.Domain}}"
nomad_endpoint  = "https://nomad.{{.Name}}.{{.Domain}}"

vault_skip_tls_verify = true
consul_insecure_https = true
ca_cert_file          = "../caravan-infra-azure/ca_certs.pem"

auth_providers = ["{{.Provider}}"]

azure_bootstrap_resource_group_name  = "{{ .AzureResourceGroup }}"
azure_bootstrap_storage_account_name = "{{ .AzureStorageAccount }}"
azure_bootstrap_client_id            = "{{ .AzureClientID }}"
azure_bootstrap_client_secret        = "{{ .AzureClientSecret }}"
azure_bootstrap_tenant_id            = "{{ .AzureTenantID }}"
azure_bootstrap_subscription_id      = "{{ .AzureSubscriptionID }}"

bootstrap_state_backend_provider   = "{{ .Provider }}"
bootstrap_state_object_name_prefix = "infraboot/terraform/state"
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
ca_cert_file          = "../caravan-infra-azure/ca_certs.pem"
`
	infraBackendTmpl = `
terraform {
  backend "azurerm" {
    resource_group_name  = "{{ .AzureResourceGroup }}"
    storage_account_name = "{{ .AzureStorageAccount }}"
    container_name       = "{{ .AzureStorageContainerName }}"
    key                  = "infraboot/terraform/state/terraform.tfstate"
    client_id            = "{{ .AzureClientID }}"
    client_secret        = "{{ .AzureClientSecret }}"
    tenant_id            = "{{ .AzureTenantID }}"
    subscription_id      = "{{ .AzureSubscriptionID }}"
  }
}
`
	platformBackendTmpl = `
terraform {
  backend "azurerm" {
    resource_group_name  = "{{ .AzureResourceGroup }}"
    storage_account_name = "{{ .AzureStorageAccount }}"
    container_name       = "{{ .AzureStorageContainerName }}"
    key                  = "platform/terraform/state/terraform.tfstate"
    client_id            = "{{ .AzureClientID }}"
    client_secret        = "{{ .AzureClientSecret }}"
    tenant_id            = "{{ .AzureTenantID }}"
    subscription_id      = "{{ .AzureSubscriptionID }}"
  }
}
`
	applicationSupportBackendTmpl = `
terraform {
  backend "azurerm" {
    resource_group_name  = "{{ .AzureResourceGroup }}"
    storage_account_name = "{{ .AzureStorageAccount }}"
    container_name       = "{{ .AzureStorageContainerName }}"
    key                  = "appsupport/terraform/state/terraform.tfstate"
    client_id            = "{{ .AzureClientID }}"
    client_secret        = "{{ .AzureClientSecret }}"
    tenant_id            = "{{ .AzureTenantID }}"
    subscription_id      = "{{ .AzureSubscriptionID }}"
  }
}
`
)
