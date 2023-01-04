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
  name        = "${each.key}-lint-trigger"
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
  name        = "${each.key}-int-trigger"
  description = "Integration tests on pull request for ${each.key}"
  for_each    = local.repo_folder
  github {
    owner = each.value.gh_org
    name  = each.key
    pull_request {
      branch = ".*"
    }
  }
  substitutions = merge(
    {
      _BILLING_ACCOUNT          = local.billing_account
      _FOLDER_ID                = each.value.folder_id
      _ORG_ID                   = local.org_id
      _BILLING_IAM_TEST_ACCOUNT = each.key == "terraform-google-iam" ? local.billing_iam_test_account : null
      _VOD_TEST_PROJECT_ID      = each.key == "terraform-google-media-cdn-vod" ? local.vod_test_project_id : null
      _FILE_LOGS_BUCKET         = lookup(local.enable_file_log, each.key, false) ? module.filelogs_bucket.url : null
    },
    # add sfb substitutions
    contains(local.bp_on_sfb, each.key) ? local.sfb_substs : {}
  )

  filename      = "build/int.cloudbuild.yaml"
  ignored_files = ["**/*.md", ".gitignore", ".github/**"]
}

resource "google_cloudbuild_trigger" "tf_validator_main_integration_tests" {
  for_each = {
    tf12 = "0.12.31"
    tf13 = "0.13.7"
  }
  name        = "tf-validator-main-integration-tests-${each.key}"
  description = "Main/release branch integration tests for terraform-validator with terraform ${each.value}. Managed by Terraform https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/infra/terraform/test-org/tf-validator/project.tf"

  provider = google-beta
  project  = local.project_id
  github {
    owner = "GoogleCloudPlatform"
    name  = "terraform-validator"
    push {
      branch = "^(main|release-.+)$"
    }
  }
  substitutions = {
    _TERRAFORM_VERSION = each.value
    _TEST_PROJECT      = local.tf_validator_project_id
    _TEST_FOLDER       = local.tf_validator_folder_id
    _TEST_ANCESTRY     = local.tf_validator_ancestry
    _TEST_ORG          = local.org_id
  }

  filename = ".ci/cloudbuild-tests-integration.yaml"
}

resource "google_cloudbuild_trigger" "tf_validator_pull_integration_tests" {
  for_each = {
    tf12 = "0.12.31"
    tf13 = "0.13.7"
  }
  name        = "tf-validator-pull-integration-tests-${each.key}"
  description = "Pull request integration tests for terraform-validator with terraform ${each.value}. Managed by Terraform https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/infra/terraform/test-org/tf-validator/project.tf"

  provider = google-beta
  project  = local.project_id
  github {
    owner = "GoogleCloudPlatform"
    name  = "terraform-validator"
    pull_request {
      branch = ".*"
    }
  }
  substitutions = {
    _TERRAFORM_VERSION = each.value
    _TEST_PROJECT      = local.tf_validator_project_id
    _TEST_FOLDER       = local.tf_validator_folder_id
    _TEST_ANCESTRY     = local.tf_validator_ancestry
    _TEST_ORG          = local.org_id
  }

  filename = ".ci/cloudbuild-tests-integration.yaml"
}

resource "google_cloudbuild_trigger" "tf_validator_pull_unit_tests" {
  name        = "tf-validator-pull-unit-tests"
  description = "Pull request unit tests for terraform-validator. Managed by Terraform https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/infra/terraform/test-org/tf-validator/project.tf"

  provider = google-beta
  project  = local.project_id
  github {
    owner = "GoogleCloudPlatform"
    name  = "terraform-validator"
    pull_request {
      branch = ".*"
    }
  }
  substitutions = {
    _TEST_PROJECT  = local.tf_validator_project_id
    _TEST_FOLDER   = local.tf_validator_folder_id
    _TEST_ANCESTRY = local.tf_validator_ancestry
    _TEST_ORG      = local.org_id
  }

  filename = ".ci/cloudbuild-tests-unit.yaml"
}

resource "google_cloudbuild_trigger" "tf_validator_main_unit_tests" {
  name        = "tf-validator-main-unit-tests"
  description = "Main/release branch unit tests for terraform-validator. Managed by Terraform https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/infra/terraform/test-org/tf-validator/project.tf"

  provider = google-beta
  project  = local.project_id
  github {
    owner = "GoogleCloudPlatform"
    name  = "terraform-validator"
    push {
      branch = "^(main|release-.+)$"
    }
  }
  substitutions = {
    _TEST_PROJECT  = local.tf_validator_project_id
    _TEST_FOLDER   = local.tf_validator_folder_id
    _TEST_ANCESTRY = local.tf_validator_ancestry
    _TEST_ORG      = local.org_id
  }

  filename = ".ci/cloudbuild-tests-unit.yaml"
}

resource "google_cloudbuild_trigger" "tf_validator_pull_license_check" {
  name        = "tf-validator-pull-license-check"
  description = "Pull request license check for terraform-validator. Managed by Terraform https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/infra/terraform/test-org/tf-validator/project.tf"

  provider = google-beta
  project  = local.project_id
  github {
    owner = "GoogleCloudPlatform"
    name  = "terraform-validator"
    pull_request {
      branch = ".*"
    }
  }

  filename = ".ci/cloudbuild-tests-go-licenses.yaml"
}

resource "google_cloudbuild_trigger" "tf_validator_main_license_check" {
  name        = "tf-validator-main-license-check"
  description = "Main/release branch license check for terraform-validator. Managed by Terraform https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/infra/terraform/test-org/tf-validator/project.tf"

  provider = google-beta
  project  = local.project_id
  github {
    owner = "GoogleCloudPlatform"
    name  = "terraform-validator"
    push {
      branch = "^(main|release-.+)$"
    }
  }

  filename = ".ci/cloudbuild-tests-go-licenses.yaml"
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

# example-foundation-int tests
resource "google_cloudbuild_trigger" "example_foundations_int_trigger" {
  provider    = google-beta
  project     = local.project_id
  name        = "terraform-example-foundation-int-trigger-${each.value}"
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
  ignored_files = ["**/*.md", ".gitignore", ".github/**"]
}
