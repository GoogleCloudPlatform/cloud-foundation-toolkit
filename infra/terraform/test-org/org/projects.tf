/**
 * Copyright 2022 Google LLC
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

module "ci_media_cdn_vod_project" {
  source  = "terraform-google-modules/project-factory/google"
  version = "~> 17.0"

  name            = "ci-media-cdn-vod-project"
  project_id      = "ci-media-cdn-vod-project"
  org_id          = local.org_id
  folder_id       = module.folders-ci.ids["ci-media-cdn-vod"]
  billing_account = local.old_billing_account

  labels = {
    cft-ci = "permanent"
  }

  activate_apis = [
    "cloudresourcemanager.googleapis.com",
    "compute.googleapis.com",
    "iam.googleapis.com",
  ]
}
