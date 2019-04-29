locals {
  log_export_required_roles = [
    # Needed for the Pubsub submodule to create a service account for the
    # subscription it creates
    "roles/iam.serviceAccountAdmin",

    # Needed for the cloud storage submodule to create/delete a bucket
    "roles/storage.admin",

    # Needed for the pubsub submodule to create/delete a pubsub topic
    "roles/pubsub.admin",

    # Needed for the bigquery submodule to create/delete a bigquery dataset
    "roles/bigquery.dataOwner",

    # Needed for the root module to activate APIs
    "roles/serviceusage.serviceUsageAdmin",

    # Needed for the Pubsub submodule to assign roles/bigquery.dataEditor to
    # the service account it creates
    "roles/resourcemanager.projectIamAdmin",
  ]

  log_export_billing_account_roles = [
    # Required to associate billing accounts to new projects
    "roles/billing.user",
  ]

  log_export_organization_roles = [
    # Required to create log sinks from the organization level on down
    "roles/logging.configWriter",

    # Required to associate billing accounts to new projects
    "roles/billing.projectManager",
  ]

  log_export_folder_roles = [
    # Required to spin up a project within the log_export folder
    "roles/resourcemanager.projectCreator",
  ]

  log_export_required_apis = [
    "cloudresourcemanager.googleapis.com",
    "oslogin.googleapis.com",
    "serviceusage.googleapis.com",
    "compute.googleapis.com",
    "bigquery-json.googleapis.com",
    "pubsub.googleapis.com",
    "storage-component.googleapis.com",
    "storage-api.googleapis.com",
    "logging.googleapis.com",
    "iam.googleapis.com",
    "cloudbilling.googleapis.com",
  ]
}

# Creating a dedicated folder to test out folder level log sink creation
resource "google_folder" "log_export" {
  provider = "google.phoogle"

  display_name = "ci-log-export"
  parent       = "${google_folder.phoogle_cloud_foundation_cicd.name}"
}

resource "google_project" "log_export" {
  provider = "google.phoogle"

  name            = "ci-log-export"
  project_id      = "ci-log-export"
  folder_id       = "${google_folder.log_export.name}"
  billing_account = "${module.variables.phoogle_billing_account}"
}

resource "google_project_service" "log_export" {
  provider = "google.phoogle"

  count   = "${length(local.log_export_required_apis)}"
  project = "${google_project.log_export.id}"
  service = "${element(local.log_export_required_apis, count.index)}"
}

resource "google_service_account" "log_export" {
  provider = "google.phoogle"

  project      = "${google_project.log_export.id}"
  account_id   = "ci-log-export"
  display_name = "ci-log-export"
}

resource "google_project_iam_member" "log_export" {
  provider = "google.phoogle"

  count = "${length(local.log_export_required_roles)}"

  project = "${google_project.log_export.project_id}"
  role    = "${element(local.log_export_required_roles, count.index)}"
  member  = "serviceAccount:${google_service_account.log_export.email}"
}

resource "google_billing_account_iam_member" "log_export" {
  provider = "google.phoogle"

  count = "${length(local.log_export_billing_account_roles)}"

  billing_account_id = "${module.variables.phoogle_billing_account}"
  role               = "${element(local.log_export_billing_account_roles, count.index)}"
  member             = "serviceAccount:${google_service_account.log_export.email}"
}

# roles/logging.configWriter is needed at the organization level to be able to
# test organization level log sinks.
resource "google_organization_iam_member" "log_export" {
  provider = "google.phoogle"

  count = "${length(local.log_export_organization_roles)}"

  org_id = "${var.phoogle_org_id}"
  role   = "${element(local.log_export_organization_roles, count.index)}"
  member = "serviceAccount:${google_service_account.log_export.email}"
}

# There is a test in the log-exports module that needs to spin up a project
# within a folder, and then reference that project within the test. Because
# of that test we need to assign roles/resourcemanager.projectCreator on the
# folder we're using for log-exports
resource "google_folder_iam_member" "log_export" {
  provider = "google.phoogle"

  count = "${length(local.log_export_folder_roles)}"

  folder = "${google_folder.log_export.name}"
  role   = "${element(local.log_export_folder_roles, count.index)}"
  member = "serviceAccount:${google_service_account.log_export.email}"
}

resource "google_service_account_key" "log_export" {
  provider = "google.phoogle"

  service_account_id = "${google_service_account.log_export.id}"
}

resource "random_id" "log_export_github_webhook_token" {
  byte_length = 20
}

data "template_file" "log_export_github_webhook_url" {
  template = "https://concourse.infra.cft.tips/api/v1/teams/cft/pipelines/$${pipeline}/resources/pull-request/check/webhook?webhook_token=$${webhook_token}"

  vars {
    pipeline      = "terraform-google-log-export"
    webhook_token = "${random_id.log_export_github_webhook_token.hex}"
  }
}

resource "kubernetes_secret" "log_export" {
  metadata {
    namespace = "concourse-cft"
    name      = "log-export"
  }

  data {
    github_webhook_token = "${random_id.log_export_github_webhook_token.hex}"
    phoogle_project_id   = "${google_project.log_export.id}"
    phoogle_folder_id    = "${google_folder.log_export.id}"
    phoogle_sa           = "${base64decode(google_service_account_key.log_export.private_key)}"
  }
}
