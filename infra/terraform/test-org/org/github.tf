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
  owners = distinct(
    flatten(
      [for repo, val in local.repos : [for owner in val.owners : lower(owner)] if try(val.owners != null, false)]
    )
  )

  org_members = setunion(
    [for login in data.github_organization.tgm.users[*].login : lower(login)],
    [for login in data.github_organization.gcp.users[*].login : lower(login)]
  )

  invalid_owners = setsubtract(local.owners, local.org_members)
}

data "github_organization" "tgm" {
  name = "terraform-google-modules"
}

data "github_organization" "gcp" {
  name = "GoogleCloudPlatform"
}
