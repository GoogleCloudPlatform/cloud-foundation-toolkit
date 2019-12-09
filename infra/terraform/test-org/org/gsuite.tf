/**
 * Copyright 2019 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

locals {
  ci_gsuite_sa_project_roles = [
    "roles/owner",
    "roles/compute.admin",
    "roles/iam.serviceAccountAdmin",
    "roles/resourcemanager.projectIamAdmin",
    "roles/storage.admin",
    "roles/iam.serviceAccountUser",
    "roles/billing.projectManager",
  ]

  ci_gsuite_sa_folder_roles = [
    "roles/owner",
    "roles/resourcemanager.projectCreator",
    "roles/resourcemanager.folderAdmin",
    "roles/resourcemanager.folderIamAdmin",
    "roles/billing.projectManager",
  ]

  ci_group_gsuite_sa_project_roles = [
    "roles/owner",
    "roles/iam.serviceAccountAdmin",
    "roles/storage.admin",
  ]

  ci_gsuite_sa_bucket      = "ci-gsuite-sa-secrets"
  ci_gsuite_sa_bucket_path = "gsuite-sa.json"
}

resource "google_folder" "ci_gsuite_sa_folder" {
  display_name = "ci-gsuite-sa-folder"
  parent       = "folders/${replace(local.folders["ci-projects"], "folders/", "")}"
}

module "ci_gsuite_sa_project" {
  source  = "terraform-google-modules/project-factory/google"
  version = "~> 4.0"

  name            = "ci-gsuite-sa-project"
  project_id      = "ci-gsuite-sa-project"
  org_id          = local.org_id
  folder_id       = google_folder.ci_gsuite_sa_folder.id
  billing_account = local.billing_account

  labels = {
    cft-ci = "permanent"
  }

  activate_apis = [
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

resource "google_service_account" "ci_gsuite_sa" {
  project      = module.ci_gsuite_sa_project.project_id
  account_id   = "ci-gsuite-sa"
  display_name = "ci-gsuite-sa"
}

resource "google_project_iam_member" "ci_gsuite_sa_project" {
  for_each = toset(local.ci_gsuite_sa_project_roles)

  project = module.ci_gsuite_sa_project.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.ci_gsuite_sa.email}"
}

resource "google_folder_iam_member" "ci_gsuite_sa_folder" {
  for_each = toset(local.ci_gsuite_sa_folder_roles)

  folder = google_folder.ci_gsuite_sa_folder.name
  role   = each.value
  member = "serviceAccount:${google_service_account.ci_gsuite_sa.email}"
}

resource "google_billing_account_iam_member" "ci_gsuite_sa_billing" {
  billing_account_id = local.billing_account
  role               = "roles/billing.user"
  member             = "serviceAccount:${google_service_account.ci_gsuite_sa.email}"
}

// Generate a json key and put it into the secrets bucket.

resource "google_service_account_key" "ci_gsuite_sa" {
  service_account_id = google_service_account.ci_gsuite_sa.id
}

resource "google_storage_bucket" "ci_gsuite_sa" {
  name          = local.ci_gsuite_sa_bucket
  storage_class = "MULTI_REGIONAL"
  project       = module.ci_gsuite_sa_project.project_id

  versioning {
    enabled = true
  }

  force_destroy = true
}

resource "google_storage_bucket_object" "ci_gsuite_sa_json" {
  name    = local.ci_gsuite_sa_bucket_path
  content = base64decode(google_service_account_key.ci_gsuite_sa.private_key)
  bucket  = google_storage_bucket.ci_gsuite_sa.name
}

# Grant G-Suite project rights to cft_ci_group.
# Required to be able to create new gsuite sa keys and to fetch
# the precreated one from the secrets bucket.

resource "google_project_iam_member" "ci_group_gsuite_sa_project" {
  for_each = toset(local.ci_group_gsuite_sa_project_roles)

  project = module.ci_gsuite_sa_project.project_id
  role    = each.value
  member  = "group:${local.cft_ci_group}"
}
