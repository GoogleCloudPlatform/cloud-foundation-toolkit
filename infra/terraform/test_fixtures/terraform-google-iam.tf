locals {
  // Roles required by CI service accounts to run the IAM
  // integration tests.
  iam_required_roles = [
    "roles/resourcemanager.organizationAdmin",
    "roles/compute.xpnAdmin",
    "roles/storage.admin",
    "roles/pubsub.admin",
    "roles/cloudkms.admin",
  ]
}

resource "google_project" "iam" {
  provider = "google.phoogle"
  name = "ci-iam"
  project_id = "ci-iam"
  folder_id = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  billing_account = "${module.variables.phoogle_billing_account}"
}

resource "google_project_services" "iam" {
  provider = "google.phoogle"

  project = "${google_project.iam.id}"
  services = [
    "admin.googleapis.com",
    "appengine.googleapis.com",
    "cloudbilling.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "compute.googleapis.com",
    "iam.googleapis.com",
    "iamcredentials.googleapis.com",
    "oslogin.googleapis.com",
    "serviceusage.googleapis.com",
  ]
}

resource "google_service_account" "iam" {
  provider = "google.phoogle"
  project      = "${google_project.iam.id}"
  account_id   = "ci-iam"
  display_name = "ci-iam"
}

resource "google_folder_iam_member" "iam" {
  provider = "google.phoogle"
  count = "${length(local.iam_required_roles)}"
  folder = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  role   = "${element(local.iam_required_roles, count.index)}"
  member = "serviceAccount:${google_service_account.iam.email}"
}

resource "google_service_account_key" "iam" {
  provider = "google.phoogle"
  service_account_id = "${google_service_account.iam.id}"
}

resource "random_id" "iam_github_webhook_token" {
  byte_length = 20
}

data "template_file" "iam_github_webhook_url" {
  template = "https://concourse.infra.cft.tips/api/v1/teams/cft/pipelines/$${pipeline}/resources/pull-request/check/webhook?webhook_token=$${webhook_token}"
  vars {
    pipeline = "terraform-google-iam"
    webhook_token = "${random_id.iam_github_webhook_token.hex}"
  }
}

resource "kubernetes_secret" "concourse_cft_iam" {
  metadata {
    namespace = "concourse-cft"
    name = "iam"
  }
  data {
    github_webhook_token = "${random_id.iam_github_webhook_token.hex}"
    phoogle_project_id = "${google_project.iam.id}"
    phoogle_sa = "${base64decode(google_service_account_key.iam.private_key)}"
  }
}