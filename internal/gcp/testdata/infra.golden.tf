terraform {
  backend "gcs" {
     bucket = "test-name-caravan-terraform-state"
     prefix = "infraboot/terraform/state"
     credentials = ".test-name-terraform-sa-key.json"
   }
}
