terraform {
  backend "s3" {
    bucket         = "test-name-caravan-terraform-state"
    key            = "infraboot/terraform/state/terraform.tfstate"
    region         = "eu-south-1"
    dynamodb_table = "test-name-caravan-terraform-state-lock"
  }
}
