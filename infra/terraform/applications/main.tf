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
    prefix = "applications"
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

data "terraform_remote_state" "postgres" {
  backend = "gcs"
  config {
    bucket = "cloud-foundation-cicd-tfstate"
    prefix = "postgres"
  }
  workspace = "${terraform.workspace}"
}

data "google_container_cluster" "cicd" {
  name   = "${data.terraform_remote_state.gke.cluster_name}"
  region = "${module.variables.region[terraform.workspace]}"
}

data "google_client_config" "current" {}

data "google_dns_managed_zone" "tips_cft_infra" {
  name = "tips-cft-infra"
}

provider "kubernetes" {
  host = "${data.google_container_cluster.cicd.endpoint}"
  cluster_ca_certificate = "${base64decode(data.google_container_cluster.cicd.master_auth.0.cluster_ca_certificate)}"
  token = "${data.google_client_config.current.access_token}"
  load_config_file = false
}

provider "helm" {
  install_tiller = false
  service_account = "${kubernetes_cluster_role_binding.tiller.metadata.0.name}"
  kubernetes {
    config_path = "/dev/null"
    host = "${data.google_container_cluster.cicd.endpoint}"
    cluster_ca_certificate = "${base64decode(data.google_container_cluster.cicd.master_auth.0.cluster_ca_certificate)}"
    token = "${data.google_client_config.current.access_token}"
  }
}
