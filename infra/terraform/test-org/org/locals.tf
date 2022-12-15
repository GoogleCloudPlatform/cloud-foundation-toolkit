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

  ci_folders       = module.folders-ci.ids
  ci_repos_folders = {
    for repo in local.repos : try(repo.short_name, trimprefix(repo.name, "terraform-google-")) => {
      folder_name = "ci-${try(repo.short_name, trimprefix(repo.name, "terraform-google-"))}",
      folder_id = replace(module.folders-ci.ids["ci-${try(repo.short_name, trimprefix(repo.name, "terraform-google-"))}"], "folders/", ""),
      gh_org = repo.org
    }
  }

/*
 *  repos schema
 *  name         = "string" (required for modules)
 *  short_name   = "string" (optional for modules, if not prefixed with 'terraform-google-')
 *  org          = "terraform-google-modules" or "GoogleCloudPlatform" (required)
 *  description  = "string" (required)
 *  owners       = "@user1 @user2" (optional)
 *  homepage_url = "string" (optional, overrides default)
 *  topics       = "string1,string2,string3" (optional)
 *  module       = BOOL (optional, default is true which includes GH repo configuration)
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
    },
    {
      name        = "terraform-google-cloud-run"
      org         = "GoogleCloudPlatform"
      description = "Deploys apps to Cloud Run, along with option to map custom domain"
      owners      = "@prabhu34 @anamer @mitchelljamie"
      topics      = "cft-fabric,cloudrun,google-cloud-platform,terraform-modules"
    },
    {
      name        = "terraform-google-secured-data-warehouse"
      org         = "GoogleCloudPlatform"
      description = "Deploys a secured BigQuery data warehouse"
      owners      = "@erlanderlo"
    },
    {
      name        = "terraform-google-anthos-vm"
      org         = "GoogleCloudPlatform"
      description = "Creates VMs on Anthos Bare Metal clusters"
      owners      = "@zhuchenwang"
      topics      = "anthos,kubernetes,terraform-module,vm"
    },
    {
      name        = "terraform-google-kubernetes-engine"
      org         = "terraform-google-modules"
      description = "Configures opinionated GKE clusters"
      owners      = "@Jberlinsky"
    },
    {
      name         = "terraform-example-foundation"
      short_name   = "example-foundation"
      org          = "terraform-google-modules"
      description  = "Shows how the CFT modules can be composed to build a secure cloud foundation"
      owners       = "@rjerrems"
      homepage_url = "https://github.com/terraform-google-modules/terraform-example-foundation"
    },
    {
      name        = "terraform-google-log-analysis"
      org         = "GoogleCloudPlatform"
      description = "Stores and analyzes log data"
      owners      = "@ryotat7"
    },
    {
      name        = "terraform-google-three-tier-web-app"
      org         = "GoogleCloudPlatform"
      description = "Creates a Cloud Storage bucket"
      owners      = "@tpryan"
    },
    {
      name        = "terraform-google-load-balanced-vms"
      org         = "GoogleCloudPlatform"
      description = "Creates a Managed Instance Group with a loadbalancer"
      owners      = "@tpryan"
    },
    {
      name        = "terraform-google-secure-cicd"
      org         = "GoogleCloudPlatform"
      description = "Builds a secure CI/CD pipeline on Google Cloud"
      owners      = "@gtsorbo"
    },
    {
      name        = "terraform-google-media-cdn-vod"
      org         = "GoogleCloudPlatform"
      description = "Deploys Media CDN video-on-demand"
      owners      = "@roddzurcher"
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
    },
    {
      name        = "terraform-google-secret-manager"
      org         = "GoogleCloudPlatform"
      description = "Creates one or more Google Secret Manager secrets and manages basic permissions for them"
      topics      = "gcp,kms,pubsub,terraform-module"
    },
    {
      name        = "terraform-google-address"
      org         = "terraform-google-modules"
      description = "Manages Google Cloud IP addresses"
    },
    {
      name        = "terraform-google-bastion-host"
      org         = "terraform-google-modules"
      description = "Generates a bastion host VM compatible with OS Login and IAP Tunneling that can be used to access internal VMs"
    },
    {
      name        = "terraform-google-bigquery"
      org         = "terraform-google-modules"
      description = "Creates opinionated BigQuery datasets and tables"
    },
    {
      name        = "terraform-google-bootstrap"
      org         = "terraform-google-modules"
      description = "Bootstraps Terraform usage and related CI/CD in a new Google Cloud organization"
    },
    {
      name        = "terraform-google-cloud-datastore"
      org         = "terraform-google-modules"
      description = "Manages Datastore"
    },
    {
      name        = "terraform-google-cloud-dns"
      org         = "terraform-google-modules"
      description = "Creates and manages Cloud DNS public or private zones and their records"
      topics      = "cft-fabric"
    },
    {
      name        = "terraform-google-cloud-nat"
      org         = "terraform-google-modules"
      description = "Creates and configures Cloud NAT"
      topics      = "cft-fabric"
    },
    {
      name        = "terraform-google-cloud-operations"
      org         = "terraform-google-modules"
      description = "Manages Google Cloud's operations suite (Cloud Logging and Cloud Monitoring)"
    },
    {
      name        = "terraform-google-cloud-router"
      org         = "terraform-google-modules"
      description = "Manages a Cloud Router on Google Cloud"
    },
    {
      name        = "terraform-google-cloud-storage"
      org         = "terraform-google-modules"
      description = "Creates one or more Cloud Storage buckets and assigns basic permissions on them to arbitrary users"
      topics      = "cft-fabric"
    },
    {
      name        = "terraform-google-composer"
      org         = "terraform-google-modules"
      description = "Manages Cloud Composer v1 and v2 along with option to manage networking"
    },
    {
      name        = "terraform-google-container-vm"
      org         = "terraform-google-modules"
      description = "Deploys containers on Compute Engine instances"
      topics      = "cft-fabric"
    },
    {
      name        = "terraform-google-data-fusion"
      org         = "terraform-google-modules"
      description = "[ALPHA] Manages Cloud Data Fusion"
    },
    {
      name        = "terraform-google-dataflow"
      org         = "terraform-google-modules"
      description = "Handles opinionated Dataflow job configuration and deployments"
    },
    {
      name        = "terraform-google-datalab"
      org         = "terraform-google-modules"
      description = "Creates DataLab instances with support for GPU instances"
    },
    {
      name        = "terraform-google-event-function"
      org         = "terraform-google-modules"
      description = "Responds to logging events with a Cloud Function"
    },
    {
      name        = "terraform-google-folders"
      org         = "terraform-google-modules"
      description = "Creates several Google Cloud folders under the same parent"
      topics      = "cft-fabric"
    },
    {
      name        = "terraform-google-gcloud"
      org         = "terraform-google-modules"
      description = "Executes Google Cloud CLI commands within Terraform"
    },
    {
      name        = "terraform-google-github-actions-runners"
      org         = "terraform-google-modules"
      description = "[ALPHA] Creates self-hosted GitHub Actions Runners on Google Cloud"
    },
    {
      name        = "terraform-google-gke-gitlab"
      org         = "terraform-google-modules"
      description = "Installs GitLab on Kubernetes Engine"
    },
    {
      name        = "terraform-google-group"
      org         = "terraform-google-modules"
      description = "Manages Google Groups"
    },
    {
      name        = "terraform-google-gsuite-export"
      org         = "terraform-google-modules"
      description = "Creates a Compute Engine VM instance and sets up a cronjob to export GSuite Admin SDK data to Cloud Logging on a schedule"
    },
    {
      name        = "terraform-google-iam"
      org         = "terraform-google-modules"
      description = "Manages multiple IAM roles for resources on Google Cloud"
    },
    {
      name        = "terraform-google-jenkins"
      org         = "terraform-google-modules"
      description = "Creates a Compute Engine instance running Jenkins"
    },
    {
      name        = "terraform-google-kms"
      org         = "terraform-google-modules"
      description = "Allows managing a keyring, zero or more keys in the keyring, and IAM role bindings on individual keys"
      topics      = "cft-fabric"
    },
    {
      name        = "terraform-google-lb"
      org         = "terraform-google-modules"
      description = "Creates a regional TCP proxy load balancer for Compute Engine by using target pools and forwarding rules"
    },
    {
      name        = "terraform-google-lb-http"
      org         = "terraform-google-modules"
      description = "Creates a global HTTP load balancer for Compute Engine by using forwarding rules"
    },
    {
      name        = "terraform-google-lb-internal"
      org         = "terraform-google-modules"
      description = "Creates an internal load balancer for Compute Engine by using forwarding rules"
    },
    {
      name        = "terraform-google-log-export"
      org         = "terraform-google-modules"
      description = "Creates log exports at the project, folder, or organization level"
      topics      = "cft-fabric"
    },
    {
      name        = "terraform-google-memorystore"
      org         = "terraform-google-modules"
      description = "Creates a fully functional Google Memorystore (redis) instance"
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
      topics      = "cft-fabric"
    },
    {
      name        = "terraform-google-org-policy"
      org         = "terraform-google-modules"
      description = "Manages Google Cloud organization policies"
      topics      = "cft-fabric"
    },
    {
      name        = "terraform-google-project-factory"
      org         = "terraform-google-modules"
      description = "Creates an opinionated Google Cloud project by using Shared VPC, IAM, and Google Cloud APIs"
      topics      = "cft-fabric"
    },
    {
      name        = "terraform-google-pubsub"
      org         = "terraform-google-modules"
      description = "Creates Pub/Sub topic and subscriptions associated with the topic"
    },
    {
      name        = "terraform-google-sap"
      org         = "terraform-google-modules"
      description = "Deploys SAP products"
    },
    {
      name        = "terraform-google-scheduled-function"
      org         = "terraform-google-modules"
      description = "Sets up a scheduled job to trigger events and run functions"
    },
    {
      name        = "terraform-google-service-accounts"
      org         = "terraform-google-modules"
      description = "Creates one or more service accounts and grants them basic roles"
      topics      = "cft-fabric"
    },
    {
      name        = "terraform-google-slo"
      org         = "terraform-google-modules"
      description = "Creates SLOs on Google Cloud from custom Stackdriver metrics capability to export SLOs to Google Cloud services and other systems"
    },
    {
      name        = "terraform-google-sql-db"
      org         = "terraform-google-modules"
      description = "Creates a Cloud SQL database instance"
    },
    {
      name        = "terraform-google-startup-scripts"
      org         = "terraform-google-modules"
      description = "Provides a library of useful startup scripts to embed in VMs"
    },
    {
      name        = "terraform-google-utils"
      org         = "terraform-google-modules"
      description = "Gets the short names for a given Google Cloud region"
    },
    {
      name        = "terraform-google-vault"
      org         = "terraform-google-modules"
      description = "Deploys Vault on Compute Engine"
      topics      = "hashicorp-vault,terraform,terraform-module"
    },
    {
      name        = "terraform-google-vm"
      org         = "terraform-google-modules"
      description = "Provisions VMs in Google Cloud"
    },
    {
      name        = "terraform-google-vpc-service-controls"
      org         = "terraform-google-modules"
      description = "Handles opinionated VPC Service Controls and Access Context Manager configuration and deployments"
    },
    {
      name        = "terraform-google-vpn"
      org         = "terraform-google-modules"
      description = "Sets up a Cloud VPN gateway"
      topics      = "cft-fabric"
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
    }
  ]
}
