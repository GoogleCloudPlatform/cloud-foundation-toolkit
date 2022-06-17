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
  terraform_validator_project_name = "ci-terraform-validator"
  terraform_validator_int_required_roles = [
    "roles/editor"
  ]

  magic_modules_cloudbuild_sa = "serviceAccount:673497134629@cloudbuild.gserviceaccount.com"

  projects_folder_id = data.terraform_remote_state.org.outputs.folders["ci-projects"]
  folder_id          = data.terraform_remote_state.org.outputs.folders["ci-terraform-validator"]
  org_id             = data.terraform_remote_state.org.outputs.org_id
  billing_account    = data.terraform_remote_state.org.outputs.billing_account
  cft_ci_group       = data.terraform_remote_state.org.outputs.cft_ci_group
}
