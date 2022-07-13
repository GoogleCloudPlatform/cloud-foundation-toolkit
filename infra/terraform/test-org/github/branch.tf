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
  remove_branch_repos = [
    # TODO: add support for non terraform-google-modules org repos
    "anthos-samples",
    "terraform-example-foundation-app",
    "terraform-google-network-forensics",
    "terraform-google-secured-data-warehouse",
    "terraform-google-secure-cicd",
    "terraform-google-cloud-run",
    "terraform-google-on-prem",
    "terraform-google-secret-manager",
    # Special CI case
    "terraform-example-foundation"
  ]
  branch_repos = setsubtract(toset(local.repos), local.remove_branch_repos)
}

data "github_repository" "repo" {
  for_each = local.branch_repos
  name = each.value
}

data "github_team" "cft-admins" {
  slug = "cft-admins"
}

resource "github_branch_protection" "master" {
  for_each      = data.github_repository.repo
  repository_id = each.value.node_id
  pattern       = "master"

  required_pull_request_reviews {
    required_approving_review_count = 1
    require_code_owner_reviews      = true
  }

  required_status_checks {
    strict   = true
    contexts = [
      "cla/google",
      "${each.value.name}-int-trigger (cloud-foundation-cicd)",
      "${each.value.name}-lint-trigger (cloud-foundation-cicd)"
    ]
  }

  enforce_admins = false

  push_restrictions = [
    data.github_team.cft-admins.node_id
  ]

}

data "github_repository" "terraform-example-foundation" {
  name = "terraform-example-foundation"
}

resource "github_branch_protection" "terraform-example-foundation" {
  repository_id = data.github_repository.terraform-example-foundation.node_id
  pattern       = "master"

  required_pull_request_reviews {
    required_approving_review_count = 1
    require_code_owner_reviews      = true
  }

  required_status_checks {
    strict   = true
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
