locals {
  cloud_nat_required_roles = [
    "roles/compute.networkAdmin",
    "roles/iam.serviceAccountUser",
  ]
}

resource "google_project" "cloud_nat" {
  provider        = "google.phoogle"
  name            = "ci-cloud-nat"
  project_id      = "ci-cloud-nat"
  folder_id       = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  billing_account = "${module.variables.phoogle_billing_account}"
}

resource "google_project_services" "cloud_nat" {
  provider = "google.phoogle"
  project  = "${google_project.cloud_nat.id}"

  services = [
    "compute.googleapis.com",
  ]
}

resource "google_service_account" "cloud_nat" {
  provider     = "google.phoogle"
  project      = "${google_project.cloud_nat.id}"
  account_id   = "ci-cloud-nat"
  display_name = "ci-cloud-nat"
}

resource "google_project_iam_binding" "cloud_nat" {
  provider = "google.phoogle"
  count    = "${length(local.cloud_nat_required_roles)}"
  project  = "${google_project_services.cloud_nat.project}"
  role     = "${element(local.cloud_nat_required_roles, count.index)}"

  members = [
    "serviceAccount:${google_service_account.cloud_nat.email}",
  ]
}

resource "google_service_account_key" "cloud_nat" {
  provider           = "google.phoogle"
  service_account_id = "${google_service_account.cloud_nat.id}"
}

resource "random_id" "cloud_nat_github_webhook_token" {
  byte_length = 20
}

data "template_file" "cloud_nat_github_webhook_url" {
  template = "https://concourse.infra.cft.tips/api/v1/teams/cft/pipelines/$${pipeline}/resources/pull-request/check/webhook?webhook_token=$${webhook_token}"

  vars {
    pipeline      = "terraform-google-cloud-nat"
    webhook_token = "${random_id.cloud_nat_github_webhook_token.hex}"
  }
}

resource "kubernetes_secret" "cloud_nat" {
  metadata {
    namespace = "concourse-cft"
    name      = "cloud-nat"
  }

  data {
    github_webhook_token = "${random_id.cloud_nat_github_webhook_token.hex}"
    phoogle_project_id   = "${google_project.cloud_nat.id}"
    phoogle_sa           = "${base64decode(google_service_account_key.cloud_nat.private_key)}"
  }
}
