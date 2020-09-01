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

module "app-engine" {
  source      = "terraform-google-modules/project-factory/google//modules/app_engine"
  version     = "~> 9.0"
  location_id = local.app_location
  project_id  = module.project.project_id
}

module "projects_cleanup" {
  source = "github.com/terraform-google-modules/terraform-google-scheduled-function//modules/project_cleanup?ref=bucket-name-default"
  # version = "~> 1.5"

  job_schedule             = "17 * * * *"
  max_project_age_in_hours = "6"
  organization_id          = local.org_id
  project_id               = module.project.project_id
  region                   = local.region
  target_excluded_labels   = local.exclude_labels
  target_folder_id         = local.org_id
}

