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
  org_id           = "943740911108"
  billing_account  = "01D904-DAF6EC-F34EF7"
  cft_ci_group     = "cft-ci-robots@test.infra.cft.tips"
  cft_dev_group    = "cft-developers@dev.infra.cft.tips"
  gcp_admins_group = "gcp-admins@test.infra.cft.tips"

  folders = {
    "ci-projects" = module.folders-root.ids["ci-projects"]
    "ci-shared"   = module.folders-root.ids["ci-shared"]
  }

  ci_folders = module.folders-ci.ids
  ci_repos_folders = merge(
    { for m in local.tgm_org_modules : m => { folder_name = "ci-${m}", folder_id = replace(module.folders-ci.ids["ci-${m}"], "folders/", ""), gh_org = "terraform-google-modules" } },
    { for m in local.gcp_org_modules : m => { folder_name = "ci-${m}", folder_id = replace(module.folders-ci.ids["ci-${m}"], "folders/", ""), gh_org = "GoogleCloudPlatform" } }
  )
  tgm_org_modules = [
    "kms",
    "network",
    "folders",
    "slo",
    "sap",
    "iam",
    "event-function",
    "kubernetes-engine",
    "vpn",
    "project-factory",
    "pubsub",
    "migrate",
    "bootstrap",
    "redis",
    "datalab",
    "mariadb",
    "jenkins",
    "container-vm",
    "lb",
    "vm",
    "memorystore",
    "airflow",
    "service-accounts",
    "cloud-storage",
    "sql-db",
    "vpc-service-controls",
    "cloud-datastore",
    "dataflow",
    "cloud-nat",
    "startup-scripts",
    "scheduled-function",
    "address",
    "bigquery",
    "bastion-host",
    "org-policy",
    "log-export",
    "on-prem",
    "cloud-dns",
    "gsuite-export",
    "secret",
    "terraform-validator",
    "lb-http",
    "gcloud",
    "lb-internal",
    "utils",
    "composer",
    "github-actions-runners",
    "healthcare",
    "gke-gitlab",
    "example-foundation", # Not module
    "anthos-platform",    # Not module
    "cloud-operations",
    "cloud-foundation-training", # Not module
    "cloud-router",
    "group",
  ]
  gcp_org_modules = [
    "example-foundation-app", # Not module
    "anthos-samples",
    "secure-cicd",
    "secured-data-warehouse",
    "cloud-run",
  ]
}
