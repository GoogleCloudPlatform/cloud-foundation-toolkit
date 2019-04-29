terraform {
  required_version = "0.11.13"

  backend "gcs" {
    bucket = "cloud-foundation-cicd-seed-projects-tfstate"
    prefix = "state/"
  }
}

provider "google" {
  version = "~> 2.1"
}

provider "google-beta" {
  version = "~> 2.1"
}

resource "google_folder" "seed_root_folder" {
  parent       = "organizations/${var.org_id}"
  display_name = "seed-projects"
}

resource "google_folder_iam_member" "users_seed_root_members" {
  folder = "${google_folder.seed_root_folder.name}"
  role   = "roles/viewer"
  member = "group:cloud-foundation-devs@phoogle.net"
}
