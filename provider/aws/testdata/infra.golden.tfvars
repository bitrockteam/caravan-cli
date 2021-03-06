region                  = "eu-south-1"
awsprofile              = "default"
shared_credentials_file = "~/.aws/credentials"
prefix                  = "test-name"
personal_ip_list        = ["0.0.0.0/0"]
use_le_staging          = true
external_domain         = "test.me"
tfstate_bucket_name     = "test-name-caravan-terraform-state"
tfstate_table_name      = "test-name-caravan-terraform-state-lock"
tfstate_region          = "eu-south-1"
ami_filter_name         = "caravan-ent-ubuntu-2204-*"
ssh_username            = "ubuntu"
