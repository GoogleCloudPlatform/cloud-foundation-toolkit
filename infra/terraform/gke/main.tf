provider "google" {
  region = "${module.variables.region[terraform.workspace]}"
}

terraform {
  backend "gcs" {
    bucket = "cloud-foundation-cicd-tfstate"
    prefix = "gke"
  }
}

data "terraform_remote_state" "networking" {
  backend = "gcs"
  config {
    bucket = "cloud-foundation-cicd-tfstate"
    prefix = "networking"
  }
  workspace = "${terraform.workspace}"
}

module "gke" {

  source  = "terraform-google-modules/kubernetes-engine/google"
  version = "0.3.0"

  project_id = "${module.variables.project_id}"

  name               = "${module.variables.name_prefix}-${terraform.workspace}"
  regional           = true
  region             = "${module.variables.region[terraform.workspace]}"
  network            = "${data.terraform_remote_state.networking.network_name}"
  subnetwork         = "${data.terraform_remote_state.networking.subnet_name}"
  ip_range_pods      = "${data.terraform_remote_state.networking.subnet_range_pods_name}"
  ip_range_services  = "${data.terraform_remote_state.networking.subnet_range_services_name}"
  kubernetes_version = "latest"

  node_pools = [
    {
      name      = "pool-00"
      min_count = "${var.pool_00_min_count[terraform.workspace]}"
      max_count = "${var.pool_00_max_count[terraform.workspace]}"
    },
  ]

  node_pools_labels = {
    all = {}
    pool-00 = {}
  }

  node_pools_taints = {
    all = []
    pool-00 = []
  }

  node_pools_tags = {
    all = []
    pool-00 = []
  }
}
