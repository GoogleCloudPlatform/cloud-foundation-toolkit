locals {
  jenkins_required_roles = [
    "roles/compute.admin",
    "roles/compute.networkAdmin",
    "roles/iam.serviceAccountUser",
  ]
}

resource "google_project" "jenkins" {
  provider = "google.phoogle"

  name            = "phoogle-ci-jenkins-project"
  project_id      = "phoogle-ci-jenkins-project"
  folder_id       = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  billing_account = "${module.variables.phoogle_billing_account}"
}

resource "google_project_services" "jenkins" {
  provider = "google.phoogle"

  project = "${google_project.jenkins.id}"

  services = [
    "compute.googleapis.com",
    "storage-api.googleapis.com",
  ]
}

resource "google_service_account" "jenkins" {
  provider = "google.phoogle"

  project      = "${google_project.jenkins.id}"
  account_id   = "phoogle-ci-jenkins-id"
  display_name = "phoogle-ci-jenkins-id"
}

resource "google_project_iam_member" "jenkins" {
  provider = "google.phoogle"

  count = "${length(local.jenkins_required_roles)}"

  project = "${google_project_services.jenkins.project}"
  role    = "${element(local.jenkins_required_roles, count.index)}"
  member  = "serviceAccount:${google_service_account.jenkins.email}"
}

resource "google_project_iam_member" "jenkins_service_account" {
  provider = "google.phoogle"

  project = "${google_project.jenkins.id}"

  role   = "roles/compute.instanceAdmin.v1"
  member = "serviceAccount:${google_project.jenkins.number}@cloudservices.gserviceaccount.com"
}

resource "google_project_iam_member" "jenkins_service_account_user" {
  provider = "google.phoogle"

  project = "${google_project.jenkins.id}"

  role   = "roles/iam.serviceAccountUser"
  member = "serviceAccount:${google_project.jenkins.number}@cloudservices.gserviceaccount.com"
}

resource "google_service_account_key" "jenkins" {
  provider = "google.phoogle"

  service_account_id = "${google_service_account.jenkins.id}"
}

resource "random_id" "jenkins_github_webhook_token" {
  byte_length = 20
}

data "template_file" "jenkins_github_webhook_url" {
  template = "https://concourse.infra.cft.tips/api/v1/teams/cft/pipelines/$${pipeline}/resources/pull-request/check/webhook?webhook_token=$${webhook_token}"

  vars {
    pipeline      = "terraform-google-jenkins"
    webhook_token = "${random_id.jenkins_github_webhook_token.hex}"
  }
}

resource "kubernetes_secret" "jenkins" {
  metadata {
    namespace = "concourse-cft"
    name      = "jenkins"
  }

  data {
    github_webhook_token = "${random_id.jenkins_github_webhook_token.hex}"
    phoogle_project_id   = "${google_project.jenkins.id}"
    phoogle_sa           = "${base64decode(google_service_account_key.jenkins.private_key)}"
  }
}