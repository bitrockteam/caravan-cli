
vault_endpoint  = "https://vault.test-name.test.me"
consul_endpoint = "https://consul.test-name.test.me"
nomad_endpoint  = "https://nomad.test-name.test.me"

vault_skip_tls_verify = true
consul_insecure_https = true
ca_cert_file          = "../caravan-infra-aws/ca_certs.pem"

auth_providers = ["aws"]

aws_region                  = eu-south-1
aws_shared_credentials_file = "~/.aws/credentials"
aws_profile                 = "default"

bootstrap_state_backend_provider   = aws
bootstrap_state_bucket_name_prefix = test-name-caravan-terraform-state
bootstrap_state_object_name_prefix = "infraboot/terraform/state"
s3_bootstrap_region                = eu-south-1
