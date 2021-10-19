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

resource "github_actions_secret" "infra_secret_gcr_project" {
  repository      = local.gh_repos.infra
  secret_name     = "GCR_PROJECT_ID"
  plaintext_value = local.project_id
}

resource "github_actions_secret" "wif_provider" {
  repository      = local.gh_repos.infra
  secret_name     = "GCP_WIF_PROVIDER"
  plaintext_value = module.oidc.provider_name
}

resource "github_actions_secret" "wif_sa" {
  repository      = local.gh_repos.infra
  secret_name     = "GCP_WIF_SA_EMAIL"
  plaintext_value = module.service_accounts.service_account.email
}
