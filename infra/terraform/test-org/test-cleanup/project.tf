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

module "cft-manager-project" {
  source  = "terraform-google-modules/project-factory/google"
  version = "~> 17.0"

  name              = "cft-project-manager"
  random_project_id = true
  org_id            = local.org_id
  folder_id         = local.folder_id
  billing_account   = local.billing_account
  labels            = local.exclude_labels

  activate_apis = [
    "cloudresourcemanager.googleapis.com",
    "cloudscheduler.googleapis.com",
    "cloudbuild.googleapis.com",
    "cloudfunctions.googleapis.com",
    "pubsub.googleapis.com",
    "storage-api.googleapis.com",
    "serviceusage.googleapis.com",
    "storage-component.googleapis.com",
    "appengine.googleapis.com",
    "securitycenter.googleapis.com",
    "cloudasset.googleapis.com"
  ]
}
