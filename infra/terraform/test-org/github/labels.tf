/**
 * Copyright 2019 Google LLC
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
  repo_labels = {
    for o in flatten([
      for repo in local.repos :
      [
        for label in local.labels :
        {
          "repo" : repo,
          "label" : label.name,
          "color" : label.color,
          "description" : label.description
        }
      ]
    ]) :
    "${o.repo}/${o.label}" => o
  }
}

# Create labels on all repos
resource "github_issue_label" "test_repo" {
  for_each    = local.repo_labels
  repository  = each.value.repo
  name        = each.value.label
  color       = each.value.color
  description = each.value.description
}
