locals {
  // Roles required by CI service accounts to run the folders
  // integration tests.
  folders_required_roles = [
    "roles/resourcemanager.folderCreator",
  ]
}

resource "google_project" "folders" {
  provider = "google.phoogle"
  name = "ci-folders"
  project_id = "ci-folders"
  folder_id = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  billing_account = "${module.variables.phoogle_billing_account}"
}

resource "google_project_services" "folders" {
  provider = "google.phoogle"

  project = "${google_project.folders.id}"
  services = [
    "cloudresourcemanager.googleapis.com",
  ]
}

resource "google_service_account" "folders" {
  provider = "google.phoogle"
  project      = "${google_project.folders.id}"
  account_id   = "ci-folders"
  display_name = "ci-folders"
}

resource "google_folder_iam_member" "folders" {
  provider = "google.phoogle"
  count = "${length(local.folders_required_roles)}"
  folder = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  role   = "${element(local.folders_required_roles, count.index)}"
  member = "serviceAccount:${google_service_account.folders.email}"
}

resource "google_service_account_key" "folders" {
  provider = "google.phoogle"
  service_account_id = "${google_service_account.folders.id}"
}

resource "random_id" "folders_github_webhook_token" {
  byte_length = 20
}

data "template_file" "folders_github_webhook_url" {
  template = "https://concourse.infra.cft.tips/api/v1/teams/cft/pipelines/$${pipeline}/resources/pull-request/check/webhook?webhook_token=$${webhook_token}"
  vars {
    pipeline = "terraform-google-folders"
    webhook_token = "${random_id.folders_github_webhook_token.hex}"
  }
}

resource "kubernetes_secret" "concourse_cft_folders" {
  metadata {
    namespace = "concourse-cft"
    name = "folders"
  }
  data {
    github_webhook_token = "${random_id.folders_github_webhook_token.hex}"
    phoogle_project_id = "${google_project.folders.id}"
    phoogle_sa = "${base64decode(google_service_account_key.folders.private_key)}"
  }
}