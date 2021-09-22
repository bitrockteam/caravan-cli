terraform {
  backend "gcs" {
     bucket = "test-name-caravan-terraform-state"
     prefix = "appsupport/terraform/state"
     credentials = ".test-name-terraform-sa-key.json"
   }
}
