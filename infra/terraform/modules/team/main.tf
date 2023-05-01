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

resource "github_team" "collaborators_team" {
  name        = "cft_collaborators_ci"
  description = "CFT collaborators team for approved CI"

  create_default_maintainer = true
}

resource "github_team_members" "collaborators" {
  team_id = github_team.collaborators_team.id

  dynamic "members" {
    for_each = var.users
    content {
      username = members.value
      role     = "member"
    }
  }

}
