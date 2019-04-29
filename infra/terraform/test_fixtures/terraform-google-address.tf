locals {
  address_required_roles = [
    "roles/compute.networkAdmin",
    "roles/dns.admin",
    "roles/iam.serviceAccountUser",
  ]
}

resource "google_project" "address" {
  provider = "google.phoogle"

  name            = "ci-address"
  project_id      = "ci-address"
  folder_id       = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  billing_account = "${module.variables.phoogle_billing_account}"
}

resource "google_project_services" "address" {
  provider = "google.phoogle"

  project = "${google_project.address.id}"

  services = [
    "oslogin.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "compute.googleapis.com",
    "dns.googleapis.com",
  ]
}

resource "google_service_account" "address" {
  provider = "google.phoogle"

  project      = "${google_project.address.id}"
  account_id   = "ci-address"
  display_name = "ci-address"
}

resource "google_project_iam_member" "address" {
  provider = "google.phoogle"

  count = "${length(local.address_required_roles)}"

  project = "${google_project_services.address.project}"
  role    = "${element(local.address_required_roles, count.index)}"
  member  = "serviceAccount:${google_service_account.address.email}"
}

resource "google_service_account_key" "address" {
  provider = "google.phoogle"

  service_account_id = "${google_service_account.address.id}"
}

resource "random_id" "address_github_webhook_token" {
  byte_length = 20
}

data "template_file" "address_github_webhook_url" {
  template = "https://concourse.infra.cft.tips/api/v1/teams/cft/pipelines/$${pipeline}/resources/pull-request/check/webhook?webhook_token=$${webhook_token}"

  vars {
    pipeline      = "terraform-google-address"
    webhook_token = "${random_id.address_github_webhook_token.hex}"
  }
}

resource "kubernetes_secret" "address" {
  metadata {
    namespace = "concourse-cft"
    name      = "address"
  }

  data {
    github_webhook_token = "${random_id.address_github_webhook_token.hex}"
    phoogle_project_id   = "${google_project.address.id}"
    phoogle_sa           = "${base64decode(google_service_account_key.address.private_key)}"
  }
}
