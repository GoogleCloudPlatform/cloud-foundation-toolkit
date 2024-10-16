/**
 * Copyright 2020 Google LLC
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

module "service_accounts" {
  source  = "terraform-google-modules/service-accounts/google"
  version = "~> 4.1"

  project_id = local.project_id

  names = ["cft-github-actions"]
  project_roles = [
    "${local.project_id}=>roles/storage.admin"
  ]
}

module "oidc" {
  source  = "terraform-google-modules/github-actions-runners/google//modules/gh-oidc"
  version = "~> 4.0"

  project_id  = local.project_id
  pool_id     = "cft-pool"
  provider_id = "cft-gh-provider"
  sa_mapping = {
    cft-github-actions = {
      sa_name   = module.service_accounts.service_account.name
      attribute = "attribute.repository/GoogleCloudPlatform/cloud-foundation-toolkit"
    }
  }
}
