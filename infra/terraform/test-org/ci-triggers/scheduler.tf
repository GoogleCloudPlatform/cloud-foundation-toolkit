locals {
  periodic_repos = toset([for m in data.terraform_remote_state.org.outputs.modules : m.name if lookup(m, "enable_periodic", false)])
}

resource "google_service_account" "periodic_sa" {
  project      = local.project_id
  account_id   = "periodic-test-trigger-sa"
  display_name = "SA used by Cloud Scheduler to trigger periodics"
}

resource "google_project_iam_member" "periodic_role" {
  # custom role we created in ci-project/cleaner for scheduled cleaner builds
  project = local.project_id
  role    = "projects/${local.project_id}/roles/CreateBuild"
  member  = "serviceAccount:${google_service_account.periodic_sa.email}"
}

resource "google_cloud_scheduler_job" "job" {
  for_each    = local.periodic_repos
  name        = "periodic-${each.value}"
  project     = local.project_id
  description = "Trigger periodic build for ${each.value}"
  region      = "us-central1"
  # run every day at 3:00
  # todo(bharathkkb): likely need to stagger run times once number of repos increase
  schedule = "0 3 * * *"

  http_target {
    http_method = "POST"
    uri         = "https://cloudbuild.googleapis.com/v1/projects/${local.project_id}/triggers/${google_cloudbuild_trigger.periodic_int_trigger[each.value].trigger_id}:run"
    body        = base64encode("{\"branchName\": \"main\"}")
    oauth_token {
      service_account_email = google_service_account.periodic_sa.email
    }
  }
}
