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

provider "github" {
  owner = var.org
}

locals {
  repos = setintersection(var.repo_list, toset(data.github_repositories.repos.names))
}

data "github_team" "admin" {
  slug = var.admin
}

data "github_repositories" "repos" {
  query = "org:${var.org} archived:no"
}

data "github_repository" "repo" {
  for_each = local.repos
  name     = each.value
}

resource "github_branch_protection" "default" {
  for_each      = data.github_repository.repo
  repository_id = each.value.node_id
  pattern       = each.value.default_branch

  required_pull_request_reviews {
    required_approving_review_count = 1
    require_code_owner_reviews      = true
  }

  required_status_checks {
    strict = true
    contexts = [
      "cla/google",
      "${each.value.name}-int-trigger (cloud-foundation-cicd)",
      "${each.value.name}-lint-trigger (cloud-foundation-cicd)"
    ]
  }

  enforce_admins   = false
  blocks_creations = false

  push_restrictions = [
    data.github_team.admin.node_id
  ]

}
