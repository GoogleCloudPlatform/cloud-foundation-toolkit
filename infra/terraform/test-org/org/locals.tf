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
  gcp_admins_group = "gcp-admins@test.infra.cft.tips"
  ci_project_id    = "cloud-foundation-cicd"

  folders = {
    "ci-projects" = module.folders-root.ids["ci-projects"]
    "ci-shared"   = module.folders-root.ids["ci-shared"]
  }

  ci_folders = module.folders-ci.ids
  ci_repos_folders = {
    for repo in local.repos : try(repo.short_name, trimprefix(repo.name, "terraform-google-")) => {
      folder_name = "ci-${try(repo.short_name, trimprefix(repo.name, "terraform-google-"))}",
      folder_id   = replace(module.folders-ci.ids["ci-${try(repo.short_name, trimprefix(repo.name, "terraform-google-"))}"], "folders/", ""),
      gh_org      = repo.org
    }
  }

  common_topics = {
    hcls       = "healthcare-life-sciences",
    e2e        = "end-to-end"
    serverless = "serverless-computing",
    compute    = "compute"
    containers = "containers",
    db         = "databases",
    da         = "data-analytics",
    storage    = "storage",
    ops        = "operations",
    net        = "networking",
    security   = "security-identity",
    devtools   = "developer-tools"
    workspace  = "workspace"
  }

  /*
 *  repos schema
 *  name         = "string" (required for modules)
 *  short_name   = "string" (optional for modules, if not prefixed with 'terraform-google-')
 *  org          = "terraform-google-modules" or "GoogleCloudPlatform" (required)
 *  description  = "string" (required)
 *  owners       = "@user1 @user2" (optional)
 *  homepage_url = "string" (optional, overrides default)
 *  module       = BOOL (optional, default is true which includes GH repo configuration)
 *  topics       = "string1,string2,string3" (one or more of local.common_topics required if module = true)
 *
 */

  repos = [
    {
      name        = "cloud-foundation-training"
      org         = "terraform-google-modules"
      description = ""
      owners      = "@marine675 @zefdelgadillo"
    },
    {
      name        = "terraform-google-healthcare"
      org         = "terraform-google-modules"
      description = "Handles opinionated Google Cloud Healthcare datasets and stores"
      owners      = "@yeweidaniel"
      topics      = local.common_topics.hcls
    },
    {
      name        = "terraform-google-cloud-run"
      org         = "GoogleCloudPlatform"
      description = "Deploys apps to Cloud Run, along with option to map custom domain"
      owners      = "@prabhu34 @anamer @mitchelljamie"
      topics      = "cloudrun,google-cloud-platform,terraform-modules,${local.common_topics.serverless}"
    },
    {
      name        = "terraform-google-secured-data-warehouse"
      org         = "GoogleCloudPlatform"
      description = "Deploys a secured BigQuery data warehouse"
      owners      = "@erlanderlo"
      topics      = join(",", [local.common_topics.da, local.common_topics.e2e])
    },
    {
      name        = "terraform-google-anthos-vm"
      org         = "GoogleCloudPlatform"
      description = "Creates VMs on Anthos Bare Metal clusters"
      owners      = "@zhuchenwang"
      topics      = "anthos,kubernetes,terraform-module,vm,${local.common_topics.compute}"
    },
    {
      name        = "terraform-google-kubernetes-engine"
      org         = "terraform-google-modules"
      description = "Configures opinionated GKE clusters"
      owners      = "@Jberlinsky @ericyz"
      topics      = join(",", [local.common_topics.compute, local.common_topics.containers])
    },
    {
      name         = "terraform-example-foundation"
      short_name   = "example-foundation"
      org          = "terraform-google-modules"
      description  = "Shows how the CFT modules can be composed to build a secure cloud foundation"
      owners       = "@rjerrems"
      homepage_url = "https://cloud.google.com/architecture/security-foundations"
      topics       = join(",", [local.common_topics.e2e, local.common_topics.ops])
    },
    {
      name        = "terraform-google-log-analysis"
      org         = "GoogleCloudPlatform"
      description = "Stores and analyzes log data"
      owners      = "@ryotat7"
      topics      = local.common_topics.da
    },
    {
      name        = "terraform-google-three-tier-web-app"
      org         = "GoogleCloudPlatform"
      description = "Deploys a three tier web application using Cloud Run and Cloud SQL"
      owners      = "@tpryan"
      topics      = join(",", [local.common_topics.serverless, local.common_topics.db])
    },
    {
      name        = "terraform-google-load-balanced-vms"
      org         = "GoogleCloudPlatform"
      description = "Creates a Managed Instance Group with a loadbalancer"
      owners      = "@tpryan"
      topics      = local.common_topics.net
    },
    {
      name        = "terraform-google-secure-cicd"
      org         = "GoogleCloudPlatform"
      description = "Builds a secure CI/CD pipeline on Google Cloud"
      owners      = "@gtsorbo"
      topics      = join(",", [local.common_topics.security, local.common_topics.devtools, local.common_topics.e2e])
    },
    {
      name        = "terraform-google-media-cdn-vod"
      org         = "GoogleCloudPlatform"
      description = "Deploys Media CDN video-on-demand"
      owners      = "@roddzurcher"
      topics      = local.common_topics.ops
    },
    {
      name        = "terraform-example-foundation-app"
      short_name  = "example-foundation-app"
      org         = "GoogleCloudPlatform"
      description = ""
    },
    {
      name        = "terraform-google-network-forensics"
      org         = "GoogleCloudPlatform"
      description = "Deploys Zeek on Google Cloud"
      topics      = local.common_topics.net
    },
    {
      name        = "terraform-google-secret-manager"
      org         = "GoogleCloudPlatform"
      description = "Creates one or more Google Secret Manager secrets and manages basic permissions for them"
      topics      = "gcp,kms,pubsub,terraform-module,${local.common_topics.security}"
    },
    {
      name        = "terraform-google-address"
      org         = "terraform-google-modules"
      description = "Manages Google Cloud IP addresses"
      topics      = local.common_topics.net
    },
    {
      name        = "terraform-google-bastion-host"
      org         = "terraform-google-modules"
      description = "Generates a bastion host VM compatible with OS Login and IAP Tunneling that can be used to access internal VMs"
      topics      = join(",", [local.common_topics.security, local.common_topics.ops, local.common_topics.devtools])
    },
    {
      name        = "terraform-google-bigquery"
      org         = "terraform-google-modules"
      description = "Creates opinionated BigQuery datasets and tables"
      topics      = local.common_topics.da
    },
    {
      name        = "terraform-google-bootstrap"
      org         = "terraform-google-modules"
      description = "Bootstraps Terraform usage and related CI/CD in a new Google Cloud organization"
      topics      = join(",", [local.common_topics.ops, local.common_topics.devtools])
    },
    {
      name        = "terraform-google-cloud-datastore"
      org         = "terraform-google-modules"
      description = "Manages Datastore"
      topics      = local.common_topics.db
    },
    {
      name        = "terraform-google-cloud-dns"
      org         = "terraform-google-modules"
      description = "Creates and manages Cloud DNS public or private zones and their records"
      topics      = local.common_topics.net
    },
    {
      name        = "terraform-google-cloud-nat"
      org         = "terraform-google-modules"
      description = "Creates and configures Cloud NAT"
      topics      = local.common_topics.net
    },
    {
      name        = "terraform-google-cloud-operations"
      org         = "terraform-google-modules"
      description = "Manages Google Cloud's operations suite (Cloud Logging and Cloud Monitoring)"
      topics      = local.common_topics.ops
    },
    {
      name        = "terraform-google-cloud-router"
      org         = "terraform-google-modules"
      description = "Manages a Cloud Router on Google Cloud"
      topics      = local.common_topics.net
    },
    {
      name        = "terraform-google-cloud-storage"
      org         = "terraform-google-modules"
      description = "Creates one or more Cloud Storage buckets and assigns basic permissions on them to arbitrary users"
      topics      = local.common_topics.storage
    },
    {
      name        = "terraform-google-composer"
      org         = "terraform-google-modules"
      description = "Manages Cloud Composer v1 and v2 along with option to manage networking"
      topics      = join(",", [local.common_topics.da, local.common_topics.ops, local.common_topics.e2e])
    },
    {
      name        = "terraform-google-container-vm"
      org         = "terraform-google-modules"
      description = "Deploys containers on Compute Engine instances"
      topics      = join(",", [local.common_topics.containers, local.common_topics.compute])
    },
    {
      name        = "terraform-google-data-fusion"
      org         = "terraform-google-modules"
      description = "Manages Cloud Data Fusion"
      topics      = local.common_topics.da
    },
    {
      name        = "terraform-google-dataflow"
      org         = "terraform-google-modules"
      description = "Handles opinionated Dataflow job configuration and deployments"
      topics      = local.common_topics.da
    },
    {
      name        = "terraform-google-datalab"
      org         = "terraform-google-modules"
      description = "Creates DataLab instances with support for GPU instances"
      topics      = local.common_topics.da
    },
    {
      name        = "terraform-google-event-function"
      org         = "terraform-google-modules"
      description = "Responds to logging events with a Cloud Function"
      topics      = local.common_topics.serverless
    },
    {
      name        = "terraform-google-folders"
      org         = "terraform-google-modules"
      description = "Creates several Google Cloud folders under the same parent"
      topics      = local.common_topics.devtools
    },
    {
      name        = "terraform-google-gcloud"
      org         = "terraform-google-modules"
      description = "Executes Google Cloud CLI commands within Terraform"
      topics      = local.common_topics.devtools
    },
    {
      name        = "terraform-google-github-actions-runners"
      org         = "terraform-google-modules"
      description = "Creates self-hosted GitHub Actions Runners on Google Cloud"
      topics      = local.common_topics.devtools
    },
    {
      name        = "terraform-google-gke-gitlab"
      org         = "terraform-google-modules"
      description = "Installs GitLab on Kubernetes Engine"
      topics      = local.common_topics.devtools
    },
    {
      name        = "terraform-google-group"
      org         = "terraform-google-modules"
      description = "Manages Google Groups"
      topics      = local.common_topics.workspace
    },
    {
      name        = "terraform-google-gsuite-export"
      org         = "terraform-google-modules"
      description = "Creates a Compute Engine VM instance and sets up a cronjob to export GSuite Admin SDK data to Cloud Logging on a schedule"
      topics      = join(",", [local.common_topics.ops, local.common_topics.workspace])
    },
    {
      name        = "terraform-google-iam"
      org         = "terraform-google-modules"
      description = "Manages multiple IAM roles for resources on Google Cloud"
      topics      = local.common_topics.security
    },
    {
      name        = "terraform-google-jenkins"
      org         = "terraform-google-modules"
      description = "Creates a Compute Engine instance running Jenkins"
      topics      = local.common_topics.devtools
    },
    {
      name        = "terraform-google-kms"
      org         = "terraform-google-modules"
      description = "Allows managing a keyring, zero or more keys in the keyring, and IAM role bindings on individual keys"
      topics      = local.common_topics.security
    },
    {
      name        = "terraform-google-lb"
      org         = "terraform-google-modules"
      description = "Creates a regional TCP proxy load balancer for Compute Engine by using target pools and forwarding rules"
      topics      = local.common_topics.net
    },
    {
      name        = "terraform-google-lb-http"
      org         = "terraform-google-modules"
      description = "Creates a global HTTP load balancer for Compute Engine by using forwarding rules"
      topics      = local.common_topics.net
    },
    {
      name        = "terraform-google-lb-internal"
      org         = "terraform-google-modules"
      description = "Creates an internal load balancer for Compute Engine by using forwarding rules"
      topics      = local.common_topics.net
    },
    {
      name        = "terraform-google-log-export"
      org         = "terraform-google-modules"
      description = "Creates log exports at the project, folder, or organization level"
      topics      = local.common_topics.ops
    },
    {
      name        = "terraform-google-memorystore"
      org         = "terraform-google-modules"
      description = "Creates a fully functional Google Memorystore (redis) instance"
      topics      = local.common_topics.db
    },
    {
      name        = "terraform-google-module-template"
      org         = "terraform-google-modules"
      description = "Provides a template for creating a Cloud Foundation Toolkit Terraform module"
    },
    {
      name        = "terraform-google-network"
      org         = "terraform-google-modules"
      description = "Sets up a new VPC network on Google Cloud"
      topics      = local.common_topics.net
    },
    {
      name        = "terraform-google-org-policy"
      org         = "terraform-google-modules"
      description = "Manages Google Cloud organization policies"
      topics      = local.common_topics.security
    },
    {
      name        = "terraform-google-project-factory"
      org         = "terraform-google-modules"
      description = "Creates an opinionated Google Cloud project by using Shared VPC, IAM, and Google Cloud APIs"
      topics      = local.common_topics.ops
    },
    {
      name        = "terraform-google-pubsub"
      org         = "terraform-google-modules"
      description = "Creates Pub/Sub topic and subscriptions associated with the topic"
      topics      = local.common_topics.da
    },
    {
      name        = "terraform-google-sap"
      org         = "terraform-google-modules"
      description = "Deploys SAP products"
      topics      = local.common_topics.compute
    },
    {
      name        = "terraform-google-scheduled-function"
      org         = "terraform-google-modules"
      description = "Sets up a scheduled job to trigger events and run functions"
      topics      = local.common_topics.serverless
    },
    {
      name        = "terraform-google-service-accounts"
      org         = "terraform-google-modules"
      description = "Creates one or more service accounts and grants them basic roles"
      topics      = local.common_topics.security
    },
    {
      name        = "terraform-google-slo"
      org         = "terraform-google-modules"
      description = "Creates SLOs on Google Cloud from custom Stackdriver metrics capability to export SLOs to Google Cloud services and other systems"
      topics      = local.common_topics.ops
    },
    {
      name        = "terraform-google-sql-db"
      org         = "terraform-google-modules"
      description = "Creates a Cloud SQL database instance"
      topics      = local.common_topics.db
    },
    {
      name        = "terraform-google-startup-scripts"
      org         = "terraform-google-modules"
      description = "Provides a library of useful startup scripts to embed in VMs"
      topics      = local.common_topics.compute
    },
    {
      name        = "terraform-google-utils"
      org         = "terraform-google-modules"
      description = "Gets the short names for a given Google Cloud region"
      topics      = local.common_topics.ops
    },
    {
      name        = "terraform-google-vault"
      org         = "terraform-google-modules"
      description = "Deploys Vault on Compute Engine"
      topics      = "hashicorp-vault,${local.common_topics.ops},${local.common_topics.devtools},${local.common_topics.security}"
    },
    {
      name        = "terraform-google-vm"
      org         = "terraform-google-modules"
      description = "Provisions VMs in Google Cloud"
      topics      = local.common_topics.compute
    },
    {
      name        = "terraform-google-vpc-service-controls"
      org         = "terraform-google-modules"
      description = "Handles opinionated VPC Service Controls and Access Context Manager configuration and deployments"
      topics      = local.common_topics.net
    },
    {
      name        = "terraform-google-vpn"
      org         = "terraform-google-modules"
      description = "Sets up a Cloud VPN gateway"
      topics      = local.common_topics.net
    },
    {
      short_name = "anthos-platform"
      org        = "terraform-google-modules"
      module     = false
    },
    {
      short_name = "anthos-samples"
      org        = "GoogleCloudPlatform"
      module     = false
    },
    {
      short_name = "blueprints"
      org        = "GoogleCloudPlatform"
      module     = false
    },
    {
      short_name = "docs-samples"
      org        = "terraform-google-modules"
      module     = false
    },
    {
      short_name = "migrate"
      org        = "terraform-google-modules"
      module     = false
    },
    {
      short_name = "policy-blueprints"
      org        = "GoogleCloudPlatform"
      module     = false
    },
    {
      short_name = "terraform-validator"
      org        = "terraform-google-modules"
      module     = false
    },
    {
      name        = "terraform-google-waap"
      org         = "GoogleCloudPlatform"
      description = "Deploys the WAAP solution on Google Cloud."
      owners      = "@gtsorbo"
      topics      = local.common_topics.ops
    },
    {
      name        = "terraform-google-cloud-workflows"
      org         = "GoogleCloudPlatform"
      description = "Manage Cloud Workflows with optional Scheduler or Event Arc triggers."
      owners      = "@anaik91"
      topics      = join(",", [local.common_topics.serverless, local.common_topics.devtools])
    },
    {
      name        = "terraform-google-cloud-armor"
      org         = "GoogleCloudPlatform"
      description = "Deploy Cloud Armor Security policy"
      owners      = "@imrannayer @belgana"
      topics      = join(",", [local.common_topics.compute, local.common_topics.net])
    },
    {
      name        = "terraform-google-cloud-deploy"
      org         = "GoogleCloudPlatform"
      description = "Create Cloud Deploy pipelines and targets"
      owners      = "@gtsorbo @niranjankl"
      topics      = join(",", [local.common_topics.devtools])
    },
  ]
}
