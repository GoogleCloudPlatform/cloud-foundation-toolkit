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

output "lr_billing_account" {
  value = local.lr_billing_account
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


# output "ci_gsuite_sa_key" {
#   value     = google_service_account_key.ci_gsuite_sa.private_key
#   sensitive = true
# }

# output "ci_gsuite_sa_bucket" {
#   value = google_storage_bucket.ci_gsuite_sa.name
# }

# output "ci_gsuite_sa_bucket_path" {
#   value = google_storage_bucket_object.ci_gsuite_sa_json.name
# }

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

output "prow_int_sa" {
  value = module.prow-int-sa-wi.gcp_service_account_email
}

output "ci_media_cdn_vod_project_id" {
  value = module.ci_media_cdn_vod_project.project_id
}

output "modules" {
  value = [for value in local.repos : value if try(value.module, true)]

  precondition {
    condition     = length(setsubtract(local.invalid_owners, var.temp_allow_invalid_owners)) == 0
    error_message = "Provided Repo Owners are not currently members of GCP or TGM Orgs: ${join(", ", setsubtract(local.invalid_owners, var.temp_allow_invalid_owners))}. You can bypass this error by setting `-var='temp_allow_invalid_owners=[\"${join("\",\"", local.invalid_owners)}\"]'` when running plan/apply."
  }

}

output "bpt_folder" {
  value = module.bpt_ci_folder.id
}

output "periodic_repos" {
  value = sort([for value in local.repos : coalesce(try(value.name, null), try(value.short_name, null)) if try(value.enable_periodic, false)])
}
