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
  project_id = "cloud-foundation-cicd"
}

module "gcf_service_account" {
  source     = "terraform-google-modules/service-accounts/google"
  version    = "~> 3.0"
  project_id = local.project_id
  names      = ["pr-comment-cf-sa"]
  project_roles = [
    "${local.project_id}=>roles/cloudbuild.builds.editor"
  ]
}

module "pr_comment_build_function" {
  source                = "terraform-google-modules/event-function/google"
  version               = "~> 1.0"
  name                  = "pr-comment-downstream-builder"
  project_id            = local.project_id
  region                = "us-central1"
  description           = "Launches a downstream build that comments on a PR."
  entry_point           = "main"
  runtime               = "python37"
  source_directory      = "${path.module}/function_source"
  service_account_email = module.gcf_service_account.email
  bucket_force_destroy  = true

  environment_variables = {
    CLOUDBUILD_PROJECT = local.project_id
  }

  event_trigger = {
    event_type = "google.pubsub.topic.publish"
    resource   = "projects/${local.project_id}/topics/cloud-builds"
  }
}
