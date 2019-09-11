
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
      "group:cft-ci-robots@test.infra.cft.tips",
    ]
    "roles/resourcemanager.organizationAdmin" = [
      "group:cft-ci-robots@test.infra.cft.tips",
    ]
    "roles/compute.xpnAdmin" = [
      "group:cft-ci-robots@test.infra.cft.tips",
    ]
  }
}

module "ci_bindings" {
  source  = "terraform-google-modules/iam/google"
  version = "~> 2.0"

  folders = [local.folders["ci"]]

  bindings = {
    "roles/resourcemanager.projectCreator" = [
      "group:cft-ci-robots@test.infra.cft.tips",
    ]

    "roles/resourcemanager.folderAdmin" = [
      "group:cft-ci-robots@test.infra.cft.tips",
    ]

    "roles/billing.projectManager" = [
      "group:cft-ci-robots@test.infra.cft.tips",
    ]
  }
}

module "ci_folders_folder_bindings" {
  source  = "terraform-google-modules/iam/google"
  version = "~> 2.0"

  folders = [local.ci_folders["ci-folders"]]

  bindings = {
    "roles/resourcemanager.folderIamAdmin" = [
      "group:cft-ci-robots@test.infra.cft.tips",
    ]
  }
}

resource "google_billing_account_iam_member" "ci-billing-user" {
  billing_account_id = local.billing_account
  role               = "roles/billing.admin"
  member             = "group:cft-ci-robots@test.infra.cft.tips"
}
