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

module "folders-root" {
  source  = "terraform-google-modules/folders/google"
  version = "~> 1.0"

  parent_id   = local.org_id
  parent_type = "organization"

  names = [
    "ci-projects"
  ]

  set_roles = false
}

module "folders-ci" {
  source  = "terraform-google-modules/folders/google"
  version = "~> 1.0"

  parent_id   = replace(local.folders["ci"], "folders/", "")
  parent_type = "folder"

  names = [for module in local.modules: "ci-${module}"]

  set_roles = false
}
