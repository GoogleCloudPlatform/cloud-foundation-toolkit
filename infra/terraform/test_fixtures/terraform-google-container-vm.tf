locals {
  container_vm_required_roles = [
    "roles/compute.admin",
    "roles/iam.serviceAccountUser",
  ]
}

resource "google_project" "container_vm" {

  provider = "google.phoogle"

  name = "ci-container-vm"
  project_id = "ci-container-vm"
  folder_id = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  billing_account = "${module.variables.phoogle_billing_account}"
}

resource "google_project_services" "container_vm" {

  provider = "google.phoogle"

  project = "${google_project.container_vm.id}"
  services = [
    "bigquery-json.googleapis.com",
    "compute.googleapis.com",
    "container.googleapis.com",
    "containerregistry.googleapis.com",
    "oslogin.googleapis.com",
    "pubsub.googleapis.com",
    "storage-api.googleapis.com",
  ]
}

resource "google_service_account" "container_vm" {

  provider = "google.phoogle"

  project = "${google_project.container_vm.id}"
  account_id = "ci-container-vm"
  display_name = "ci-container-vm"
}

resource "google_project_iam_binding" "container_vm" {

  provider = "google.phoogle"

  count = "${length(local.container_vm_required_roles)}"

  project = "${google_project_services.container_vm.project}"
  role = "${element(local.container_vm_required_roles, count.index)}"

  members = [
    "serviceAccount:${google_service_account.container_vm.email}",
  ]
}

resource "google_service_account_key" "container_vm" {

  provider = "google.phoogle"

  service_account_id = "${google_service_account.container_vm.id}"
}

resource "random_id" "container_vm_github_webhook_token" {
  byte_length = 20
}

data "template_file" "container_vm_github_webhook_url" {

  template = "https://concourse.infra.cft.tips/api/v1/teams/cft/pipelines/$${pipeline}/resources/pull-request/check/webhook?webhook_token=$${webhook_token}"

  vars {
    pipeline = "terraform-google-container-vm"
    webhook_token = "${random_id.container_vm_github_webhook_token.hex}"
  }
}

resource "kubernetes_secret" "container_vm" {
  metadata {
    namespace = "concourse-cft"
    name = "container-vm"
  }
  data {
    github_webhook_token = "${random_id.container_vm_github_webhook_token.hex}"
    phoogle_project_id = "${google_project.container_vm.id}"
    phoogle_sa = "${base64decode(google_service_account_key.container_vm.private_key)}"
  }
}
