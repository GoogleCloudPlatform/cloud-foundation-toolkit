# https://github.com/terraform-google-modules/terraform-docs-samples/blob/main/bigquery/bigquery_create_view/main.tf

# [START bigquery_create_view]
resource "google_bigquery_dataset" "default" {
  dataset_id                      = var.dataset_id
  default_partition_expiration_ms = var.expr
  default_table_expiration_ms     = 31536000000 # 365 days
  description                     = "dataset description"
  location                        = "US"
  max_time_travel_hours           = 96 # 4 days

  labels = {
    billing_group = "accounting",
    pii           = "sensitive"
  }
}

output "creation_time" {
  value = google_bigquery_dataset.default.creation_time
}

module "bigquery" {
  source  = "terraform-google-modules/bigquery/google"
  version = "~> 7.0"

  dataset_id                 = "foo"
  dataset_name               = "foo"
  description                = "some description"
  project_id                 = var.project_id
  location                   = "US"
  delete_contents_on_destroy = true
  tables = [
    {
      table_id           = "bar",
      time_partitioning  = null,
      range_partitioning = null,
      expiration_time    = 2524604400000, # 2050/01/01
      clustering         = [],
      labels = {
        env      = "devops"
        billable = "true"
        owner    = "joedoe"
      },
    }
  ]
  dataset_labels = {
    env      = "dev"
    billable = "true"
    owner    = "janesmith"
  }
}
# [END bigquery_create_view]
