locals {
  seed_service_account_required_org_roles = [
    "roles/resourcemanager.organizationViewer",
    "roles/resourcemanager.projectCreator",
    "roles/billing.user",
    "roles/compute.xpnAdmin",
    "roles/compute.networkAdmin",
  ]

  seed_service_account_required_folder_roles = [
    "roles/resourcemanager.folderViewer",
  ]

  seed_service_account_required_shared_vpc_roles = [
    "roles/browser",
    "roles/resourcemanager.projectIamAdmin",
  ]

  seed_service_account_required_bucket_project_roles = [
    "roles/storage.admin",
  ]

  seed_service_account_required_billing_account_roles = [
    "roles/billing.user",
  ]
}

module "project_factory" {
  source  = "terraform-google-modules/project-factory/google"
  version = "~> 14.0"

  random_project_id           = "true"
  name                        = "${var.username}-seed"
  org_id                      = "${var.org_id}"
  billing_account             = "${var.billing_account}"
  activate_apis               = "${var.seed_project_services}"
  folder_id                   = "${var.seed_folder_id}"
  disable_services_on_destroy = "false"
}

resource "google_compute_shared_vpc_host_project" "main" {
  project = "${module.project_factory.project_id}"
}

resource "google_folder" "users_seed_root" {
  display_name = "${var.username}"
  parent       = "${var.seed_folder_id}"
}

// This account should be used for provisioning test projects. Provisioning
// resources within those test projects should be done using service accounts
// associated with those projects to ensure that required roles are properly
// isolated.

resource "google_service_account" "seed_service_account" {
  project      = "${module.project_factory.project_id}"
  account_id   = "${var.username}-seed"
  display_name = "Project Factory seed service account"
}

resource "google_folder_iam_member" "seed_service_account_folder_roles" {
  count  = "${length(local.seed_service_account_required_folder_roles)}"
  folder = "${google_folder.users_seed_root.name}"
  role   = "${element(local.seed_service_account_required_folder_roles, count.index)}"
  member = "serviceAccount:${google_service_account.seed_service_account.email}"
}

resource "google_organization_iam_member" "seed_service_account_organization_roles" {
  count  = "${length(local.seed_service_account_required_org_roles)}"
  org_id = "${var.org_id}"
  role   = "${element(local.seed_service_account_required_org_roles, count.index)}"
  member = "serviceAccount:${google_service_account.seed_service_account.email}"
}

resource "google_project_iam_member" "seed_service_account_shared_vpc_roles" {
  count   = "${length(local.seed_service_account_required_shared_vpc_roles)}"
  project = "${module.project_factory.project_id}"
  role    = "${element(local.seed_service_account_required_shared_vpc_roles, count.index)}"
  member  = "serviceAccount:${google_service_account.seed_service_account.email}"
}

resource "google_project_iam_member" "seed_service_account_bucket_project_roles" {
  count   = "${length(local.seed_service_account_required_bucket_project_roles)}"
  project = "${module.project_factory.project_id}"
  role    = "${element(local.seed_service_account_required_bucket_project_roles, count.index)}"
  member  = "serviceAccount:${google_service_account.seed_service_account.email}"
}

resource "google_billing_account_iam_member" "seed_service_account_billing_account_roles" {
  count              = "${length(local.seed_service_account_required_billing_account_roles)}"
  billing_account_id = "${var.billing_account}"
  role               = "${element(local.seed_service_account_required_billing_account_roles, count.index)}"
  member             = "serviceAccount:${google_service_account.seed_service_account.email}"
}

resource "google_folder" "users_pf_test_projects" {
  display_name = "pf-test-projects"
  parent       = "${google_folder.users_seed_root.id}"
}

resource "google_project_iam_member" "project_owners_roles" {
  count   = "${length(var.owner_emails)}"
  project = "${module.project_factory.project_id}"
  role    = "roles/owner"
  member  = "${var.owner_emails[count.index]}"
}
