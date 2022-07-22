
region                = "europe-west6"
zone                  = "europe-west6-a"
project_id            = "test-name"
prefix                = "test-name"
external_domain       = "test.me"
use_le_staging        = true
dc_name               = "gcp-dc"
control_plane_sa_name = "control-plane"
worker_plane_sa_name  = "worker-plane"
image                 = "projects/parent-project/global/images/family/caravan-centos-image-ent"
parent_dns_project_id = "parent-project"
parent_dns_zone_name  = "dns-zone"
google_account_file   = ".test-name-terraform-sa-key.json"
enable_nomad          = false
