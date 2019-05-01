locals {
  network_required_roles = [
    "roles/compute.networkAdmin",
    "roles/iam.serviceAccountUser",
  ]

  network_required_apis = [
    "cloudresourcemanager.googleapis.com",
    "compute.googleapis.com",
  ]
}

resource "google_project" "network" {
  provider = "google.phoogle"

  name            = "ci-network"
  project_id      = "ci-network"
  folder_id       = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  billing_account = "${module.variables.phoogle_billing_account}"
}

resource "google_project_service" "network" {
  provider = "google.phoogle"

  count   = "${length(local.network_required_apis)}"
  project = "${google_project.network.id}"
  service = "${element(local.network_required_apis, count.index)}"
}

resource "google_service_account" "network" {
  provider = "google.phoogle"

  project      = "${google_project.network.id}"
  account_id   = "ci-network"
  display_name = "ci-network"
}

resource "google_project_iam_member" "network" {
  provider = "google.phoogle"

  count = "${length(local.network_required_roles)}"

  project = "${google_project.network.project_id}"
  role    = "${element(local.network_required_roles, count.index)}"
  member  = "serviceAccount:${google_service_account.network.email}"
}

resource "google_service_account_key" "network" {
  provider = "google.phoogle"

  service_account_id = "${google_service_account.network.id}"
}

resource "random_id" "network_github_webhook_token" {
  byte_length = 20
}

data "template_file" "network_github_webhook_url" {
  template = "https://concourse.infra.cft.tips/api/v1/teams/cft/pipelines/$${pipeline}/resources/pull-request/check/webhook?webhook_token=$${webhook_token}"

  vars {
    pipeline      = "terraform-google-network"
    webhook_token = "${random_id.network_github_webhook_token.hex}"
  }
}

resource "kubernetes_secret" "network" {
  metadata {
    namespace = "concourse-cft"
    name      = "network"
  }

  data {
    github_webhook_token = "${random_id.network_github_webhook_token.hex}"
    phoogle_project_id   = "${google_project.network.id}"
    phoogle_sa           = "${base64decode(google_service_account_key.network.private_key)}"
  }
}
