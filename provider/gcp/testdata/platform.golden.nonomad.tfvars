
vault_endpoint  = "https://vault.test-name.test.me"
consul_endpoint = "https://consul.test-name.test.me"
nomad_endpoint  = ""
enable_nomad    = false

bootstrap_state_backend_provider = "gcp"
auth_providers                   = ["gcp", "gsuite"]
gcp_project_id                   = "test-name"
gcp_csi                          = true
gcp_region                       = "europe-west6"
google_account_file              = "../caravan-infra-gcp/.test-name-terraform-sa-key.json"

gsuite_domain                = ""
gsuite_client_id             = ""
gsuite_client_secret         = ""
gsuite_default_role          = "bitrock"
gsuite_default_role_policies = ["default", "bitrock", "vault-admin-role"]
gsuite_allowed_redirect_uris = ["https://vault.test-name.test.me/ui/vault/auth/gsuite/oidc/callback", "https://vault.test-name.test.me/ui/vault/auth/oidc/oidc/callback"]

bootstrap_state_bucket_name        = "test-name-caravan-terraform-state"
bootstrap_state_object_name_prefix = "infraboot/terraform/state"
control_plane_role_name            = "control-plane"

vault_skip_tls_verify = true
consul_insecure_https = true
ca_cert_file          = "../caravan-infra-gcp/ca_certs.pem"
