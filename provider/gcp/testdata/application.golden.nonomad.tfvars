
vault_endpoint  = "https://vault.test-name.test.me"
consul_endpoint = "https://consul.test-name.test.me"
domain          = "test-name.test.me"

artifacts_source_prefix = "gcs::https://www.googleapis.com/storage/v1/cfgs-test-name"
services_domain         = "service.consul"
dc_names                = ["gcp-dc"]
cloud                   = "gcp"

jenkins_volume_external_id = ""

vault_skip_tls_verify = true
consul_insecure_https = true
ca_cert_file          = "../caravan-infra-gcp/ca_certs.pem"
