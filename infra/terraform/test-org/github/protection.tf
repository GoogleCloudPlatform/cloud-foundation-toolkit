/**
 * Copyright 2022-2023 Google LLC
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
  users           = distinct(flatten([for value in local.modules : value.owners if try(value.owners, null) != null]))
}

module "repos_tgm" {
  source    = "../../modules/repositories"
  repos_map = local.tgm_modules_map
  team_id   = module.collaborators_tgm.id
  providers = {
    github = github
  }
}

module "repos_gcp" {
  source    = "../../modules/repositories"
  repos_map = local.gcp_modules_map
  # TODO: filter team on users already part of the org
  # team_id   = module.collaborators_gcp.id
  providers = {
    github = github.gcp
  }
}

// terraform-example-foundation CI is a special case - below
module "branch_protection_tgm" {
  source    = "../../modules/branch_protection"
  org       = "terraform-google-modules"
  repo_list = setsubtract(module.repos_tgm.repos, ["terraform-example-foundation"])
  admin     = "cft-admins"
}

module "branch_protection_gcp" {
  source    = "../../modules/branch_protection"
  org       = "GoogleCloudPlatform"
  repo_list = module.repos_gcp.repos
  admin     = "blueprint-solutions"
}

// terraform-example-foundation renovate is a special case - manual
module "renovate_json_tgm" {
  source    = "../../modules/repo_file"
  org       = "terraform-google-modules"
  repo_list = setsubtract(module.repos_tgm.repos, ["terraform-example-foundation"])
  filename  = ".github/renovate.json"
  content   = file("${path.module}/resources/renovate.json")
}

module "renovate_json_gcp" {
  source    = "../../modules/repo_file"
  org       = "GoogleCloudPlatform"
  repo_list = module.repos_gcp.repos
  filename  = ".github/renovate.json"
  content   = file("${path.module}/resources/renovate.json")
}

module "stale_yml_tgm" {
  source    = "../../modules/repo_file"
  org       = "terraform-google-modules"
  repo_list = module.repos_tgm.repos
  filename  = ".github/workflows/stale.yml"
  content   = file("${path.module}/resources/stale.yml")
}

module "stale_yml_gcp" {
  source    = "../../modules/repo_file"
  org       = "GoogleCloudPlatform"
  repo_list = module.repos_gcp.repos
  filename  = ".github/workflows/stale.yml"
  content   = file("${path.module}/resources/stale.yml")
}

module "conventional-commit-lint_yaml_tgm" {
  source    = "../../modules/repo_file"
  org       = "terraform-google-modules"
  repo_list = module.repos_tgm.repos
  filename  = ".github/conventional-commit-lint.yaml"
  content   = file("${path.module}/resources/conventional-commit-lint.yaml")
}

module "conventional-commit-lint_yaml_gcp" {
  source    = "../../modules/repo_file"
  org       = "GoogleCloudPlatform"
  repo_list = module.repos_gcp.repos
  filename  = ".github/conventional-commit-lint.yaml"
  content   = file("${path.module}/resources/conventional-commit-lint.yaml")
}

module "trusted-contribution_yml_tgm" {
  source    = "../../modules/repo_file"
  org       = "terraform-google-modules"
  repo_list = module.repos_tgm.repos
  filename  = ".github/trusted-contribution.yml"
  content   = file("${path.module}/resources/trusted-contribution.yml")
}

module "trusted-contribution_yml_gcp" {
  source    = "../../modules/repo_file"
  org       = "GoogleCloudPlatform"
  repo_list = module.repos_gcp.repos
  filename  = ".github/trusted-contribution.yml"
  content   = file("${path.module}/resources/trusted-contribution.yml")
}

module "codeowners_tgm" {
  source = "../../modules/codeowners_file"
  org    = "terraform-google-modules"
  providers = {
    github = github
  }
  owner     = "cft-admins"
  repos_map = local.tgm_modules_map
}

module "codeowners_gcp" {
  source = "../../modules/codeowners_file"
  org    = "GoogleCloudPlatform"
  providers = {
    github = github.gcp
  }
  owner     = "blueprint-solutions"
  repos_map = local.gcp_modules_map
}

module "lint_yaml_tgm" {
  source    = "../../modules/lint_file"
  repos_map = local.tgm_modules_map
  providers = {
    github = github
  }
}

module "lint_yaml_gcp" {
  source    = "../../modules/lint_file"
  repos_map = local.gcp_modules_map
  providers = {
    github = github.gcp
  }
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
      "lint",
      "conventionalcommits.org"
    ]
  }

  enforce_admins = false

  push_restrictions = [
    data.github_team.cft-admins.node_id
  ]

}

# collaborator teams for approved CI
module "collaborators_tgm" {
  source = "../../modules/team"
  users  = local.users
  providers = {
    github = github
  }
}
