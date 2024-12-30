/**
 * Copyright 2022-2024 Google LLC
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
  tgm_modules_map = { for value in local.modules : value.name => value if value.org == "terraform-google-modules" }
  gcp_modules_map = { for value in local.modules : value.name => value if value.org == "GoogleCloudPlatform" }
}

data "github_team" "cft-admins" {
  slug     = "cft-admins"
  provider = github
}

data "github_team" "blueprint-solutions" {
  slug     = "blueprint-solutions"
  provider = github.gcp
}

module "repos_tgm" {
  source    = "../../modules/repositories"
  repos_map = local.tgm_modules_map
  ci_teams  = ["blueprint-contributors"]
  providers = { github = github }
}

module "repos_gcp" {
  source    = "../../modules/repositories"
  repos_map = local.gcp_modules_map
  ci_teams  = ["blueprint-contributors"]
  providers = { github = github.gcp }
}

// All new repos are created in advance in the GCP org
import {
  for_each = local.gcp_modules_map
  to       = module.repos_gcp.github_repository.repo[each.value.name]
  id       = each.value.name
}

// terraform-example-foundation CI is a special case - below
module "branch_protection_tgm" {
  source    = "../../modules/branch_protection"
  repo_list = { for k, v in module.repos_tgm.repos : k => v if k != "terraform-example-foundation" }
  repos_map = local.tgm_modules_map
  admin     = data.github_team.cft-admins.node_id
  providers = { github = github }
}

module "branch_protection_gcp" {
  source    = "../../modules/branch_protection"
  repo_list = module.repos_gcp.repos
  repos_map = local.gcp_modules_map
  admin     = data.github_team.blueprint-solutions.node_id
  providers = { github = github.gcp }
}

// terraform-example-foundation renovate is a special case
module "renovate_json_tgm" {
  source    = "../../modules/repo_file"
  repo_list = { for k, v in module.repos_tgm.repos : k => v if k != "terraform-example-foundation" }
  filename  = ".github/renovate.json"
  content   = file("${path.module}/resources/renovate-repo-config.json")
  providers = { github = github }
}

module "renovate_json_gcp" {
  source    = "../../modules/repo_file"
  repo_list = module.repos_gcp.repos
  filename  = ".github/renovate.json"
  content   = file("${path.module}/resources/renovate-repo-config.json")
  providers = { github = github.gcp }
}

module "stale_yml_tgm" {
  source    = "../../modules/repo_file"
  repo_list = module.repos_tgm.repos
  filename  = ".github/workflows/stale.yml"
  content   = file("${path.module}/resources/stale.yml")
  providers = { github = github }
}

module "stale_yml_gcp" {
  source    = "../../modules/repo_file"
  repo_list = module.repos_gcp.repos
  filename  = ".github/workflows/stale.yml"
  content   = file("${path.module}/resources/stale.yml")
  providers = { github = github.gcp }
}

module "conventional-commit-lint_yaml_tgm" {
  source    = "../../modules/repo_file"
  repo_list = module.repos_tgm.repos
  filename  = ".github/conventional-commit-lint.yaml"
  content   = file("${path.module}/resources/conventional-commit-lint.yaml")
  providers = { github = github }
}

module "conventional-commit-lint_yaml_gcp" {
  source    = "../../modules/repo_file"
  repo_list = module.repos_gcp.repos
  filename  = ".github/conventional-commit-lint.yaml"
  content   = file("${path.module}/resources/conventional-commit-lint.yaml")
  providers = { github = github.gcp }
}

module "trusted-contribution_yml_tgm" {
  source    = "../../modules/repo_file"
  repo_list = module.repos_tgm.repos
  filename  = ".github/trusted-contribution.yml"
  content   = file("${path.module}/resources/trusted-contribution.yml")
  providers = { github = github }
}

module "trusted-contribution_yml_gcp" {
  source    = "../../modules/repo_file"
  repo_list = module.repos_gcp.repos
  filename  = ".github/trusted-contribution.yml"
  content   = file("${path.module}/resources/trusted-contribution.yml")
  providers = { github = github.gcp }
}

module "codeowners_tgm" {
  source    = "../../modules/codeowners_file"
  org       = "terraform-google-modules"
  providers = { github = github }
  owner     = "cft-admins"
  repos_map = local.tgm_modules_map
  repo_list = module.repos_tgm.repos
}

module "codeowners_gcp" {
  source    = "../../modules/codeowners_file"
  org       = "GoogleCloudPlatform"
  providers = { github = github.gcp }
  owner     = "blueprint-solutions"
  repos_map = local.gcp_modules_map
  repo_list = module.repos_gcp.repos
}

module "lint_yaml_tgm" {
  source    = "../../modules/workflow_files"
  repos_map = local.tgm_modules_map
  repo_list = module.repos_tgm.repos
  providers = { github = github }
}

module "lint_yaml_gcp" {
  source    = "../../modules/workflow_files"
  repos_map = local.gcp_modules_map
  repo_list = module.repos_gcp.repos
  providers = { github = github.gcp }
}

# Special CI/branch protection case

resource "github_branch_protection" "terraform-example-foundation" {
  repository_id = module.repos_tgm.repos["terraform-example-foundation"].node_id
  pattern       = module.repos_tgm.repos["terraform-example-foundation"].default_branch

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
      "lint",
      "conventionalcommits.org"
    ]
  }

  enforce_admins = false

  restrict_pushes {
    push_allowances = [data.github_team.cft-admins.node_id]
  }
}
