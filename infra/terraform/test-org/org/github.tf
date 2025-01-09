/**
 * Copyright 2023-2025 Google LLC
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
provider "github" {}

locals {
  repo_members = flatten(
    [for repo, val in local.repos : [for member in setunion(lookup(val, "admins", []), lookup(val, "maintainers", [])) : member]]
  )

  org_members = [for login in setunion(data.github_organization.tgm.users[*].login, data.github_organization.gcp.users[*].login) : login]
}

data "github_organization" "tgm" {
  name = "terraform-google-modules"
}

data "github_organization" "gcp" {
  name = "GoogleCloudPlatform"
}
