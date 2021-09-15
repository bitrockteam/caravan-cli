terraform {
  backend "s3" {
     bucket = "test-name-caravan-terraform-state"
     prefix = "platform/terraform/state"
     credentials = ".test-name-terraform-sa-key.json"
  }
}
