module "artifact_registry" {
  source  = "GoogleCloudPlatform/artifact-registry/google"
  version = "~> 0.3"

  location      = "us"
  project_id    = local.project_id
  format        = "DOCKER"
  repository_id = "gcr.io"

  members = {
    readers = [
      "allUsers",
    ]

    writers = [
      "serviceAccount:${module.service_accounts.service_account.email}",
    ]
  }
}
