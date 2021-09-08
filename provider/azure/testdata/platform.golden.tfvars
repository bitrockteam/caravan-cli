vault_endpoint  = "https://vault.test-name.test.me"
consul_endpoint = "https://consul.test-name.test.me"
nomad_endpoint  = "https://nomad.test-name.test.me"

vault_skip_tls_verify = true
consul_insecure_https = true
ca_cert_file          = "../caravan-infra-azure/ca_certs.pem"

auth_providers = ["azure"]

azure_bootstrap_resource_group_name  = "caravan-test-rg"
azure_bootstrap_storage_account_name = "sg-test-01"
azure_bootstrap_client_id            = "client1"
azure_bootstrap_client_secret        = "pass1"
azure_bootstrap_tenant_id            = "my-tenant-111"
azure_bootstrap_subscription_id      = "111-222-333"

bootstrap_state_backend_provider   = "azure"
bootstrap_state_object_name_prefix = "infraboot/terraform/state"
