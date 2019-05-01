locals {
  event_function_required_roles = [
    "roles/cloudfunctions.developer",
    "roles/compute.instanceAdmin.v1",
    "roles/iam.serviceAccountUser",
    "roles/logging.configWriter",
    "roles/pubsub.admin",
    "roles/source.admin",
    "roles/storage.admin",
  ]

  name = "ci-event-function"
}

resource "google_project" "event_function" {
  provider = "google.phoogle"

  name            = "${local.name}"
  project_id      = "${local.name}"
  folder_id       = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  billing_account = "${module.variables.phoogle_billing_account}"
}

resource "google_project_services" "event_function" {
  provider = "google.phoogle"

  project = "${google_project.event_function.id}"

  services = [
    "cloudfunctions.googleapis.com",
    "compute.googleapis.com",
    "logging.googleapis.com",
    "oslogin.googleapis.com",
    "pubsub.googleapis.com",
    "sourcerepo.googleapis.com",
    "storage-api.googleapis.com",
    "storage-component.googleapis.com",
  ]
}

resource "google_service_account" "event_function" {
  provider = "google.phoogle"

  project      = "${google_project.event_function.id}"
  account_id   = "${local.name}"
  display_name = "${local.name}"
}

resource "google_project_iam_member" "event_function" {
  provider = "google.phoogle"

  count = "${length(local.event_function_required_roles)}"

  project = "${google_project_services.event_function.project}"
  role    = "${element(local.event_function_required_roles, count.index)}"
  member  = "serviceAccount:${google_service_account.event_function.email}"
}

resource "google_service_account_key" "event_function" {
  provider = "google.phoogle"

  service_account_id = "${google_service_account.event_function.id}"
}

resource "random_id" "event_function_github_webhook_token" {
  byte_length = 20
}

data "template_file" "event_function_github_webhook_url" {
  template = "https://concourse.infra.cft.tips/api/v1/teams/cft/pipelines/$${pipeline}/resources/pull-request/check/webhook?webhook_token=$${webhook_token}"

  vars {
    pipeline      = "terraform-google-event-function"
    webhook_token = "${random_id.event_function_github_webhook_token.hex}"
  }
}

resource "kubernetes_secret" "event_function" {
  metadata {
    namespace = "concourse-cft"
    name      = "event-function"
  }

  data {
    github_webhook_token = "${random_id.event_function_github_webhook_token.hex}"
    phoogle_project_id   = "${google_project.event_function.id}"
    phoogle_sa           = "${base64decode(google_service_account_key.event_function.private_key)}"
  }
}
