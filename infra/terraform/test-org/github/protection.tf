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
  gcp_repos      = setintersection(local.filtered_repos, toset(data.github_repositories.repos_gcp.names))
  tgm_repos      = setintersection(local.filtered_repos, toset(data.github_repositories.repos_tgm.names))
}

data "github_repositories" "repos_gcp" {
  query = "org:GoogleCloudPlatform archived:no"
}

data "github_repositories" "repos_tgm" {
  query = "org:terraform-google-modules archived:no"
}

module "branch_protection_tgm" {
  source    = "../../modules/branch_protection"
  org       = "terraform-google-modules"
  repo_list = local.tgm_repos
  admin     = "cft-admins"
}

module "branch_protection_gcp" {
  source    = "../../modules/branch_protection"
  org       = "GoogleCloudPlatform"
  repo_list = local.gcp_repos
  admin     = "blueprint-solutions"
}

module "renovate_json_tgm" {
  source    = "../../modules/repo_file"
  org       = "terraform-google-modules"
  repo_list = local.tgm_repos
  filename  = ".github/renovate.json"
  content   = file("${path.module}/renovate.json")
}

module "renovate_json_gcp" {
  source    = "../../modules/repo_file"
  org       = "GoogleCloudPlatform"
  repo_list = local.gcp_repos
  filename  = ".github/renovate.json"
  content   = file("${path.module}/renovate.json")
}

module "stale_yml_tgm" {
  source    = "../../modules/repo_file"
  org       = "terraform-google-modules"
  repo_list = local.tgm_repos
  filename  = ".github/workflows/stale.yml"
  content   = file("${path.module}/stale.yml")
}

module "stale_yml_gcp" {
  source    = "../../modules/repo_file"
  org       = "GoogleCloudPlatform"
  repo_list = local.gcp_repos
  filename  = ".github/workflows/stale.yml"
  content   = file("${path.module}/stale.yml")
}

module "conventional-commit-lint_yaml_tgm" {
  source    = "../../modules/repo_file"
  org       = "terraform-google-modules"
  repo_list = setunion(local.tgm_repos, ["terraform-example-foundation"])
  filename  = ".github/conventional-commit-lint.yaml"
  content   = file("${path.module}/resources/conventional-commit-lint.yaml")
}

module "conventional-commit-lint_yaml_gcp" {
  source    = "../../modules/repo_file"
  org       = "GoogleCloudPlatform"
  repo_list = local.gcp_repos
  filename  = ".github/conventional-commit-lint.yaml"
  content   = file("${path.module}/resources/conventional-commit-lint.yaml")
}

module "codeowners_tgm" {
  source = "../../modules/codeowners_file"
  org    = "terraform-google-modules"
  providers = {
    github = github
  }
  # Add in terraform-example-foundation for CODEOWNERS
  repo_list  = setunion(local.tgm_repos, ["terraform-example-foundation"])
  owner      = "cft-admins"
  add_owners = local.add_owners
}

module "codeowners_gcp" {
  source = "../../modules/codeowners_file"
  org    = "GoogleCloudPlatform"
  providers = {
    github = github.gcp
  }
  repo_list  = local.gcp_repos
  owner      = "blueprint-solutions"
  add_owners = local.add_owners
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
      "terraform-example-foundation-lint-trigger (cloud-foundation-cicd)",
      "conventionalcommits.org"
    ]
  }

  enforce_admins = false

  push_restrictions = [
    data.github_team.cft-admins.node_id
  ]

}
