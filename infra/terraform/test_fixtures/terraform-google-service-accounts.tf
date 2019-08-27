locals {
  service_accounts_required_roles = [
    "roles/iam.serviceAccountAdmin",
  ]
}

resource "google_project" "service_accounts" {
  provider = "google.phoogle"
  name = "ci-service-accounts"
  project_id = "ci-service-accounts"
  folder_id = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  billing_account = "${module.variables.phoogle_billing_account}"
}

resource "google_project_services" "service_accounts" {
  provider = "google.phoogle"

  project = "${google_project.service_accounts.id}"
  services = [
    "iam.googleapis.com",
    "iamcredentials.googleapis.com",
  ]
}

resource "google_service_account" "service_accounts" {
  provider = "google.phoogle"
  project      = "${google_project.service_accounts.id}"
  account_id   = "ci-service-accounts"
  display_name = "ci-service-accounts"
}

resource "google_folder_iam_member" "service_accounts" {
  provider = "google.phoogle"
  count = "${length(local.service_accounts_required_roles)}"
  folder = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  role   = "${element(local.service_accounts_required_roles, count.index)}"
  member = "serviceAccount:${google_service_account.service_accounts.email}"
}

resource "google_service_account_key" "service_accounts" {
  provider = "google.phoogle"
  service_account_id = "${google_service_account.service_accounts.id}"
}

resource "random_id" "service_accounts_github_webhook_token" {
  byte_length = 20
}

data "template_file" "service_accounts_github_webhook_url" {
  template = "https://concourse.infra.cft.tips/api/v1/teams/cft/pipelines/$${pipeline}/resources/pull-request/check/webhook?webhook_token=$${webhook_token}"
  vars {
    pipeline = "terraform-google-service-accounts"
    webhook_token = "${random_id.service_accounts_github_webhook_token.hex}"
  }
}

resource "kubernetes_secret" "concourse_cft_service_accounts" {
  metadata {
    namespace = "concourse-cft"
    name = "service-accounts"
  }
  data {
    github_webhook_token = "${random_id.service_accounts_github_webhook_token.hex}"
    phoogle_project_id = "${google_project.service_accounts.id}"
    phoogle_sa = "${base64decode(google_service_account_key.service_accounts.private_key)}"
  }
}