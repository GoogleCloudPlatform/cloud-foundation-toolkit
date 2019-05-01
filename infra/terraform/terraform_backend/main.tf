resource "google_storage_bucket" "terraform_backend" {
  name     = "cloud-foundation-cicd-tfstate"
  location = "US"
  project  = "${module.variables.project_id}"

  versioning {
    enabled = true
  }
}

resource "google_storage_bucket" "terraform_backend_seed_projects" {
  name     = "cloud-foundation-cicd-seed-projects-tfstate"
  location = "US"
  project  = "${module.variables.project_id}"

  versioning {
    enabled = true
  }
}
