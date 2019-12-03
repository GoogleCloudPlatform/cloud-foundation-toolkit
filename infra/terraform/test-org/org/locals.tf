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
  modules = [
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
  ]
}
