resource_group_name        = "caravan-test-rg"
image_resource_group_name  = "caravan-admin"
parent_resource_group_name = "caravan-test-rg"
storage_account_name       = "sg-test-01"
prefix                     = "test-name"
location                   = "europewest"
external_domain            = "test.me"
client_id                  = "client1"
client_secret              = "pass1"
tenant_id                  = "my-tenant-111"
subscription_id            = "111-222-333"
image_name_regex           = "caravan-centos-image-ent-*"
tags = {
  project   = "caravan-test-name"
  managedBy = "terraform"
  repo      = "github.com/bitrockteam/caravan-infra-azure"
}
use_le_staging = true
