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
  commit_author = "CFT Bot"
  commit_email  = "cloud-foundation-bot@google.com"
}

resource "github_repository_file" "file" {
  for_each            = var.repo_list
  repository          = each.key
  branch              = each.value.default_branch
  file                = var.filename
  commit_message      = "chore: update ${var.filename}"
  commit_author       = local.commit_author
  commit_email        = local.commit_email
  overwrite_on_create = true
  content             = var.content
}
