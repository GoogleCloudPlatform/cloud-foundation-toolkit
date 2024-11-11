/**
 * Copyright 2021-2024 Google LLC
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
  int_required_roles = [
    "roles/compute.networkAdmin",
    "roles/compute.securityAdmin",
    "roles/iam.serviceAccountUser",
    "roles/vpcaccess.admin",
    "roles/serviceusage.serviceUsageAdmin",
    "roles/container.admin",
    "roles/cloudasset.viewer",
    "roles/serviceusage.serviceUsageConsumer"
  ]
}

module "project" {
  source  = "terraform-google-modules/project-factory/google"
  version = "~> 17.0"

  name              = "ci-bptest"
  random_project_id = "true"
  org_id            = var.org_id
  folder_id         = var.folder_id
  billing_account   = var.billing_account

  default_service_account  = "DEPRIVILEGE"
  deletion_policy          = "DELETE"

  activate_apis = [
    "cloudresourcemanager.googleapis.com",
    "compute.googleapis.com",
    "serviceusage.googleapis.com",
    "vpcaccess.googleapis.com",
    "container.googleapis.com",
    "cloudasset.googleapis.com"
  ]
}

resource "google_service_account" "sa" {
  project      = module.project.project_id
  account_id   = "ci-account"
  display_name = "ci-account"
}

resource "google_project_iam_member" "roles" {
  for_each = toset(local.int_required_roles)

  project = module.project.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_service_account_key" "key" {
  service_account_id = google_service_account.sa.id
}

module "kubernetes-engine_example_simple_autopilot_public" {
  source  = "terraform-google-modules/kubernetes-engine/google//examples/simple_autopilot_public"
  version                     = "~> 34.0"
  project_id                  = module.project.project_id
}
