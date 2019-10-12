/*
Since testing TFV doesn't actually require creating any resources,
there are no reason to create an ephemeral test project + service account each build.
*/

locals {
  terraform_validator_project_name = "ci-terraform-validator"
  terraform_validator_int_required_roles = [
    "roles/editor"
  ]
}

module "terraform_validator_test_project" {
  source  = "terraform-google-modules/project-factory/google"
  version = "~> 3.0"

  name                = local.terraform_validator_project_name
  random_project_id   = false
  org_id              = local.org_id
  folder_id           = local.ci_folders[local.terraform_validator_project_name]
  billing_account     = local.billing_account

  activate_apis = [
    "cloudresourcemanager.googleapis.com"
  ]
}

resource "google_project_iam_member" "int_test" {
  count = length(local.terraform_validator_int_required_roles)

  project = module.terraform_validator_test_project.project_id
  role    = local.terraform_validator_int_required_roles[count.index]
  member  = "serviceAccount:${local.cicd_project_number}@cloudbuild.gserviceaccount.com"
}
