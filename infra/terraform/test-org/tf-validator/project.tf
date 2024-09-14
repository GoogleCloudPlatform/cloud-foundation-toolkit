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


// Since testing TFV doesn't actually require creating any resources,
// there are no reason to create an ephemeral test project + service account each build.
module "terraform_validator_test_project" {
  source  = "terraform-google-modules/project-factory/google"
  version = "~> 17.0"

  name              = local.terraform_validator_project_name
  random_project_id = true
  org_id            = local.org_id
  folder_id         = local.folder_id
  billing_account   = local.billing_account
  labels            = data.terraform_remote_state.cleaner.outputs.excluded_labels

  activate_apis = [
    "cloudresourcemanager.googleapis.com"
  ]
}
