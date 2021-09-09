/**
 * Copyright 2021 Google LLC
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
  prow_project_id = "blueprints-prow"
  test_ns         = "test-pods"
}

data "google_container_cluster" "prow_build_cluster" {
  name     = "blueprints-prow"
  location = "us-west1-b"
  project  = local.prow_project_id
}

data "google_client_config" "default" {}

provider "kubernetes" {
  host                   = "https://${data.google_container_cluster.prow_build_cluster.endpoint}"
  token                  = data.google_client_config.default.access_token
  cluster_ca_certificate = base64decode(data.google_container_cluster.prow_build_cluster.master_auth[0].cluster_ca_certificate)
}


module "prow-int-sa-wi" {
  source     = "terraform-google-modules/kubernetes-engine/google//modules/workload-identity"
  version    = "~> 16.0"
  name       = "int-test-sa"
  namespace  = local.test_ns
  project_id = local.prow_project_id
}

resource "kubernetes_config_map" "test-constants" {
  metadata {
    name      = "test-constants"
    namespace = local.test_ns
  }

  data = {
    ORG_ID          = local.org_id
    BILLING_ACCOUNT = local.billing_account
    FOLDER_ID       = replace(module.folders-ci.ids["ci-blueprints"], "folders/", "")
  }
}
