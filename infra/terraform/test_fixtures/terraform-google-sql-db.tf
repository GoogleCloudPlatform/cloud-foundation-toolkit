locals {
  sql_db_required_roles = [
    "roles/compute.networkAdmin",
    "roles/cloudsql.admin",
  ]
}

resource "google_project" "sql_db" {
  provider = "google.phoogle"
  name = "ci-sql-db"
  project_id = "ci-sql-db"
  folder_id = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  billing_account = "${module.variables.phoogle_billing_account}"
}

resource "google_project_services" "sql_db" {
  provider = "google.phoogle"

  project = "${google_project.sql_db.id}"
  services = [
    "sqladmin.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "compute.googleapis.com",
    "servicenetworking.googleapis.com",
  ]
}

resource "google_service_account" "sql_db" {
  provider = "google.phoogle"
  project      = "${google_project.sql_db.id}"
  account_id   = "ci-sql-db"
  display_name = "ci-sql-db"
}

resource "google_folder_iam_member" "sql_db" {
  provider = "google.phoogle"
  count = "${length(local.sql_db_required_roles)}"
  folder = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  role   = "${element(local.sql_db_required_roles, count.index)}"
  member = "serviceAccount:${google_service_account.sql_db.email}"
}

resource "google_service_account_key" "sql_db" {
  provider = "google.phoogle"
  service_account_id = "${google_service_account.sql_db.id}"
}

resource "random_id" "sql_db_github_webhook_token" {
  byte_length = 20
}

data "template_file" "sql_db_github_webhook_url" {
  template = "https://concourse.infra.cft.tips/api/v1/teams/cft/pipelines/$${pipeline}/resources/pull-request/check/webhook?webhook_token=$${webhook_token}"
  vars {
    pipeline = "terraform-google-sql-db"
    webhook_token = "${random_id.sql_db_github_webhook_token.hex}"
  }
}

resource "kubernetes_secret" "concourse_cft_sql_db" {
  metadata {
    namespace = "concourse-cft"
    name = "sql-db"
  }
  data {
    github_webhook_token = "${random_id.sql_db_github_webhook_token.hex}"
    phoogle_project_id = "${google_project.sql_db.id}"
    phoogle_sa = "${base64decode(google_service_account_key.sql_db.private_key)}"
  }
}