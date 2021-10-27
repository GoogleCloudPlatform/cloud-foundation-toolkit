/**
 * Copyright 2021 Google LLC
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

resource "google_cloudbuild_trigger" "foundation_reconciler" {
  provider = google-beta

  description = "Reset test foundation"
  github {
    owner = local.gh_orgs.infra
    name  = local.gh_repos.infra
    # this will be invoked via cloud scheduler, hence using a regex that will not match any branch
    push {
      branch = ".^"
    }
  }

  filename = "infra/terraform/test-org/ci-foundation/cloudbuild.yaml"
  substitutions = {
    _FOUNDATION_CICD_PROJECT_ID = local.foundation_cicd_project_id
  }
}

resource "google_cloud_scheduler_job" "reconcile_job" {
  name        = "trigger-test-foundation-reconcile"
  description = "Trigger reset test foundation"
  region      = "us-central1"
  # run every everyday
  schedule = "0 23 * * *"

  http_target {
    http_method = "POST"
    uri         = "https://cloudbuild.googleapis.com/v1/projects/${local.project_id}/triggers/${google_cloudbuild_trigger.foundation_reconciler.trigger_id}:run"
    body        = base64encode("{\"branchName\": \"master\"}")
    oauth_token {
      service_account_email = google_service_account.service_account.email
    }
  }
  depends_on = [google_project_iam_member.project]
}
