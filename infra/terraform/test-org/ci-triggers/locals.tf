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
  exclude_folders = [
    "ci",
    "ci-terraform-validator",
    "ci-projects",
    "ci-shared",
    "ci-anthos-platform",
    "ci-example-foundation",
    "ci-gcp-example-foundation-app"
  ]
  gcp_org_prefix                    = "gcp-"
  example_foundation                = { "terraform-example-foundation" = { folder_id = replace(data.terraform_remote_state.org.outputs.folders["ci-example-foundation"], "folders/", ""), gh_org = "terraform-google-modules" } }
  example_foundation_app            = { "terraform-example-foundation-app" = { folder_id = replace(data.terraform_remote_state.org.outputs.folders["ci-gcp-example-foundation-app"], "folders/", ""), gh_org = "GoogleCloudPlatform" } }
  example_foundation_int_test_modes = ["default", "HubAndSpoke"]
  training_repos = [
    "ci-cloud-foundation-training"
  ]
  repo_folder              = { for key, value in data.terraform_remote_state.org.outputs.folders : contains(local.training_repos, key) ? replace(key, "ci-", "") : replace(key, "ci-", "terraform-google-") => { "folder_id" = replace(value, "folders/", ""), "gh_org" = "terraform-google-modules" } if ! contains(local.exclude_folders, key) && ! (replace(key, local.gcp_org_prefix, "") != key) }
  repo_folder_gcp_org      = merge({ for key, value in data.terraform_remote_state.org.outputs.folders : replace(key, "ci-gcp-", "terraform-google-") => { "folder_id" = replace(value, "folders/", ""), "gh_org" = "GoogleCloudPlatform" } if ! contains(local.exclude_folders, key) && replace(key, local.gcp_org_prefix, "") != key }, local.example_foundation_app)
  org_id                   = data.terraform_remote_state.org.outputs.org_id
  billing_account          = data.terraform_remote_state.org.outputs.billing_account
  tf_validator_project_id  = data.terraform_remote_state.tf-validator.outputs.project_id
  project_id               = "cloud-foundation-cicd"
  forseti_ci_folder_id     = "542927601143"
  billing_iam_test_account = "0151A3-65855E-5913CF"
}
