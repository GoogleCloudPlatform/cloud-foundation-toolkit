/**
 * Copyright 2023 Google LLC
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

resource "google_secret_manager_secret" "tfe_token" {
  project   = local.project_id
  secret_id = "tke-token"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_iam_member" "tfe_token_member" {
  project   = google_secret_manager_secret.tfe_token.project
  secret_id = google_secret_manager_secret.tfe_token.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "group:${data.terraform_remote_state.org.outputs.cft_ci_group}"
}

resource "google_secret_manager_secret" "im_github_pat" {
  project   = local.project_id
  secret_id = "im_github_pat"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_iam_member" "im_github_pat_member" {
  project   = google_secret_manager_secret.im_github_pat.project
  secret_id = google_secret_manager_secret.im_github_pat.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "group:${data.terraform_remote_state.org.outputs.cft_ci_group}"
}

resource "google_secret_manager_secret" "im_gitlab_pat" {
  project   = local.project_id
  secret_id = "im_gitlab_pat"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_iam_member" "im_gitlab_pat_member" {
  project   = google_secret_manager_secret.im_gitlab_pat.project
  secret_id = google_secret_manager_secret.im_gitlab_pat.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "group:${data.terraform_remote_state.org.outputs.cft_ci_group}"
}
