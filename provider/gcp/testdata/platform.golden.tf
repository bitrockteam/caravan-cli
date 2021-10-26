terraform {
  backend "gcs" {
    bucket      = "test-name-caravan-terraform-state"
    prefix      = "platform/terraform/state"
    credentials = "../caravan-infra-gcp/.test-name-terraform-sa-key.json"
  }
}
