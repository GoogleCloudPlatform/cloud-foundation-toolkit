# https://github.com/terraform-google-modules/terraform-docs-samples/blob/main/bigquery/bigquery_create_view/main.tf

variable "dataset_id" {

}

# [START bigquery_create_view]
resource "google_bigquery_dataset" "default" {
  dataset_id                      = var.dataset_id
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
# [END bigquery_create_view]
