
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

module "org_bindings" {
  source  = "terraform-google-modules/iam/google"
  version = "~> 2.0"

  organizations = [local.org_id]

  bindings = {
    "roles/resourcemanager.organizationViewer" = [
      "group:${local.cft_ci_group}",
    ]
    "roles/resourcemanager.organizationAdmin" = [
      "group:${local.cft_ci_group}",
    ]
    "roles/compute.xpnAdmin" = [
      "group:${local.cft_ci_group}",
    ]
    "roles/viewer" = [
      "group:${local.cft_dev_group}"
    ]
  }
}

module "admin_bindings" {
  source  = "terraform-google-modules/iam/google"
  version = "~> 2.0"

  folders = [local.folders["ci-projects"]]

  bindings = {
    "roles/resourcemanager.projectCreator" = [
      "group:${local.gcp_admins_group}",
    ]

    "roles/resourcemanager.folderAdmin" = [
      "group:${local.gcp_admins_group}",
    ]

    "roles/billing.projectManager" = [
      "group:${local.gcp_admins_group}",
    ]
  }
}

module "ci_bindings" {
  source  = "terraform-google-modules/iam/google"
  version = "~> 2.0"

  folders = [local.folders["ci-projects"]]

  bindings = {
    "roles/resourcemanager.projectCreator" = [
      "group:${local.cft_ci_group}",
    ]

    "roles/resourcemanager.folderAdmin" = [
      "group:${local.cft_ci_group}",
    ]

    "roles/billing.projectManager" = [
      "group:${local.cft_ci_group}",
    ]
  }
}

module "ci_folders_folder_bindings" {
  source  = "terraform-google-modules/iam/google"
  version = "~> 2.0"

  folders = [local.ci_folders["ci-folders"]]

  bindings = {
    "roles/resourcemanager.folderIamAdmin" = [
      "group:${local.cft_ci_group}",
    ]
  }
}

resource "google_billing_account_iam_member" "ci-billing-user" {
  billing_account_id = local.billing_account
  role               = "roles/billing.admin"
  member             = "group:${local.cft_ci_group}"
}
