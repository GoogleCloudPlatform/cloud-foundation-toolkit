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
  ]
  # custom mapping of the form name => repo_name used for overriding `terraform-google` prefix
  custom_repo_mapping = {
    "cloud-foundation-training" = "cloud-foundation-training",
    "example-foundation-app"    = "terraform-example-foundation-app",
    "anthos-samples"            = "anthos-samples"
  }
  # example foundation has custom test modes
  example_foundation                = { "terraform-example-foundation" = data.terraform_remote_state.org.outputs.ci_repos_folders["example-foundation"] }
  example_foundation_int_test_modes = ["default", "HubAndSpoke"]

  repo_folder              = { for key, value in data.terraform_remote_state.org.outputs.ci_repos_folders : contains(keys(local.custom_repo_mapping), key) ? local.custom_repo_mapping[key] : "terraform-google-${key}" => value if !contains(local.exclude_folders, value.folder_name) }
  org_id                   = data.terraform_remote_state.org.outputs.org_id
  billing_account          = data.terraform_remote_state.org.outputs.billing_account
  tf_validator_project_id  = data.terraform_remote_state.tf-validator.outputs.project_id
  project_id               = "cloud-foundation-cicd"
  forseti_ci_folder_id     = "542927601143"
  billing_iam_test_account = "0151A3-65855E-5913CF"
}
