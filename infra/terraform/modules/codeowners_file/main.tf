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
  commit_author = "CFT Bot"
  commit_email  = "cloud-foundation-bot@google.com"
  owners        = { for value in var.repos_map : value.name => "${join(" ", formatlist("@%s", value.owners))} " if length(value.owners) > 0 }
  groups        = { for value in var.repos_map : value.name => "${join(" ", formatlist("@${value.org}/%s", value.groups))} " if length(value.groups) > 0 }
  header        = "# NOTE: This file is automatically generated from values at:\n# https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/infra/terraform/test-org/org/locals.tf\n"
  footer_prefix = "# NOTE: GitHub CODEOWNERS locations:\n# https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners#codeowners-and-branch-protection\n
  footer_code   = "CODEOWNERS @${var.org}/${var.owner}\n .github/CODEOWNERS @${var.org}/${var.owner}\n docs/CODEOWNERS @${var.org}/${var.owner}\n"
  footer        = "\n${local.footer_prefix}\n{local.footer_code}\n"
}

resource "github_repository_file" "CODEOWNERS" {
  for_each            = var.repo_list
  repository          = each.key
  branch              = each.value.default_branch
  file                = "CODEOWNERS"
  commit_message      = "chore: update CODEOWNERS"
  commit_author       = local.commit_author
  commit_email        = local.commit_email
  overwrite_on_create = true
  content             = "${trimspace("${local.header}\n* @${var.org}/${var.owner} ${try(local.owners[each.key], "")}${try(local.groups[each.key], "")}")}\n${local.footer}"
}
