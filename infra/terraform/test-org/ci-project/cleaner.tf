

resource "google_cloudbuild_trigger" "org_iam_cleaner_trigger" {
  provider = google-beta

  description = "Reset org IAM"
  github {
    owner = local.gh_orgs.infra
    name  = local.gh_repos.infra
    # this will be invoked via cloud scheduler, hence using a regex that will not match any branch
    push {
      branch = ".^"
    }
  }

  filename = "infra/terraform/test-org/org-iam-policy/cloudbuild.yaml"
}

resource "google_service_account" "service_account" {
  account_id   = "org-iam-clean-trigger"
  display_name = "SA used by Cloud Scheduler to trigger cleanup builds"
}

resource "google_project_iam_member" "project" {
  role   = "roles/cloudbuild.builds.editor"
  member = "serviceAccount:${google_service_account.service_account.email}"
}

module "app-engine" {
  source      = "terraform-google-modules/project-factory/google//modules/app_engine"
  version     = "~> 9.0"
  location_id = "us-central"
  project_id  = local.project_id
}

resource "google_cloud_scheduler_job" "job" {
  name        = "trigger-test-org-iam-reset-build"
  description = "Trigger reset test org IAM build"
  region      = "us-central1"
  # run every week at 13:00 on Saturday
  schedule = "0 13 * * 6"

  http_target {
    http_method = "POST"
    uri         = "https://cloudbuild.googleapis.com/v1/projects/${local.project_id}/triggers/${google_cloudbuild_trigger.org_iam_cleaner_trigger.trigger_id}:run"
    body        = base64encode("{\"branchName\": \"master\"}")
    oauth_token {
      service_account_email = google_service_account.service_account.email
    }
  }
  depends_on = [google_project_iam_member.project, module.app-engine]
}
