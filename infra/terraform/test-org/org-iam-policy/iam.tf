/**
 * Copyright 2020 Google LLC
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
  cft_ci_group          = "cft-ci-robots@test.infra.cft.tips"
  cft_dev_group         = "cft-developers@dev.infra.cft.tips"
  gcp_admins_group_test = "gcp-admins@test.infra.cft.tips"
  project_cleaner       = "project-cleaner-function@${data.terraform_remote_state.project_cleaner.outputs.project_id}.iam.gserviceaccount.com"
  policy = {
    "roles/billing.admin" : ["group:${local.gcp_admins_group_test}"],
    "roles/compute.xpnAdmin" : ["group:${local.cft_ci_group}"],
    "roles/containeranalysis.admin" : ["group:${local.cft_ci_group}"],
    "roles/orgpolicy.policyAdmin" : ["group:${local.gcp_admins_group_test}"],
    "roles/resourcemanager.folderAdmin" : ["group:${local.gcp_admins_group_test}"],
    "roles/resourcemanager.folderViewer" : ["serviceAccount:${local.project_cleaner}"],
    "roles/resourcemanager.lienModifier" : ["serviceAccount:${local.project_cleaner}"],
    "roles/resourcemanager.organizationAdmin" : ["group:${local.cft_ci_group}", "group:${local.gcp_admins_group_test}", ],
    "roles/resourcemanager.organizationViewer" : ["group:${local.cft_ci_group}"],
    "roles/resourcemanager.projectDeleter" : ["serviceAccount:${local.project_cleaner}"],
    "roles/owner" : ["group:${local.gcp_admins_group_test}", "serviceAccount:${local.project_cleaner}"],
    "roles/browser" : ["group:${local.cft_dev_group}"],
    "roles/viewer" : ["group:${local.cft_dev_group}"]
  }
}

resource "google_organization_iam_policy" "organization" {
  org_id      = data.terraform_remote_state.org.outputs.org_id
  policy_data = data.google_iam_policy.admin.policy_data
}

data "google_iam_policy" "admin" {
  dynamic "binding" {
    for_each = local.policy
    content {
      role    = binding.key
      members = binding.value
    }
  }
}
