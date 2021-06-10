terraform {
  backend "s3" {
    bucket         = "test-bucket"
    key            = "infraboot/terraform/state/terraform.tfstate"
    region         = "test-region"
    dynamodb_table = "test-table"
  }
}
