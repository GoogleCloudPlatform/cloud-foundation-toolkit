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
  owners = flatten([
    for repo, val in var.repos_map : [
      for owner in val.owners : {
        "repo" : repo
        "owner" : owner
      }
    ]
  ])

  groups = flatten([
    for repo, val in var.repos_map : [
      for group in val.groups : {
        "repo" : repo
        "group" : group
      }
    ]
  ])

  teams = flatten([
    for repo, val in var.repos_map : [
      for team in var.ci_teams : {
        "repo" : repo
        "team" : team
      }
    ]
  ])

}

resource "github_repository" "repo" {
  for_each     = var.repos_map
  name         = each.value.name
  description  = each.value.description
  homepage_url = coalesce(each.value.homepage_url, "https://registry.terraform.io/modules/${each.value.org}/${trimprefix(each.value.name, "terraform-google-")}/google")
  topics       = setunion(["cft-terraform"], try(split(",", trimspace(each.value.topics)), []))

  allow_merge_commit          = false
  allow_rebase_merge          = false
  allow_update_branch         = true
  delete_branch_on_merge      = true
  has_issues                  = true
  has_projects                = false
  has_wiki                    = false
  vulnerability_alerts        = true
  has_downloads               = false
  squash_merge_commit_message = "BLANK"
  squash_merge_commit_title   = "PR_TITLE"
}

resource "github_repository_collaborator" "dpebot" {
  for_each   = github_repository.repo
  repository = each.value.name
  username   = "dpebot"
  permission = "pull"
}

resource "github_repository_collaborator" "owners" {
  for_each = {
    for v in local.owners : "${v.repo}/${v.owner}" => v
  }
  repository = each.value.repo
  username   = each.value.owner
  permission = "admin"
}

resource "github_team_repository" "groups" {
  for_each = {
    for v in local.groups : "${v.repo}/${v.group}" => v
  }
  repository = each.value.repo
  team_id    = each.value.group
  permission = "admin"
}

resource "github_team_repository" "ci_teams" {
  for_each = {
    for v in local.teams : "${v.repo}/${v.team}" => v
  }
  repository = each.value.repo
  team_id    = each.value.team
  permission = "triage"
}
