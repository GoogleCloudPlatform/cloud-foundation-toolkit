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

output "folders" {
  value = merge(local.folders, local.ci_folders)
}

output "ci_repos_folders" {
  value = local.ci_repos_folders
}

output "org_id" {
  value = local.org_id
}

output "billing_account" {
  value = local.billing_account
}

output "cft_ci_group" {
  value = local.cft_ci_group
}

output "ci_gsuite_sa_id" {
  value = google_service_account.ci_gsuite_sa.id
}

output "ci_gsuite_sa_email" {
  value = google_service_account.ci_gsuite_sa.email
}

output "ci_gsuite_sa_folder_id" {
  value = google_folder.ci_gsuite_sa_folder.id
}

output "ci_gsuite_sa_project_id" {
  value = module.ci_gsuite_sa_project.project_id
}

output "ci_gsuite_sa_key" {
  value     = google_service_account_key.ci_gsuite_sa.private_key
  sensitive = true
}

output "ci_gsuite_sa_bucket" {
  value = google_storage_bucket.ci_gsuite_sa.name
}

output "ci_gsuite_sa_bucket_path" {
  value = google_storage_bucket_object.ci_gsuite_sa_json.name
}

output "ci_bq_external_data_folder_id" {
  value = google_folder.ci_bq_external_data_folder.id
}

output "ci_bq_external_data_project_id" {
  value = module.ci_bq_external_data_project.project_id
}

output "ci_bq_external_data_storage_bucket" {
  value = google_storage_bucket.ci_bq_external_data_storage_bucket.name
}

output "ci_bq_external_csv_file" {
  value = google_storage_bucket_object.ci_bq_external_csv_file.name
}

output "ci_bq_external_hive_file_foo" {
  value = google_storage_bucket_object.ci_bq_external_hive_file_foo.name
}

output "ci_bq_external_hive_file_bar" {
  value = google_storage_bucket_object.ci_bq_external_hive_file_bar.name
}
