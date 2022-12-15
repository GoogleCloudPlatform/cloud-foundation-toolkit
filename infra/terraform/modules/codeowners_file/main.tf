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
  commit_author = "CFT Bot"
  commit_email  = "cloud-foundation-bot@google.com"
  owners        = { for value in var.repos_map : value.name => value.owners if can(value.owners) }
}

data "github_repository" "repo" {
  for_each = var.repos_map
  name     = each.value.name
}

resource "github_repository_file" "CODEOWNERS" {
  for_each            = data.github_repository.repo
  repository          = each.value.name
  branch              = each.value.default_branch
  file                = "CODEOWNERS"
  commit_message      = "chore: update CODEOWNERS"
  commit_author       = local.commit_author
  commit_email        = local.commit_email
  overwrite_on_create = true
  content             = "${trimspace("* @${var.org}/${var.owner} ${try(local.owners[each.value.name], "")}")}\n"
}
