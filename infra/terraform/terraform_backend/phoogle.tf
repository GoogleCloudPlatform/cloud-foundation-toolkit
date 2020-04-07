# Legacy Terraform state for deleted Phoogle projects
module "phoogle-backend" {
  source  = "terraform-google-modules/cloud-storage/google//modules/simple_bucket"
  version = "~> 1.3"

  name       = "cloud-foundation-cicd-tfstate"
  project_id = module.variables.project_id
  location   = "US"
}

module "phoogle-seed" {
  source  = "terraform-google-modules/cloud-storage/google//modules/simple_bucket"
  version = "~> 1.3"

  name       = "cloud-foundation-cicd-seed-projects-tfstate"
  project_id = module.variables.project_id
  location   = "US"
}
