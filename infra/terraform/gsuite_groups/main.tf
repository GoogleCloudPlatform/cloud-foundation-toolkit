locals {
  // Note - Terraform v0.11 list handling means that when we change elements in the middle
  // of the list we'll destroy/create all following entries. Changing this might mean a
  // brief downtime.
  cloud-foundation-devs = [
    "${var.core_group}"
  ]

  cloud-foundation-devs-roles = [
    "roles/billing.user",
    "roles/resourcemanager.folderAdmin",
    "roles/resourcemanager.organizationAdmin",
    "roles/resourcemanager.projectCreator",
  ]
}

provider "gsuite" {
  impersonated_user_email = "${var.impersonated_user_email}"
  credentials = "${file(var.credentials_path)}"
  oauth_scopes = [
    "https://www.googleapis.com/auth/admin.directory.group",
    "https://www.googleapis.com/auth/admin.directory.group.member",
  ]

  version = "~> 0.1.9"
}

provider "google" {
  version = "~> 1.20"
}

resource "gsuite_group" "cloud-foundation-devs" {
  email = "cloud-foundation-devs@${var.domain}"
  name  = "Cloud Foundation Developers"
}

resource "gsuite_group_member" "cloud-foundation-devs" {
  count = "${length(local.cloud-foundation-devs)}"
  group = "${gsuite_group.cloud-foundation-devs.email}"
  email = "${element(local.cloud-foundation-devs, count.index)}"
  role  = "MEMBER"
}

resource "google_organization_iam_member" "cloud-foundation-devs" {
  count  = "${length(local.cloud-foundation-devs-roles)}"
  org_id = "${var.org_id}"
  member = "group:${gsuite_group.cloud-foundation-devs.email}"
  role   = "${element(local.cloud-foundation-devs-roles, count.index)}"
}
