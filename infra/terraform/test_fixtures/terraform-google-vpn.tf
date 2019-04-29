locals {
  vpn_required_service_account_project_roles = [
    "roles/compute.networkAdmin",
    "roles/compute.instanceAdmin",
    "roles/iam.serviceAccountUser",
  ]

  vpn_required_apis = [
    "cloudresourcemanager.googleapis.com",
    "compute.googleapis.com",
  ]
}

resource "google_project" "vpn" {
  provider = "google.phoogle"

  name            = "ci-vpn"
  project_id      = "ci-vpn"
  folder_id       = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  billing_account = "${module.variables.phoogle_billing_account}"
}

resource "google_project_service" "vpn" {
  provider = "google.phoogle"
  project  = "${google_project.vpn.id}"
  count    = "${length(local.vpn_required_apis)}"
  service  = "${element(local.vpn_required_apis, count.index)}"
}

resource "google_service_account" "vpn" {
  provider     = "google.phoogle"
  project      = "${google_project.vpn.id}"
  account_id   = "ci-vpn"
  display_name = "ci-vpn"
}

resource "google_project_iam_member" "vpn" {
  provider = "google.phoogle"
  project  = "${google_project.vpn.id}"
  count    = "${length(local.vpn_required_service_account_project_roles)}"
  role     = "${element(local.vpn_required_service_account_project_roles, count.index)}"
  member   = "serviceAccount:${google_service_account.vpn.email}"
}

resource "google_service_account_key" "vpn" {
  provider = "google.phoogle"

  service_account_id = "${google_service_account.vpn.id}"
}

resource "random_id" "vpn_github_webhook_token" {
  byte_length = 20
}

data "template_file" "vpn_github_webhook_url" {
  template = "https://concourse.infra.cft.tips/api/v1/teams/cft/pipelines/$${pipeline}/resources/pull-request/check/webhook?webhook_token=$${webhook_token}"

  vars {
    pipeline      = "terraform-google-vpn"
    webhook_token = "${random_id.vpn_github_webhook_token.hex}"
  }
}

resource "kubernetes_secret" "ci_vpn" {
  metadata {
    namespace = "concourse-cft"
    name      = "vpn"
  }

  data {
    github_webhook_token = "${random_id.vpn_github_webhook_token.hex}"
    phoogle_project_id   = "${google_project.vpn.id}"
    phoogle_sa           = "${base64decode(google_service_account_key.vpn.private_key)}"
  }
}
