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
  remove_repo_labels = [
    "anthos-samples"
  ]
  sub_repos_labels     = setsubtract(toset(local.repos), local.remove_repo_labels)
  sub_repos_labels_gcp = setintersection(local.sub_repos_labels, toset(data.github_repositories.repos_gcp.names))
  sub_repos_labels_tgm = setintersection(local.sub_repos_labels, toset(data.github_repositories.repos_tgm.names))
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
      name : "duplicate",
      color : "cfd3d7"
      description : "This issue or pull request already exists"
    },
    {
      name : "good first issue",
      color : "7057ff"
      description : "Good for newcomers"
    },
    {
      name : "help wanted",
      color : "008672",
      description : "Extra attention is needed"
    },
    {
      name : "invalid",
      color : "e4e669",
      description : "Something doesn't seem right"
    },
    {
      name : "question",
      color : "d876e3",
      description : "Further information is requested"
    },
    {
      name : "wontfix",
      color : "db643d",
      description : "This will not be worked on"
    },
    {
      name : "triaged",
      color : "322560",
      description : "Scoped and ready for work"
    },
    {
      name : "upstream",
      color : "B580D1",
      description : "Work required on Terraform core or provider"
    },
    {
      name : "security",
      color : "801336",
      description : "Fixes a security vulnerability or lapse in best practice"
    },
    {
      name : "refactor",
      color : "004445",
      description : "Updates for readability, code cleanliness, DRYness, etc. Only needs Terraform exp."
    },
    {
      name : "blocked",
      color : "ef4339",
      description : "Blocked by some other work"
    },
    {
      name : "P1",
      color : "b01111",
      description : "highest priority issues"
    },
    {
      name : "P2",
      color : "b4451f",
      description : "high priority issues"
    },
    {
      name : "P3",
      color : "e7d87d",
      description : "medium priority issues"
    },
    {
      name : "P4",
      color : "62a1db",
      description : "low priority issues"
    },
    {
      name : "release-please:force-run",
      color : "e7d87d",
      description : "Force release-please to check for changes."
    },
    {
      name : "waiting-response",
      color : "5319e7",
      description : "Waiting for issue author to respond."
    },
    {
      name : "v0.13",
      color : "edb761",
      description : "Terraform v0.13 issue."
    },
  ]
}

module "repo_labels_gcp" {
  source    = "../../modules/repo_labels"
  org       = "GoogleCloudPlatform"
  repo_list = local.sub_repos_labels_gcp
  labels    = local.labels
}

module "repo_labels_tgm" {
  source    = "../../modules/repo_labels"
  org       = "terraform-google-modules"
  repo_list = local.sub_repos_labels_tgm
  labels    = local.labels
}
