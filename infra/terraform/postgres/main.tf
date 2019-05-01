provider "google-beta" {
  project = "${module.variables.project_id}"
}

terraform {
  required_version = "0.11.13"

  backend "gcs" {
    bucket = "cloud-foundation-cicd-tfstate"
    prefix = "postgres"
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

resource "google_sql_database_instance" "postgres" {
  provider = "google-beta"

  name             = "${module.variables.name_prefix}-postgres-${terraform.workspace}"
  region           = "${module.variables.region[terraform.workspace]}"
  database_version = "POSTGRES_9_6"

  depends_on = ["google_service_networking_connection.postgres"]

  settings {
    tier = "db-g1-small"

    ip_configuration {
      ipv4_enabled    = "false"
      private_network = "${data.terraform_remote_state.networking.network_self_link}"
    }
  }
}

resource "google_sql_user" "concourse" {
  provider = "google-beta"

  name     = "concourse"
  instance = "${google_sql_database_instance.postgres.name}"
  password = "${var.postgres_concourse_user_password}"
}

resource "google_sql_database" "atc" {
  provider = "google-beta"

  name     = "atc"
  instance = "${google_sql_database_instance.postgres.name}"
}
