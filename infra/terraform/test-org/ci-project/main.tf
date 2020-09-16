provider "github" {
  version      = "~> 2.2"
  organization = local.gh_orgs.infra
}

provider "google" {
  version = "~> 3.39"
  project = local.project_id
}

provider "google-beta" {
  version = "~> 3.39"
  project = local.project_id
}