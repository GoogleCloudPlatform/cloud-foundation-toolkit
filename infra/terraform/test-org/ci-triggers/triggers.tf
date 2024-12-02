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

resource "google_cloudbuild_trigger" "int_trigger" {
  provider    = google-beta
  project     = local.project_id
  name        = "${substr(each.key, 0, 50)}-int-trigger"
  description = "Integration tests on pull request for ${each.key}"
  for_each    = local.repo_folder
  github {
    owner = each.value.gh_org
    name  = each.key
    pull_request {
      branch          = ".*"
      comment_control = "COMMENTS_ENABLED_FOR_EXTERNAL_CONTRIBUTORS_ONLY"
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
      _LR_BILLING_ACCOUNT       = local.lr_billing_account
      _TFE_TOKEN_SECRET_ID      = each.key == "terraform-google-tf-cloud-agents" ? google_secret_manager_secret.tfe_token.id : null
      _IM_GITHUB_PAT_SECRET_ID  = each.key == "terraform-google-bootstrap" ? google_secret_manager_secret.im_github_pat.id : null
      _IM_GITLAB_PAT_SECRET_ID  = each.key == "terraform-google-bootstrap" ? google_secret_manager_secret.im_gitlab_pat.id : null
    },
    # add sfb substitutions
    contains(local.bp_on_sfb, each.key) ? local.sfb_substs : {}
  )

  filename      = "build/int.cloudbuild.yaml"
  ignored_files = ["**/*.md", ".gitignore", ".github/**", "**/metadata.yaml", "**/metadata.display.yaml", "assets/**", "infra/assets/**"]
}

# pull_request triggers do not support run trigger, so we have a shadow periodic trigger
resource "google_cloudbuild_trigger" "periodic_int_trigger" {
  provider    = google-beta
  project     = local.project_id
  name        = substr("${each.key}-periodic-int-trigger", 0, 64)
  description = "Periodic integration tests on pull request for ${each.key}"
  for_each    = { for k, v in local.repo_folder : k => v if contains(local.periodic_repos, k) }
  github {
    owner = each.value.gh_org
    name  = each.key
    # this will be invoked via cloud scheduler, hence using a regex that will not match any branch
    push {
      branch = ".^"
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
      _LR_BILLING_ACCOUNT       = local.lr_billing_account
      _PERIODIC                 = true
    },
    # add sfb substitutions
    contains(local.bp_on_sfb, each.key) ? local.sfb_substs : {}
  )

  filename      = "build/int.cloudbuild.yaml"
  ignored_files = ["**/*.md", ".gitignore", ".github/**", "**/metadata.yaml"]
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
      branch          = ".*"
      comment_control = "COMMENTS_ENABLED_FOR_EXTERNAL_CONTRIBUTORS_ONLY"
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
      branch          = ".*"
      comment_control = "COMMENTS_ENABLED_FOR_EXTERNAL_CONTRIBUTORS_ONLY"
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
      branch          = ".*"
      comment_control = "COMMENTS_ENABLED_FOR_EXTERNAL_CONTRIBUTORS_ONLY"
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

resource "google_cloudbuild_trigger" "tgc_main_integration_tests" {
  for_each = {
    tf12 = "0.12.31"
    tf13 = "0.13.7"
  }
  name        = "tgc-main-integration-tests-${each.key}"
  description = "Main/release branch integration tests for terraform-google-conversion with terraform ${each.value}. Managed by Terraform https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/infra/terraform/test-org/tf-validator/project.tf"

  provider = google-beta
  project  = local.project_id
  github {
    owner = "GoogleCloudPlatform"
    name  = "terraform-google-conversion"
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

resource "google_cloudbuild_trigger" "tgc_pull_integration_tests" {
  for_each = {
    tf12 = "0.12.31"
    tf13 = "0.13.7"
  }
  name        = "tgc-pull-integration-tests-${each.key}"
  description = "Pull request integration tests for terraform-google-conversion with terraform ${each.value}. Managed by Terraform https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/infra/terraform/test-org/tf-validator/project.tf"

  provider = google-beta
  project  = local.project_id
  github {
    owner = "GoogleCloudPlatform"
    name  = "terraform-google-conversion"
    pull_request {
      branch          = ".*"
      comment_control = "COMMENTS_ENABLED_FOR_EXTERNAL_CONTRIBUTORS_ONLY"
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

resource "google_cloudbuild_trigger" "tgc_pull_unit_tests" {
  name        = "tgc-pull-unit-tests"
  description = "Pull request unit tests for terraform-google-conversion. Managed by Terraform https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/infra/terraform/test-org/tf-validator/project.tf"

  provider = google-beta
  project  = local.project_id
  github {
    owner = "GoogleCloudPlatform"
    name  = "terraform-google-conversion"
    pull_request {
      branch          = ".*"
      comment_control = "COMMENTS_ENABLED_FOR_EXTERNAL_CONTRIBUTORS_ONLY"
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

resource "google_cloudbuild_trigger" "tgc_main_unit_tests" {
  name        = "tgc-main-unit-tests"
  description = "Main/release branch unit tests for terraform-google-conversion. Managed by Terraform https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/infra/terraform/test-org/tf-validator/project.tf"

  provider = google-beta
  project  = local.project_id
  github {
    owner = "GoogleCloudPlatform"
    name  = "terraform-google-conversion"
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

resource "google_cloudbuild_trigger" "tgc_pull_license_check" {
  name        = "tgc-pull-license-check"
  description = "Pull request license check for terraform-google-conversion. Managed by Terraform https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/infra/terraform/test-org/tf-validator/project.tf"

  provider = google-beta
  project  = local.project_id
  github {
    owner = "GoogleCloudPlatform"
    name  = "terraform-google-conversion"
    pull_request {
      branch          = ".*"
      comment_control = "COMMENTS_ENABLED_FOR_EXTERNAL_CONTRIBUTORS_ONLY"
    }
  }

  filename = ".ci/cloudbuild-tests-go-licenses.yaml"
}

resource "google_cloudbuild_trigger" "tgc_main_license_check" {
  name        = "tgc-main-license-check"
  description = "Main/release branch license check for terraform-google-conversion. Managed by Terraform https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/infra/terraform/test-org/tf-validator/project.tf"

  provider = google-beta
  project  = local.project_id
  github {
    owner = "GoogleCloudPlatform"
    name  = "terraform-google-conversion"
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
      branch          = ".*"
      comment_control = "COMMENTS_ENABLED_FOR_EXTERNAL_CONTRIBUTORS_ONLY"
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
      branch          = ".*"
      comment_control = "COMMENTS_ENABLED_FOR_EXTERNAL_CONTRIBUTORS_ONLY"
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
      branch          = ".*"
      comment_control = "COMMENTS_ENABLED_FOR_EXTERNAL_CONTRIBUTORS_ONLY"
    }
  }
  substitutions = {
    _BILLING_ACCOUNT               = local.billing_account
    _FOLDER_ID                     = values(local.example_foundation)[0]["folder_id"]
    _ORG_ID                        = local.org_id
    _EXAMPLE_FOUNDATIONS_TEST_MODE = each.value
  }

  filename      = "build/int.cloudbuild.yaml"
  ignored_files = ["**/*.md", "**/*.png", ".gitignore", ".github/**", "**/*.example.tfvars", "helpers/foundation-deployer/**"]
}


resource "google_cloudbuild_trigger" "bpt_int_trigger" {
  provider    = google-beta
  project     = local.project_id
  name        = "bpt-int-trigger"
  description = "Integration tests on pull request for blueprint test framework"
  github {
    owner = "GoogleCloudPlatform"
    name  = "cloud-foundation-toolkit"
    pull_request {
      branch          = ".*"
      comment_control = "COMMENTS_ENABLED_FOR_EXTERNAL_CONTRIBUTORS_ONLY"
    }
  }
  substitutions = {
    _BILLING_ACCOUNT = local.billing_account
    _FOLDER_ID       = data.terraform_remote_state.org.outputs.bpt_folder
    _ORG_ID          = local.org_id
  }

  filename       = "infra/blueprint-test/build/int.cloudbuild.yaml"
  included_files = ["infra/blueprint-test/**"]
}
