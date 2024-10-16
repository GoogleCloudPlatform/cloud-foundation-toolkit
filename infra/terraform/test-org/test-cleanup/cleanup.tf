/**
 * Copyright 2019-2024 Google LLC
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

module "scheduler-app-engine" {
  source      = "terraform-google-modules/project-factory/google//modules/app_engine"
  version     = "~> 17.0"
  location_id = local.app_location
  project_id  = module.cft-manager-project.project_id
}

module "projects_cleaner" {
  source  = "terraform-google-modules/scheduled-function/google//modules/project_cleanup"
  version = "~> 6.0"

  job_schedule                = "17 * * * *"
  max_project_age_in_hours    = "6"
  organization_id             = local.org_id
  project_id                  = module.cft-manager-project.project_id
  region                      = local.region
  target_excluded_labels      = local.exclude_labels
  target_folder_id            = local.cleanup_folder
  clean_up_org_level_tag_keys = true

  clean_up_org_level_cai_feeds = true
  target_included_feeds        = [".*/feeds/fd-cai-monitoring-.*"]

  clean_up_org_level_scc_notifications = true
  target_included_scc_notifications    = [".*/notificationConfigs/scc-notify-.*"]

  clean_up_billing_sinks   = true
  target_billing_sinks     = [".*/sinks/sk-c-logging-.*-billing-.*"]
  billing_account          = local.billing_account
  function_docker_registry = "ARTIFACT_REGISTRY"
}
