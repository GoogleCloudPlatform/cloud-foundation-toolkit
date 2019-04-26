provider "google" {
  project = "${module.variables.project_id}"
  region  = "${module.variables.region[terraform.workspace]}"
}

provider "google" {
  alias       = "phoogle"
  credentials = "${var.phoogle_credentials_path}"
}

terraform {
  backend "gcs" {
    bucket = "cloud-foundation-cicd-tfstate"
    prefix = "test_fixtures"
  }
}

data "terraform_remote_state" "gke" {
  backend = "gcs"
  config {
    bucket = "cloud-foundation-cicd-tfstate"
    prefix = "gke"
  }
  workspace = "${terraform.workspace}"
}

data "google_container_cluster" "cicd" {
  name   = "${data.terraform_remote_state.gke.cluster_name}"
  region = "${module.variables.region[terraform.workspace]}"
}

data "google_client_config" "current" {}

provider "kubernetes" {
  host = "${data.google_container_cluster.cicd.endpoint}"
  cluster_ca_certificate = "${base64decode(data.google_container_cluster.cicd.master_auth.0.cluster_ca_certificate)}"
  token = "${data.google_client_config.current.access_token}"
  load_config_file = false
}

# All fixtures are to be created under this folder.

resource "google_folder" "phoogle_cloud_foundation_cicd" {

  provider = "google.phoogle"

  display_name = "cloud-foundation-cicd"
  parent       = "organizations/${var.phoogle_org_id}"
}
