module "folders-root" {
  source  = "terraform-google-modules/folders/google"
  version = "~> 1.0"

  parent_id   = local.org_id
  parent_type = "organization"

  names = [
    "ci-projects"
  ]

  set_roles = false
}

module "folders-ci" {
  source  = "terraform-google-modules/folders/google"
  version = "~> 1.0"

  parent_id   = replace(module.folders-root.names_and_ids["ci-projects"], "folders/", "")
  parent_type = "folder"

  names = [
    "ci-kms"
  ]

  set_roles = false
}
