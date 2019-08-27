locals {
  vm_required_roles = [
    "roles/compute.admin",
    "roles/compute.networkAdmin",
    "roles/iam.serviceAccountUser",
  ]
  vm_required_api_services = [
    "compute.googleapis.com",
    "iam.googleapis.com",
  ]
}

resource "google_project" "vm" {
  provider = "google.phoogle"

  name            = "ci-vm"
  project_id      = "ci-vm-module"
  folder_id       = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  billing_account = "${module.variables.phoogle_billing_account}"
}

resource "google_project_service" "vm" {
  provider = "google.phoogle"
  count    = "${length(local.vm_required_api_services)}"
  service  = "${element(local.vm_required_api_services, count.index)}"
  project  = "${google_project.vm.id}"
}

resource "google_service_account" "vm" {
  provider = "google.phoogle"

  project      = "${google_project.vm.id}"
  account_id   = "ci-google-vm"
  display_name = "ci-vm"
}

resource "google_project_iam_member" "vm" {
  provider = "google.phoogle"

  count = "${length(local.vm_required_roles)}"

  project = "${google_project.vm.project_id}"
  role    = "${element(local.vm_required_roles, count.index)}"
  member  = "serviceAccount:${google_service_account.vm.email}"
}

resource "google_project_iam_member" "vm_service_account" {
  provider = "google.phoogle"

  project = "${google_project.vm.id}"

  role   = "roles/compute.instanceAdmin.v1"
  member = "serviceAccount:${google_project.vm.number}@cloudservices.gserviceaccount.com"
}

resource "google_project_iam_member" "vm_service_account_user" {
  provider = "google.phoogle"

  project = "${google_project.vm.id}"

  role   = "roles/iam.serviceAccountUser"
  member = "serviceAccount:${google_project.vm.number}@cloudservices.gserviceaccount.com"
}

resource "google_service_account_key" "vm" {
  provider = "google.phoogle"

  service_account_id = "${google_service_account.vm.id}"
}

resource "random_id" "vm_github_webhook_token" {
  byte_length = 20
}

data "template_file" "vm_github_webhook_url" {
  template = "https://concourse.infra.cft.tips/api/v1/teams/cft/pipelines/$${pipeline}/resources/pull-request/check/webhook?webhook_token=$${webhook_token}"

  vars {
    pipeline      = "terraform-google-vm"
    webhook_token = "${random_id.vm_github_webhook_token.hex}"
  }
}

resource "kubernetes_secret" "vm" {
  metadata {
    namespace = "concourse-cft"
    name      = "vm"
  }

  data {
    github_webhook_token = "${random_id.vm_github_webhook_token.hex}"
    phoogle_project_id   = "${google_project.vm.id}"
    phoogle_sa           = "${base64decode(google_service_account_key.vm.private_key)}"
  }
}
