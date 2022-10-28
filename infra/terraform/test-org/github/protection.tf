/**
 * Copyright 2022 Google LLC
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

locals {
  remove_special_repos = [
    # Exclude from CI/branch protection
    "anthos-samples",
    "terraform-docs-samples",
    # Special CI/branch protection case
    "terraform-example-foundation"
  ]
  filtered_repos = setsubtract(toset(local.repos), local.remove_special_repos)
}

module "branch_protection_tgm" {
  source    = "../../modules/branch_protection"
  org       = "terraform-google-modules"
  repo_list = local.filtered_repos
  admin     = "cft-admins"
}

module "branch_protection_gcp" {
  source    = "../../modules/branch_protection"
  org       = "GoogleCloudPlatform"
  repo_list = local.filtered_repos
  admin     = "blueprint-solutions"
}

module "renovate_json_tgm" {
  source    = "../../modules/repo_file"
  org       = "terraform-google-modules"
  repo_list = local.filtered_repos
  filename  = ".github/renovate.json"
  content   = file("${path.module}/renovate.json")
}

module "renovate_json_gcp" {
  source    = "../../modules/repo_file"
  org       = "GoogleCloudPlatform"
  repo_list = local.filtered_repos
  filename  = ".github/renovate.json"
  content   = file("${path.module}/renovate.json")
}

# Special CI/branch protection case

data "github_team" "cft-admins" {
  slug = "cft-admins"
}

data "github_repository" "terraform-example-foundation" {
  name = "terraform-example-foundation"
}

resource "github_branch_protection" "terraform-example-foundation" {
  repository_id = data.github_repository.terraform-example-foundation.node_id
  pattern       = data.github_repository.terraform-example-foundation.default_branch

  required_pull_request_reviews {
    required_approving_review_count = 1
    require_code_owner_reviews      = true
  }

  required_status_checks {
    strict = true
    contexts = [
      "cla/google",
      "terraform-example-foundation-int-trigger-default (cloud-foundation-cicd)",
      "terraform-example-foundation-int-trigger-HubAndSpoke (cloud-foundation-cicd)",
      "terraform-example-foundation-lint-trigger (cloud-foundation-cicd)"
    ]
  }

  enforce_admins = false

  push_restrictions = [
    data.github_team.cft-admins.node_id
  ]

}
