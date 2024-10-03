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

resource "google_folder" "ci_bq_external_data_folder" {
  display_name = "ci-bq-external-data-folder"
  parent       = "folders/${replace(local.folders["ci-projects"], "folders/", "")}"
}

module "ci_bq_external_data_project" {
  source  = "terraform-google-modules/project-factory/google"
  version = "~> 17.0"

  name            = "ci-bq-external-data-project"
  project_id      = "ci-bq-external-data-project"
  org_id          = local.org_id
  folder_id       = google_folder.ci_bq_external_data_folder.id
  billing_account = local.old_billing_account

  labels = {
    cft-ci = "permanent"
  }

  activate_apis = [
    "storage.googleapis.com",
  ]
}

resource "google_storage_bucket" "ci_bq_external_data_storage_bucket" {
  name     = "ci-bq-external-data"
  project  = module.ci_bq_external_data_project.project_id
  location = "US"
}

resource "google_storage_bucket_iam_member" "ci_bq_external_data_storage_bucket_member" {
  bucket = google_storage_bucket.ci_bq_external_data_storage_bucket.name
  role   = "roles/storage.legacyObjectReader"
  member = "allUsers"
}

resource "google_storage_bucket_object" "ci_bq_external_csv_file" {
  name   = "bigquery-external-table-test.csv"
  source = "external_data/bigquery-external-table-test.csv"
  bucket = google_storage_bucket.ci_bq_external_data_storage_bucket.name
}

resource "google_storage_bucket_object" "ci_bq_external_hive_file_foo" {
  name   = "hive_partition_example/year=2012/foo.csv"
  source = "external_data/foo.csv"
  bucket = google_storage_bucket.ci_bq_external_data_storage_bucket.name
}

resource "google_storage_bucket_object" "ci_bq_external_hive_file_bar" {
  name   = "hive_partition_example/year=2013/bar.csv"
  source = "external_data/bar.csv"
  bucket = google_storage_bucket.ci_bq_external_data_storage_bucket.name
}
