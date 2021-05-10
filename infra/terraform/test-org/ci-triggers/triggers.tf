/**
 * Copyright 2019 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

resource "google_cloudbuild_trigger" "lint_trigger" {
  provider    = google-beta
  project     = local.project_id
  description = "Lint tests on pull request for ${each.key}"
  for_each    = merge(local.repo_folder, local.example_foundation)
  github {
    owner = each.value.gh_org
    name  = each.key
    pull_request {
      branch = ".*"
    }
  }

  filename = "build/lint.cloudbuild.yaml"
}

resource "google_cloudbuild_trigger" "int_trigger" {
  provider    = google-beta
  project     = local.project_id
  description = "Integration tests on pull request for ${each.key}"
  for_each    = local.repo_folder
  github {
    owner = each.value.gh_org
    name  = each.key
    pull_request {
      branch = ".*"
    }
  }
  substitutions = {
    _BILLING_ACCOUNT          = local.billing_account
    _FOLDER_ID                = each.value.folder_id
    _ORG_ID                   = local.org_id
    _BILLING_IAM_TEST_ACCOUNT = each.key == "terraform-google-iam" ? local.billing_iam_test_account : null
  }

  filename      = "build/int.cloudbuild.yaml"
  ignored_files = ["*.md", ".gitignore"]
}

resource "google_cloudbuild_trigger" "tf_validator" {
  provider    = google-beta
  project     = local.project_id
  description = "Pull request build for tf-validator"
  github {
    owner = "GoogleCloudPlatform"
    name  = "terraform-validator"
    pull_request {
      branch = ".*"
    }
  }
  substitutions = {
    _TEST_PROJECT = local.tf_validator_project_id
  }

  filename = "build/int.cloudbuild.yaml"
}

resource "google_cloudbuild_trigger" "forseti_lint" {
  provider    = google-beta
  project     = local.project_id
  description = "Lint tests on pull request for forseti"
  github {
    owner = "forseti-security"
    name  = "terraform-google-forseti"
    pull_request {
      branch = ".*"
    }
  }

  filename = "build/lint.cloudbuild.yaml"
}

resource "google_cloudbuild_trigger" "forseti_int" {
  provider    = google-beta
  project     = local.project_id
  description = "Integration tests on pull request for forseti"
  github {
    owner = "forseti-security"
    name  = "terraform-google-forseti"
    pull_request {
      branch = ".*"
    }
  }
  substitutions = {
    _BILLING_ACCOUNT = local.billing_account
    _FOLDER_ID       = local.forseti_ci_folder_id
    _ORG_ID          = local.org_id
  }

  filename = "build/int.cloudbuild.yaml"
}

resource "google_cloudbuild_trigger" "tf_py_test_helper_lint" {
  provider    = google-beta
  project     = local.project_id
  description = "Lint tests on pull request for terraform-python-testing-helper"
  github {
    owner = "GoogleCloudPlatform"
    name  = "terraform-python-testing-helper"
    pull_request {
      branch = ".*"
    }
  }

  filename = ".ci/cloudbuild.lint.yaml"
}

resource "google_cloudbuild_trigger" "tf_py_test_helper_test" {
  provider    = google-beta
  project     = local.project_id
  description = "Test on pull request for terraform-python-testing-helper"
  github {
    owner = "GoogleCloudPlatform"
    name  = "terraform-python-testing-helper"
    pull_request {
      branch = ".*"
    }
  }

  filename = ".ci/cloudbuild.test.yaml"
  included_files = [
    "**/*.tf",
    "**/*.py"
  ]
}

# example-foundation-int tests
resource "google_cloudbuild_trigger" "example_foundations_int_trigger" {
  provider    = google-beta
  project     = local.project_id
  description = "Integration tests on pull request for example_foundations in ${each.value} mode"
  for_each    = toset(local.example_foundation_int_test_modes)
  github {
    owner = values(local.example_foundation)[0]["gh_org"]
    name  = keys(local.example_foundation)[0]
    pull_request {
      branch = ".*"
    }
  }
  substitutions = {
    _BILLING_ACCOUNT               = local.billing_account
    _FOLDER_ID                     = values(local.example_foundation)[0]["folder_id"]
    _ORG_ID                        = local.org_id
    _EXAMPLE_FOUNDATIONS_TEST_MODE = each.value
  }

  filename      = "build/int.cloudbuild.yaml"
  ignored_files = ["*.md", ".gitignore"]
}
