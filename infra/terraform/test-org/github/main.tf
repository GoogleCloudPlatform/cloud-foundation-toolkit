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
  gh_org = "terraform-google-modules"
  repos = keys(data.terraform_remote_state.triggers.outputs.repo_folder)
}

provider "github" {
  token        = var.github_token
  organization = local.gh_org
}

# Create a new, test colored label
resource "github_issue_label" "test_repo" {
  repository = local.repos[0]
  name       = "Test"
  color      = "FF0000"
}
