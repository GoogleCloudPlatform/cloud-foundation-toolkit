resource "google_storage_bucket" "terraform_backend" {
  name     = "cloud-foundation-cicd-tfstate"
  location = "US"
  project  = "${module.variables.project_id}"

  versioning { enabled = true }
}
