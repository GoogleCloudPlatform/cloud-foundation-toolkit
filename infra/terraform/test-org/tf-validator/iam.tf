
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

resource "google_project_iam_member" "int_test" {
  count = length(local.terraform_validator_int_required_roles)

  project = module.terraform_validator_test_project.project_id
  role    = local.terraform_validator_int_required_roles[count.index]
  member  = "group:${local.cft_ci_group}"
}

resource "google_project_iam_member" "kokoro_test_0" {
  project = module.terraform_validator_test_project.project_id
  role    = "roles/editor"
  member  = "serviceAccount:kokoro-build@magic-modules.iam.gserviceaccount.com"
}

resource "google_project_iam_member" "kokoro_test_1" {
  project = module.terraform_validator_test_project.project_id
  role    = "roles/editor"
  member  = "serviceAccount:kokoro-trampoline@cloud-devrel-kokoro-resources.iam.gserviceaccount.com"
}

resource "google_folder_iam_member" "magic_modules_cloudbuild_sa_folder_viewer" {
  folder = local.projects_folder_id
  role   = "roles/resourcemanager.folderViewer"
  member = local.magic_modules_cloudbuild_sa
}

resource "google_folder_iam_member" "magic_modules_cloudbuild_sa_security_reviewer" {
  folder = local.projects_folder_id
  role   = "roles/iam.securityReviewer"
  member = local.magic_modules_cloudbuild_sa
}
