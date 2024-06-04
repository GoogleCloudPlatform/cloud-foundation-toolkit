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
  cft_ci_group          = "cft-ci-robots@test.blueprints.joonix.net"
  cft_dev_group         = "cft-developers@develop.blueprints.joonix.net"
  gcp_admins_group_test = "gcp-admins@test.blueprints.joonix.net"
  project_cleaner       = "project-cleaner-function@${data.terraform_remote_state.project_cleaner.outputs.project_id}.iam.gserviceaccount.com"
  billing_admin_group   = "billing-admin@test.blueprints.joonix.net"

  ci_gsuite_sa           = "ci-gsuite-sa@ci-gsuite-sa-project.iam.gserviceaccount.com"
  cft_admin              = "cft-admin@test.blueprints.joonix.net"
  foundation_leads_group = "cloud-foundation-leads@google.com"

  policy = {
    "roles/billing.admin" : ["group:${local.gcp_admins_group_test}"],
    "roles/compute.xpnAdmin" : ["group:${local.cft_ci_group}"],
    "roles/containeranalysis.admin" : ["group:${local.cft_ci_group}"],
    "roles/orgpolicy.policyAdmin" : ["group:${local.gcp_admins_group_test}"],
    "roles/resourcemanager.folderAdmin" : ["group:${local.gcp_admins_group_test}"],
    "roles/resourcemanager.folderViewer" : ["serviceAccount:${local.project_cleaner}"],
    "roles/resourcemanager.lienModifier" : ["serviceAccount:${local.project_cleaner}"],
    "roles/resourcemanager.organizationAdmin" : [
      "group:${local.cft_ci_group}",
      "group:${local.gcp_admins_group_test}",
      "serviceAccount:${data.google_secret_manager_secret_version.org-admin-sa.secret_data}",
    ],
    "roles/iam.organizationRoleAdmin" : ["serviceAccount:${data.google_secret_manager_secret_version.org-role-admin-sa.secret_data}", ],
    "roles/resourcemanager.organizationViewer" : ["group:${local.cft_ci_group}"],
    "roles/resourcemanager.projectDeleter" : ["serviceAccount:${local.project_cleaner}"],
    "roles/owner" : ["group:${local.gcp_admins_group_test}", "serviceAccount:${local.project_cleaner}"],
    "roles/browser" : ["group:${local.cft_dev_group}"],
    "roles/viewer" : ["group:${local.cft_dev_group}", "serviceAccount:${local.project_cleaner}"],
    "roles/compute.orgSecurityPolicyAdmin" : ["serviceAccount:${local.project_cleaner}"],
    "roles/compute.orgSecurityResourceAdmin" : ["serviceAccount:${local.project_cleaner}"],
    "roles/resourcemanager.folderEditor" : ["serviceAccount:${local.project_cleaner}"],
    "roles/serviceusage.serviceUsageAdmin" : ["serviceAccount:${local.project_cleaner}"],
    "roles/accesscontextmanager.policyReader" : ["group:${local.cft_ci_group}"],
    "roles/assuredworkloads.admin" : ["group:${local.cft_ci_group}"],
    "roles/iam.denyAdmin" : ["group:${local.cft_ci_group}"],
    "roles/resourcemanager.tagAdmin" : ["group:${local.cft_ci_group}"],
  }

  billing_policy = {
    "roles/billing.admin" : [
      "group:${local.cft_ci_group}",
      "group:${local.gcp_admins_group_test}",
      "user:${local.cft_admin}",
      "group:${local.foundation_leads_group}",
      "group:${data.google_secret_manager_secret_version.ba-admin-1.secret_data}",
      "group:${data.google_secret_manager_secret_version.ba-admin-2.secret_data}",
      "group:${local.billing_admin_group}",
    ],
    "roles/logging.configWriter" : [
      "serviceAccount:${local.project_cleaner}",
      "user:${local.cft_admin}",
    ],
    "roles/billing.user" : concat([
      "serviceAccount:${local.ci_gsuite_sa}",
    ], jsondecode(data.google_storage_bucket_object_content.ba-users.content))
  }
}

data "google_secret_manager_secret_version" "org-admin-sa" {
  project = "cloud-foundation-cicd"
  secret  = "org-admin-sa"
}

data "google_secret_manager_secret_version" "org-role-admin-sa" {
  project = "cloud-foundation-cicd"
  secret  = "org-role-admin-sa"
}

data "google_secret_manager_secret_version" "ba-admin-1" {
  project = "cloud-foundation-cicd"
  secret  = "ba-admin-1"
}

data "google_secret_manager_secret_version" "ba-admin-2" {
  project = "cloud-foundation-cicd"
  secret  = "ba-admin-2"
}

data "google_storage_bucket_object_content" "ba-users" {
  name   = "ba-users.json"
  bucket = "tf-data-199f44ed6f9a7f22"
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

resource "google_billing_account_iam_policy" "billing" {
  billing_account_id = data.terraform_remote_state.org.outputs.billing_account
  policy_data        = data.google_iam_policy.billing.policy_data
}

data "google_iam_policy" "billing" {
  dynamic "binding" {
    for_each = local.billing_policy
    content {
      role    = binding.key
      members = binding.value
    }
  }
}
