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

resource "github_repository" "repo" {
  for_each     = var.repos_map
  name         = each.value.name
  description  = try(each.value.description, null)
  homepage_url = try(each.value.url, "https://registry.terraform.io/modules/${each.value.org}/${trimprefix(each.value.name, "terraform-google-")}/google")
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
