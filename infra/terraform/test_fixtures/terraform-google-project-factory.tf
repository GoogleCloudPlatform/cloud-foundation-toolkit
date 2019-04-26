locals {
  project_factory_required_roles = [
    "roles/resourcemanager.projectCreator",
    "roles/compute.xpnAdmin",
    "roles/compute.networkAdmin",
    "roles/iam.serviceAccountAdmin",
    "roles/resourcemanager.projectIamAdmin",
  ]
}

resource "google_project" "project_factory" {

  provider = "google.phoogle"

  name = "ci-project-factory"
  project_id = "ci-project-factory"
  folder_id = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  billing_account = "${module.variables.phoogle_billing_account}"
}

resource "google_compute_shared_vpc_host_project" "shared_vpc_host" {

  provider = "google.phoogle"

  project = "${google_project.project_factory.project_id}"
}

resource "google_project_services" "project_factory" {

  provider = "google.phoogle"

  project = "${google_project.project_factory.id}"
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

resource "google_service_account" "project_factory" {

  provider = "google.phoogle"

  project = "${google_project.project_factory.id}"
  account_id = "ci-project-factory"
  display_name = "ci-project-factory"
}

resource "google_folder_iam_binding" "project_factory" {

  provider = "google.phoogle"

  count = "${length(local.project_factory_required_roles)}"

  folder = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  role   = "${element(local.project_factory_required_roles, count.index)}"

  members = [
    "serviceAccount:${google_service_account.project_factory.email}",
  ]
}

resource "google_service_account_key" "project_factory" {

  provider = "google.phoogle"

  service_account_id = "${google_service_account.project_factory.id}"
}

resource "random_id" "project_factory_github_webhook_token" {
  byte_length = 20
}

data "template_file" "project_factory_github_webhook_url" {

  template = "https://concourse.infra.cft.tips/api/v1/teams/cft/pipelines/$${pipeline}/resources/pull-request/check/webhook?webhook_token=$${webhook_token}"

  vars {
    pipeline = "terraform-google-project-factory"
    webhook_token = "${random_id.project_factory_github_webhook_token.hex}"
  }
}

resource "google_folder" "project_factory_ci_projects" {

  provider = "google.phoogle"

  display_name = "project-factory-ci-projects"
  parent       = "${google_folder.phoogle_cloud_foundation_cicd.name}"
}

resource "kubernetes_secret" "concourse_cft_project_factory" {
  metadata {
    namespace = "concourse-cft"
    name = "project-factory"
  }
  data {
    github_webhook_token = "${random_id.project_factory_github_webhook_token.hex}"
    phoogle_folder_id = "${google_folder.project_factory_ci_projects.name}"
    phoogle_project_id = "${google_project.project_factory.id}"
    phoogle_sa = "${base64decode(google_service_account_key.project_factory.private_key)}"
  }
}
