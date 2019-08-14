locals {
  kubernetes_engine_required_roles = [
    "roles/compute.networkAdmin",
    "roles/compute.viewer",
    "roles/container.clusterAdmin",
    "roles/container.developer",
    "roles/iam.serviceAccountAdmin",
    "roles/iam.serviceAccountUser",
  ]
}

resource "google_project" "ci_kubernetes_engine" {
  provider = "google.phoogle"

  name            = "ci-kubernetes-engine"
  project_id      = "ci-kubernetes-engine"
  folder_id       = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  billing_account = "${module.variables.phoogle_billing_account}"
}

resource "google_project_services" "ci_kubernetes_engine" {
  provider = "google.phoogle"

  project = "${google_project.ci_kubernetes_engine.id}"

  services = [
    "bigquery-json.googleapis.com",
    "cloudkms.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "compute.googleapis.com",
    "container.googleapis.com",
    "containerregistry.googleapis.com",
    "iam.googleapis.com",
    "iamcredentials.googleapis.com",
    "oslogin.googleapis.com",
    "pubsub.googleapis.com",
    "serviceusage.googleapis.com",
    "storage-api.googleapis.com",
  ]
}

resource "google_service_account" "ci_kubernetes_engine" {
  provider = "google.phoogle"

  project      = "${google_project.ci_kubernetes_engine.id}"
  account_id   = "ci-kubernetes-engine"
  display_name = "ci-kubernetes-engine"
}

resource "google_project_iam_binding" "ci_kubernetes_engine" {
  provider = "google.phoogle"

  count = "${length(local.kubernetes_engine_required_roles)}"

  project = "${google_project_services.ci_kubernetes_engine.project}"
  role    = "${element(local.kubernetes_engine_required_roles, count.index)}"

  members = [
    "serviceAccount:${google_service_account.ci_kubernetes_engine.email}",
  ]
}

resource "google_project_iam_binding" "ci_kubernetes_engine_kms_access" {
  provider = "google.phoogle"

  project = "${google_project_services.ci_kubernetes_engine.project}"
  role    = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  members = [
    "serviceAccount:service-${google_project.ci_kubernetes_engine.number}@container-engine-robot.iam.gserviceaccount.com",
  ]
}

resource "google_project_iam_binding" "ci_kubernetes_engine_kms_admin_access" {
  provider = "google.phoogle"

  project = "${google_project_services.ci_kubernetes_engine.project}"
  role    = "roles/cloudkms.admin"

  members = [
    "serviceAccount:service-${google_project.ci_kubernetes_engine.number}@container-engine-robot.iam.gserviceaccount.com",
  ]
}

resource "google_service_account_key" "ci_kubernetes_engine" {
  provider = "google.phoogle"

  service_account_id = "${google_service_account.ci_kubernetes_engine.id}"
}

resource "random_id" "kubernetes_engine_github_webhook_token" {
  byte_length = 20
}

data "template_file" "kubernetes_engine_github_webhook_url" {
  template = "https://concourse.infra.cft.tips/api/v1/teams/cft/pipelines/$${pipeline}/resources/pull-request/check/webhook?webhook_token=$${webhook_token}"

  vars {
    pipeline      = "terraform-google-kubernetes-engine"
    webhook_token = "${random_id.kubernetes_engine_github_webhook_token.hex}"
  }
}

resource "kubernetes_secret" "ci_kubernetes_engine" {
  metadata {
    namespace = "concourse-cft"
    name      = "kubernetes-engine"
  }

  data {
    github_webhook_token = "${random_id.kubernetes_engine_github_webhook_token.hex}"
    phoogle_project_id   = "${google_project.ci_kubernetes_engine.id}"
    phoogle_sa           = "${base64decode(google_service_account_key.ci_kubernetes_engine.private_key)}"
  }
}
