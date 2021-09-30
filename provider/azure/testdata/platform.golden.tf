terraform {
  backend "azurerm" {
    resource_group_name  = "caravan-test-rg"
    storage_account_name = "sg-test-01"
    container_name       = "tfstate"
    key                  = "platform/terraform/state/terraform.tfstate"
    client_id            = "client1"
    client_secret        = "pass1"
    tenant_id            = "my-tenant-111"
    subscription_id      = "111-222-333"
  }
}
