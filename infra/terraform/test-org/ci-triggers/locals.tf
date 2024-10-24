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
  exclude_folders = [
    "ci",
    "ci-terraform-validator",
    "ci-projects",
    "ci-shared",
    "ci-anthos-platform",
    "ci-example-foundation",
    "ci-blueprints",
    "ci-policy-blueprints",
  ]
  # custom mapping of the form name => repo_name used for overriding `terraform-google` prefix
  custom_repo_mapping = {
    "cloud-foundation-training"  = "cloud-foundation-training",
    "example-foundation-app"     = "terraform-example-foundation-app",
    "anthos-samples"             = "anthos-samples"
    "docs-samples"               = "terraform-docs-samples"
    "dynamic-python-webapp"      = "terraform-dynamic-python-webapp"
    "dynamic-javascript-webapp"  = "terraform-dynamic-javascript-webapp"
    "deploy-java-multizone"      = "terraform-example-deploy-java-multizone"
    "ecommerce-microservices"    = "terraform-ecommerce-microservices-on-gke"
    "deploy-java-gke"            = "terraform-example-deploy-java-gke"
    "java-dynamic-point-of-sale" = "terraform-example-java-dynamic-point-of-sale"
    "ml-image-annotation-gcf"    = "terraform-ml-image-annotation-gcf"
    "genai-doc-summarization"    = "terraform-genai-doc-summarization"
    "genai-knowledge-base"       = "terraform-genai-knowledge-base"
    "secured-notebook"           = "notebooks-blueprint-security"
    "sdw-onprem-ingest"          = "terraform-google-secured-data-warehouse-onprem-ingest"
    "pubsub-golang-app"          = "terraform-pubsub-integration-golang"
    "pubsub-java-app"            = "terraform-pubsub-integration-java"
    "genai-rag"                  = "terraform-genai-rag"
    "cloud-client-api"           = "terraform-cloud-client-api"
    "dataanalytics-eventdriven"  = "terraform-dataanalytics-eventdriven"
    "kms-solutions"              = "kms-solutions"
  }
  # example foundation has custom test modes
  example_foundation                = { "terraform-example-foundation" = data.terraform_remote_state.org.outputs.ci_repos_folders["example-foundation"] }
  example_foundation_int_test_modes = ["default", "HubAndSpoke"]

  repo_folder             = { for key, value in data.terraform_remote_state.org.outputs.ci_repos_folders : contains(keys(local.custom_repo_mapping), key) ? local.custom_repo_mapping[key] : "terraform-google-${key}" => value if !contains(local.exclude_folders, value.folder_name) }
  org_id                  = data.terraform_remote_state.org.outputs.org_id
  billing_account         = data.terraform_remote_state.org.outputs.billing_account
  lr_billing_account      = data.terraform_remote_state.org.outputs.lr_billing_account
  tf_validator_project_id = data.terraform_remote_state.tf-validator.outputs.project_id
  tf_validator_folder_id  = trimprefix(data.terraform_remote_state.org.outputs.folders["ci-terraform-validator"], "folders/")
  # tf validator "ancestry path" expects non-plural type names for historical reasons
  tf_validator_ancestry    = "organization/${local.org_id}/folder/${trimprefix(data.terraform_remote_state.org.outputs.folders["ci-projects"], "folders/")}/folder/${local.tf_validator_folder_id}"
  project_id               = "cloud-foundation-cicd"
  forseti_ci_folder_id     = "542927601143"
  billing_iam_test_account = "0151A3-65855E-5913CF"
  # blueprints which can be layered on top of SFB
  bp_on_sfb = [
    "terraform-google-cloud-run"
  ]
  # SFB deployment info
  sfb_substs = {
    _SFB_ORG_ID : "413973101099",
    _SFB_SEED_PROJECT_ID : data.terraform_remote_state.sfb-bootstrap.outputs.seed_project_id,
    _SFB_CLOUDBUILD_PROJECT_ID : data.terraform_remote_state.sfb-bootstrap.outputs.cloudbuild_project_id,
    _SFB_TF_SA_NAME : data.terraform_remote_state.sfb-bootstrap.outputs.terraform_sa_name,
  }
  # vod test project id
  vod_test_project_id = data.terraform_remote_state.org.outputs.ci_media_cdn_vod_project_id
  # file logger opt-in
  enable_file_log = { "terraform-docs-samples" : true }
}
