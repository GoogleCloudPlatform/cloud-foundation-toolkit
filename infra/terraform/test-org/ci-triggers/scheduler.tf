/**
 * Copyright 2023-2024 Google LLC
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
  periodic_repos = toset([for item in data.terraform_remote_state.org.outputs.periodic_repos : contains(keys(local.custom_repo_mapping), item) ? local.custom_repo_mapping[item] : item])
}

resource "google_service_account" "periodic_sa" {
  project      = local.project_id
  account_id   = "periodic-test-trigger-sa"
  display_name = "SA used by Cloud Scheduler to trigger periodics"
}

resource "google_project_iam_member" "periodic_role" {
  # custom role we created in ci-project/cleaner for scheduled cleaner builds
  project = local.project_id
  role    = "projects/${local.project_id}/roles/CreateBuild"
  member  = "serviceAccount:${google_service_account.periodic_sa.email}"
}

resource "google_cloud_scheduler_job" "job" {
  for_each    = local.periodic_repos
  name        = "periodic-${each.value}"
  project     = local.project_id
  description = "Trigger periodic build for ${each.value}"
  region      = "us-central1"
  # run every day at 3:00
  # todo(bharathkkb): likely need to stagger run times once number of repos increase
  schedule = "0 3 * * *"

  http_target {
    http_method = "POST"
    uri         = "https://cloudbuild.googleapis.com/v1/projects/${local.project_id}/triggers/${google_cloudbuild_trigger.periodic_int_trigger[each.value].trigger_id}:run"
    body        = base64encode("{\"branchName\": \"main\"}")
    oauth_token {
      service_account_email = google_service_account.periodic_sa.email
    }
  }
}
