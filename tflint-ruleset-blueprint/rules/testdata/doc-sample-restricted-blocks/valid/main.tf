# https://github.com/terraform-google-modules/terraform-docs-samples/blob/main/bigquery/bigquery_create_view/main.tf

# [START bigquery_create_view]
resource "google_bigquery_dataset" "default" {
  dataset_id                      = "mydataset"
  default_partition_expiration_ms = 2592000000  # 30 days
  default_table_expiration_ms     = 31536000000 # 365 days
  description                     = "dataset description"
  location                        = "US"
  max_time_travel_hours           = 96 # 4 days

  labels = {
    billing_group = "accounting",
    pii           = "sensitive"
  }
}

resource "google_bigquery_table" "default" {
  dataset_id          = google_bigquery_dataset.default.dataset_id
  table_id            = "myview"
  deletion_protection = false # set to "true" in production

  view {
    query          = "SELECT global_id, faa_identifier, name, latitude, longitude FROM `bigquery-public-data.faa.us_airports`"
    use_legacy_sql = false
  }

}
# [END bigquery_create_view]
