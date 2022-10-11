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

resource "google_cloudbuild_trigger" "org_iam_cleaner_trigger" {
  provider = google-beta

  description = "Reset org IAM"
  github {
    owner = local.gh_orgs.infra
    name  = local.gh_repos.infra
    # this will be invoked via cloud scheduler, hence using a regex that will not match any branch
    push {
      branch = ".^"
    }
  }

  filename = "infra/terraform/test-org/org-iam-policy/cloudbuild.yaml"
}

resource "google_service_account" "service_account" {
  account_id   = "org-iam-clean-trigger"
  display_name = "SA used by Cloud Scheduler to trigger cleanup builds"
}

resource "google_project_iam_custom_role" "create_build_role" {
  project     = local.project_id
  role_id     = "CreateBuild"
  title       = "Create Builds"
  description = "Role used by Cloud Scheduler SA for triggering IAM cleanup builds"

  permissions = [
    "cloudbuild.builds.create",
    "cloudbuild.builds.get",
    "cloudbuild.builds.list"
  ]
}

resource "google_project_iam_member" "project" {
  role    = google_project_iam_custom_role.create_build_role.id
  member  = "serviceAccount:${google_service_account.service_account.email}"
  project = local.project_id
}

resource "google_cloud_scheduler_job" "job" {
  name        = "trigger-test-org-iam-reset-build"
  description = "Trigger reset test org IAM build"
  region      = "us-central1"
  # run every day at 3:00
  schedule = "0 3 * * *"

  http_target {
    http_method = "POST"
    uri         = "https://cloudbuild.googleapis.com/v1/projects/${local.project_id}/triggers/${google_cloudbuild_trigger.org_iam_cleaner_trigger.trigger_id}:run"
    body        = base64encode("{\"branchName\": \"master\"}")
    oauth_token {
      service_account_email = google_service_account.service_account.email
    }
  }
  depends_on = [google_project_iam_member.project]
}
