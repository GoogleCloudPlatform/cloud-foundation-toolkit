locals {
  # roles/storage.admin used by stdlib::get_from_bucket tests.
  required_service_account_project_roles = [
    "roles/compute.instanceAdmin.v1",
    "roles/iam.serviceAccountUser",
    "roles/storage.admin",
  ]
}

resource "google_project" "startup_scripts" {
  provider = "google.phoogle"

  name            = "ci-startup-scripts-v2"
  project_id      = "ci-startup-scripts-v2"
  folder_id       = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  billing_account = "${module.variables.phoogle_billing_account}"
}

resource "google_project_services" "startup_scripts" {
  provider = "google.phoogle"

  project = "${google_project.startup_scripts.id}"

  services = [
    "cloudresourcemanager.googleapis.com",
    "storage-api.googleapis.com",
    "compute.googleapis.com",
    "oslogin.googleapis.com",
  ]
}

resource "google_service_account" "startup_scripts" {
  provider = "google.phoogle"

  project      = "${google_project.startup_scripts.id}"
  account_id   = "ci-startup-scripts"
  display_name = "ci-startup-scripts"
}

resource "google_project_iam_member" "startup_scripts" {
  provider = "google.phoogle"
  count    = "${length(local.required_service_account_project_roles)}"
  project  = "${google_project_services.startup_scripts.project}"
  role     = "${element(local.required_service_account_project_roles, count.index)}"
  member   = "serviceAccount:${google_service_account.startup_scripts.email}"
}

resource "google_service_account_key" "startup_scripts" {
  provider = "google.phoogle"

  service_account_id = "${google_service_account.startup_scripts.id}"
}

resource "random_id" "startup_scripts_github_webhook_token" {
  byte_length = 20
}

data "template_file" "startup_scripts_github_webhook_url" {
  template = "https://concourse.infra.cft.tips/api/v1/teams/cft/pipelines/$${pipeline}/resources/pull-request/check/webhook?webhook_token=$${webhook_token}"

  vars {
    pipeline      = "terraform-google-startup-scripts"
    webhook_token = "${random_id.startup_scripts_github_webhook_token.hex}"
  }
}

resource "kubernetes_secret" "ci_startup_scripts" {
  metadata {
    namespace = "concourse-cft"
    name      = "startup-scripts"
  }

  data {
    github_webhook_token = "${random_id.startup_scripts_github_webhook_token.hex}"
    phoogle_project_id   = "${google_project.startup_scripts.id}"
    phoogle_sa           = "${base64decode(google_service_account_key.startup_scripts.private_key)}"
  }
}
