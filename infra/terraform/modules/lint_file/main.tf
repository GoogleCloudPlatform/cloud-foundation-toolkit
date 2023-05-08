/**
 * Copyright 2023 Google LLC
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
}

data "github_repository" "repo" {
  for_each = toset(flatten([for value in var.repos_map : value.name if value.disable_lint_yaml != true]))
  name     = each.value
}

resource "github_repository_file" "file" {
  for_each            = data.github_repository.repo
  repository          = each.value.name
  branch              = each.value.default_branch
  file                = ".github/workflows/lint.yaml"
  commit_message      = "chore: update .github/workflows/lint.yaml"
  commit_author       = local.commit_author
  commit_email        = local.commit_email
  overwrite_on_create = true
  content             = templatefile("${path.module}/lint.yaml.tftpl", { branch = each.value.default_branch, lint_env = var.repos_map[each.value.name].lint_env })
}
