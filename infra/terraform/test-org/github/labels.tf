/**
 * Copyright 2019-2024 Google LLC
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
  labels = [
    {
      name : "enhancement",
      color : "a2eeef",
      description : "New feature or request"
    },
    {
      name : "bug",
      color : "d73a4a"
      description : "Something isn't working"
    },
    {
      name : "good first issue",
      color : "7057ff"
      description : "Good for newcomers"
    },
    {
      name : "triaged",
      color : "322560",
      description : "Scoped and ready for work"
    },
    {
      name : "release-please:force-run",
      color : "e7d87d",
      description : "Force release-please to check for changes."
    },
  ]
}

module "repo_labels_gcp" {
  source    = "../../modules/repo_labels"
  org       = "GoogleCloudPlatform"
  repo_list = [for k, v in module.repos_gcp.repos : k]
  labels    = local.labels
}

module "repo_labels_tgm" {
  source    = "../../modules/repo_labels"
  org       = "terraform-google-modules"
  repo_list = setunion([for k, v in module.repos_tgm.repos : k], ["terraform-docs-samples"])
  labels    = local.labels
}
