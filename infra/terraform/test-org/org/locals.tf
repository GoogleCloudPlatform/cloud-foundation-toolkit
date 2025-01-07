/**
 * Copyright 2019-2024 Google LLC
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
  org_id              = "943740911108"
  old_billing_account = "01D904-DAF6EC-F34EF7"
  billing_account     = "0138EF-C93849-98B0B5"
  lr_billing_account  = "01108A-537F1E-A5BFFC"
  cft_ci_group        = "cft-ci-robots@test.blueprints.joonix.net"
  gcp_admins_group    = "gcp-admins@test.blueprints.joonix.net"
  ci_project_id       = "cloud-foundation-cicd"

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
  jss_common_group = "jump-start-solutions-admins"

  adc_common_admins = ["q2w"]

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
 *  name              = "string" (required for modules)
 *  short_name        = "string" (optional for modules, if not prefixed with 'terraform-google-')
 *  org               = "terraform-google-modules" or "GoogleCloudPlatform" (required)
 *  description       = "string" (required)
 *  maintainers       = "list(string)" ["user1", "user2", "CASE SENSATIVE"] (optional)
 *  admins            = "list(string)" ["user1", "user2", "CASE SENSATIVE"] (optional)
 *  groups            = "list(string)" ["group1", "group1"] (optional)
 *  homepage_url      = "string" (optional, overrides default)
 *  module            = BOOL (optional, default is true which includes GH repo configuration)
 *  topics            = "string1,string2,string3" (one or more of local.common_topics required if module = true)
 *  lint_env          = "map(string)" (optional)
 *  disable_lint_yaml = BOOL (optional, default is true)
 *  enable_periodic   = BOOL (optional, if enabled runs a daily periodic test. Defaults to false )
 *
 */

  repos = [
    {
      name        = "cloud-foundation-training"
      org         = "terraform-google-modules"
      description = ""
      maintainers = ["marine675"]
    },
    {
      name        = "terraform-google-healthcare"
      org         = "terraform-google-modules"
      description = "Handles opinionated Google Cloud Healthcare datasets and stores"
      maintainers = ["yeweidaniel"]
      topics      = local.common_topics.hcls
    },
    {
      name        = "terraform-google-cloud-run"
      org         = "GoogleCloudPlatform"
      description = "Deploys apps to Cloud Run, along with option to map custom domain"
      maintainers = concat(["prabhu34", "anamer", "gtsorbo"], local.adc_common_admins)
      topics      = "cloudrun,google-cloud-platform,terraform-modules,${local.common_topics.serverless}"
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name        = "terraform-google-secured-data-warehouse"
      org         = "GoogleCloudPlatform"
      description = "Deploys a secured BigQuery data warehouse"
      maintainers = ["lanre-OG"]
      topics      = join(",", [local.common_topics.da, local.common_topics.e2e])
      lint_env = {
        SHELLCHECK_OPTS = "-e SC2154 -e SC2171 -e SC2086"
      }
    },
    {
      name        = "terraform-google-anthos-vm"
      org         = "GoogleCloudPlatform"
      description = "Creates VMs on Anthos Bare Metal clusters"
      maintainers = ["zhuchenwang"]
      topics      = "anthos,kubernetes,terraform-module,vm,${local.common_topics.compute}"
    },
    {
      name        = "terraform-google-kubernetes-engine"
      org         = "terraform-google-modules"
      description = "Configures opinionated GKE clusters"
      maintainers = ["ericyz"]
      admins      = ["apeabody"]
      topics      = join(",", [local.common_topics.compute, local.common_topics.containers])
    },
    {
      name            = "terraform-ecommerce-microservices-on-gke"
      short_name      = "ecommerce-microservices"
      org             = "GoogleCloudPlatform"
      description     = "Deploys a web-based ecommerce app into a multi-cluster Google Kubernetes Engine setup."
      groups          = ["dee-platform-ops", local.jss_common_group]
      enable_periodic = true
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name        = "terraform-example-java-dynamic-point-of-sale"
      short_name  = "java-dynamic-point-of-sale"
      org         = "GoogleCloudPlatform"
      description = "Deploys a dynamic Java webapp into a Google Kubernetes Engine cluster."
      maintainers = ["Shabirmean", "Mukamik"]
      groups      = ["dee-platform-ops", local.jss_common_group]
      lint_env = {
        "EXCLUDE_HEADER_CHECK" = "\\./infra/sql-schema"
      }
      enable_periodic = true
    },
    {
      name         = "terraform-example-foundation"
      short_name   = "example-foundation"
      org          = "terraform-google-modules"
      description  = "Shows how the CFT modules can be composed to build a secure cloud foundation"
      maintainers  = ["rjerrems", "gtsorbo", "eeaton", "sleighton2022"]
      homepage_url = "https://cloud.google.com/architecture/security-foundations"
      topics       = join(",", [local.common_topics.e2e, local.common_topics.ops])
      lint_env = {
        "EXCLUDE_LINT_DIRS" = "\\./3-networks/modules/transitivity/assets",
        "ENABLE_PARALLEL"   = "0",
        "DISABLE_TFLINT"    = "1"
      }
    },
    {
      name            = "terraform-google-log-analysis"
      org             = "GoogleCloudPlatform"
      description     = "Stores and analyzes log data"
      maintainers     = ["ryotat7"]
      topics          = local.common_topics.da
      groups          = [local.jss_common_group]
      enable_periodic = true
    },
    {
      name            = "terraform-google-three-tier-web-app"
      org             = "GoogleCloudPlatform"
      description     = "Deploys a three tier web application using Cloud Run and Cloud SQL"
      maintainers     = ["tpryan"]
      topics          = join(",", [local.common_topics.serverless, local.common_topics.db])
      groups          = [local.jss_common_group]
      enable_periodic = true
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name        = "terraform-google-load-balanced-vms"
      org         = "GoogleCloudPlatform"
      description = "Creates a Managed Instance Group with a loadbalancer"
      maintainers = ["tpryan"]
      topics      = local.common_topics.net
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name            = "terraform-google-secure-cicd"
      org             = "GoogleCloudPlatform"
      description     = "Builds a secure CI/CD pipeline on Google Cloud"
      maintainers     = ["gtsorbo"]
      topics          = join(",", [local.common_topics.security, local.common_topics.devtools, local.common_topics.e2e])
      enable_periodic = true
      groups          = [local.jss_common_group]
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name            = "terraform-google-media-cdn-vod"
      org             = "GoogleCloudPlatform"
      description     = "Deploys Media CDN video-on-demand"
      maintainers     = ["roddzurcher"]
      topics          = local.common_topics.ops
      groups          = [local.jss_common_group]
      enable_periodic = true
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
      maintainers = ["gtsorbo"]
      topics      = local.common_topics.net
    },
    {
      name        = "terraform-google-secret-manager"
      org         = "GoogleCloudPlatform"
      description = "Creates one or more Google Secret Manager secrets and manages basic permissions for them"
      maintainers = local.adc_common_admins
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
      maintainers = ["davenportjw", "shanecglass"]
      groups      = [local.jss_common_group]
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name        = "terraform-google-bootstrap"
      org         = "terraform-google-modules"
      description = "Bootstraps Terraform usage and related CI/CD in a new Google Cloud organization"
      topics      = join(",", [local.common_topics.ops, local.common_topics.devtools])
      maintainers = ["josephdt12"]
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
      maintainers = ["imrannayer"]
    },
    {
      name        = "terraform-google-cloud-nat"
      org         = "terraform-google-modules"
      description = "Creates and configures Cloud NAT"
      topics      = local.common_topics.net
      maintainers = ["imrannayer"]
    },
    {
      name        = "terraform-google-cloud-operations"
      org         = "terraform-google-modules"
      description = "Manages Cloud Logging and Cloud Monitoring"
      topics      = local.common_topics.ops
      maintainers = ["imrannayer"]
      groups      = ["stackdriver-committers"]
    },
    {
      name        = "terraform-google-cloud-router"
      org         = "terraform-google-modules"
      description = "Manages a Cloud Router on Google Cloud"
      topics      = local.common_topics.net
      maintainers = ["imrannayer"]
    },
    {
      name        = "terraform-google-cloud-storage"
      org         = "terraform-google-modules"
      description = "Creates one or more Cloud Storage buckets and assigns basic permissions on them to arbitrary users"
      topics      = local.common_topics.storage
      maintainers = local.adc_common_admins
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name        = "terraform-google-composer"
      org         = "terraform-google-modules"
      description = "Manages Cloud Composer v1 and v2 along with option to manage networking"
      topics      = join(",", [local.common_topics.da, local.common_topics.ops])
      maintainers = ["imrannayer"]
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
      lint_env    = { "EXCLUDE_LINT_DIRS" = "\\./cache" }
    },
    {
      name        = "terraform-google-github-actions-runners"
      org         = "terraform-google-modules"
      description = "Creates self-hosted GitHub Actions Runners on Google Cloud"
      topics      = local.common_topics.devtools
      maintainers = ["gtsorbo"]
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
      maintainers = ["imrannayer"]
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
      maintainers = ["imrannayer"]
    },
    {
      name        = "terraform-google-lb-http"
      org         = "terraform-google-modules"
      description = "Creates a global HTTP load balancer for Compute Engine by using forwarding rules"
      topics      = local.common_topics.net
      maintainers = concat(["imrannayer"], local.adc_common_admins)
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name        = "terraform-google-lb-internal"
      org         = "terraform-google-modules"
      description = "Creates an internal load balancer for Compute Engine by using forwarding rules"
      topics      = local.common_topics.net
      maintainers = ["imrannayer"]
    },
    {
      name        = "terraform-google-log-export"
      org         = "terraform-google-modules"
      description = "Creates log exports at the project, folder, or organization level"
      topics      = local.common_topics.ops
      maintainers = ["imrannayer"]
    },
    {
      name        = "terraform-google-memorystore"
      org         = "terraform-google-modules"
      description = "Creates a fully functional Google Memorystore (redis) instance"
      topics      = local.common_topics.db
      maintainers = concat(["imrannayer"], local.adc_common_admins)
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name              = "terraform-google-module-template"
      org               = "terraform-google-modules"
      description       = "Provides a template for creating a Cloud Foundation Toolkit Terraform module"
      disable_lint_yaml = true
    },
    {
      name        = "terraform-google-network"
      org         = "terraform-google-modules"
      description = "Sets up a new VPC network on Google Cloud"
      topics      = local.common_topics.net
      maintainers = ["imrannayer"]
    },
    {
      name        = "terraform-google-org-policy"
      org         = "terraform-google-modules"
      description = "Manages Google Cloud organization policies"
      topics      = local.common_topics.security
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name        = "terraform-google-project-factory"
      org         = "terraform-google-modules"
      description = "Creates an opinionated Google Cloud project by using Shared VPC, IAM, and Google Cloud APIs"
      topics      = local.common_topics.ops
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name        = "terraform-google-pubsub"
      org         = "terraform-google-modules"
      description = "Creates Pub/Sub topic and subscriptions associated with the topic"
      topics      = local.common_topics.da
      maintainers = ["imrannayer"]
    },
    {
      name        = "terraform-google-sap"
      org         = "terraform-google-modules"
      description = "Deploys SAP products"
      topics      = local.common_topics.compute
      maintainers = ["sjswerdlow", "megelatim"]
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
      maintainers = local.adc_common_admins
      topics      = local.common_topics.security
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
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
      maintainers = concat(["isaurabhuttam", "imrannayer"], local.adc_common_admins)
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
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
      maintainers = concat(["erlanderlo"], local.adc_common_admins)
      topics      = local.common_topics.compute
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name        = "terraform-google-vpc-service-controls"
      org         = "terraform-google-modules"
      description = "Handles opinionated VPC Service Controls and Access Context Manager configuration and deployments"
      topics      = local.common_topics.net
      maintainers = ["imrannayer"]
    },
    {
      name        = "terraform-google-vpn"
      org         = "terraform-google-modules"
      description = "Sets up a Cloud VPN gateway"
      topics      = local.common_topics.net
      maintainers = ["imrannayer"]
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
      short_name      = "docs-samples"
      org             = "terraform-google-modules"
      module          = false
      enable_periodic = true
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
      description = "Deploys the WAAP solution on Google Cloud"
      maintainers = ["gtsorbo"]
      topics      = local.common_topics.ops
    },
    {
      name        = "terraform-google-cloud-workflows"
      org         = "GoogleCloudPlatform"
      description = "Manage Workflows with optional Scheduler or Event Arc triggers"
      maintainers = ["anaik91"]
      topics      = join(",", [local.common_topics.serverless, local.common_topics.devtools])
    },
    {
      name        = "terraform-google-vertex-ai"
      org         = "GoogleCloudPlatform"
      description = "Deploy Vertex AI resources"
      maintainers = ["imrannayer"]
      topics      = join(",", [local.common_topics.compute])
    },
    {
      name        = "terraform-google-cloud-armor"
      org         = "GoogleCloudPlatform"
      description = "Deploy Cloud Armor security policy"
      maintainers = ["imrannayer"]
      topics      = join(",", [local.common_topics.compute, local.common_topics.net])
    },
    {
      name        = "terraform-google-pam"
      org         = "GoogleCloudPlatform"
      description = "Deploy Privileged Access Manager"
      maintainers = ["imrannayer", "mgaur10"]
      topics      = local.common_topics.security
    },
    {
      name        = "terraform-google-netapp-volumes"
      org         = "GoogleCloudPlatform"
      description = "Deploy NetApp Storage Volumes"
      maintainers = ["imrannayer"]
      topics      = join(",", [local.common_topics.compute, local.common_topics.net])
    },
    {
      name        = "terraform-google-cloud-deploy"
      org         = "GoogleCloudPlatform"
      description = "Create Cloud Deploy pipelines and targets"
      maintainers = ["gtsorbo", "niranjankl"]
      topics      = join(",", [local.common_topics.devtools])
    },
    {
      name        = "terraform-google-cloud-functions"
      org         = "GoogleCloudPlatform"
      description = "Deploys Cloud Functions (Gen 2)"
      maintainers = ["prabhu34", "gtsorbo"]
      topics      = "cloudfunctions,functions,google-cloud-platform,terraform-modules,${local.common_topics.serverless}"
    },
    {
      name            = "terraform-dynamic-python-webapp"
      short_name      = "dynamic-python-webapp"
      org             = "GoogleCloudPlatform"
      description     = "Deploy a dynamic python webapp"
      maintainers     = ["glasnt", "donmccasland"]
      homepage_url    = "avocano.dev"
      groups          = [local.jss_common_group, "team-egg"]
      enable_periodic = true
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name            = "terraform-dynamic-javascript-webapp"
      short_name      = "dynamic-javascript-webapp"
      org             = "GoogleCloudPlatform"
      description     = "Deploy a dynamic javascript webapp"
      maintainers     = ["LukeSchlangen", "donmccasland"]
      homepage_url    = "avocano.dev"
      groups          = [local.jss_common_group, "team-egg", "developer-journey-app-approvers"]
      enable_periodic = true
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name            = "terraform-example-deploy-java-multizone"
      short_name      = "deploy-java-multizone"
      org             = "GoogleCloudPlatform"
      description     = "Deploy a multizone Java application"
      maintainers     = ["donmccasland"]
      groups          = [local.jss_common_group]
      enable_periodic = false
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name        = "terraform-google-itar-architectures"
      org         = "GoogleCloudPlatform"
      description = "Includes use cases for deploying ITAR-aligned architectures on Google Cloud"
      maintainers = ["gtsorbo"]
      topics      = join(",", [local.common_topics.compute], ["compliance"])
    },
    {
      name            = "terraform-google-analytics-lakehouse"
      org             = "GoogleCloudPlatform"
      description     = "Deploys a Lakehouse Architecture Solution"
      maintainers     = ["davenportjw", "bradmiro"]
      topics          = local.common_topics.da
      groups          = [local.jss_common_group]
      enable_periodic = true
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name        = "terraform-google-alloy-db"
      org         = "GoogleCloudPlatform"
      description = "Creates an Alloy DB instance"
      maintainers = ["anaik91", "imrannayer"]
      topics      = local.common_topics.db
    },
    {
      name        = "terraform-google-cloud-ids"
      org         = "GoogleCloudPlatform"
      description = "Deploys a Cloud IDS instance and associated resources."
      maintainers = ["gtsorbo", "mgaur10"]
      topics      = join(",", [local.common_topics.security, local.common_topics.net])
    },
    {
      name            = "terraform-example-deploy-java-gke"
      short_name      = "deploy-java-gke"
      org             = "GoogleCloudPlatform"
      description     = "Deploy a Legacy Java App GKE"
      groups          = ["dee-platform-ops", local.jss_common_group]
      enable_periodic = true
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }

    },
    {
      name        = "terraform-google-crmint"
      org         = "GoogleCloudPlatform"
      description = "Deploy the marketing analytics application, CRMint"
      maintainers = ["dulacp"]
      topics      = join(",", [local.common_topics.da, local.common_topics.e2e], ["marketing"])
    },
    {
      name            = "terraform-ml-image-annotation-gcf"
      short_name      = "ml-image-annotation-gcf"
      org             = "GoogleCloudPlatform"
      description     = "Deploys an app for ml image annotation using gcf"
      maintainers     = ["xsxm", "ivanmkc", "balajismaniam", "donmccasland"]
      groups          = ["dee-data-ai", local.jss_common_group]
      enable_periodic = true
    },
    {
      name        = "terraform-google-out-of-band-security"
      org         = "GoogleCloudPlatform"
      description = "Creates a 3P out-of-band security appliance deployment"
      maintainers = ["Saipriyavk", "ChrisBarefoot"]
      topics      = local.common_topics.net
    },
    {
      name        = "notebooks-blueprint-security"
      short_name  = "secured-notebook"
      org         = "GoogleCloudPlatform"
      description = "Opinionated setup for securely using AI Platform Notebooks."
      maintainers = ["gtsorbo", "erlanderlo"]
      topics      = join(",", [local.common_topics.da, local.common_topics.security])
    },
    {
      name            = "terraform-genai-doc-summarization"
      short_name      = "genai-doc-summarization"
      org             = "GoogleCloudPlatform"
      description     = "Summarizes document using OCR and Vertex Generative AI LLM"
      maintainers     = ["asrivas", "davidcavazos"]
      groups          = [local.jss_common_group]
      enable_periodic = true
    },
    {
      name            = "terraform-genai-knowledge-base"
      short_name      = "genai-knowledge-base"
      org             = "GoogleCloudPlatform"
      description     = "Fine tune an LLM model to answer questions from your documents."
      maintainers     = ["davidcavazos"]
      groups          = [local.jss_common_group]
      enable_periodic = true
    },
    {
      name        = "terraform-google-secured-data-warehouse-onprem-ingest"
      short_name  = "sdw-onprem-ingest"
      org         = "GoogleCloudPlatform"
      description = "Deploys a secured data warehouse variant for ingesting encrypted data from on-prem sources"
      maintainers = ["lanre-OG"]
      topics      = join(",", [local.common_topics.da, local.common_topics.security, local.common_topics.e2e])
    },
    {
      name        = "terraform-google-tf-cloud-agents"
      org         = "GoogleCloudPlatform"
      description = "Creates self-hosted Terraform Cloud Agent on Google Cloud"
      topics      = join(",", [local.common_topics.ops, local.common_topics.security, local.common_topics.devtools])
    },
    {
      name        = "terraform-google-cloud-spanner"
      org         = "GoogleCloudPlatform"
      description = "Deploy Spanner instances"
      maintainers = ["anaik91", "imrannayer"]
      topics      = local.common_topics.db
    },
    {
      name            = "terraform-pubsub-integration-golang"
      org             = "GoogleCloudPlatform"
      short_name      = "pubsub-golang-app"
      maintainers     = ["Shabirmean", "Mukamik"]
      groups          = ["dee-platform-ops", local.jss_common_group]
      enable_periodic = true
    },
    {
      name            = "terraform-pubsub-integration-java"
      org             = "GoogleCloudPlatform"
      short_name      = "pubsub-java-app"
      maintainers     = ["Shabirmean", "Mukamik"]
      groups          = ["dee-platform-ops", local.jss_common_group]
      enable_periodic = true
    },
    {
      name        = "terraform-google-backup-dr"
      org         = "GoogleCloudPlatform"
      short_name  = "backup-dr"
      description = "Deploy Backup and DR appliances"
      maintainers = ["umeshkumhar"]
      topics      = join(",", [local.common_topics.compute, local.common_topics.ops])
    },
    {
      name        = "terraform-google-tags"
      org         = "GoogleCloudPlatform"
      description = "Create and manage Google Cloud Tags"
      maintainers = ["nidhi0710"]
      topics      = join(",", [local.common_topics.security, local.common_topics.ops])
    },
    {
      name        = "terraform-google-dataplex-auto-data-quality"
      org         = "GoogleCloudPlatform"
      description = "Move data between environments using Dataplex"
      maintainers = ["bradmiro"]
      topics      = local.common_topics.da
    },
    {
      name        = "terraform-google-enterprise-application"
      org         = "GoogleCloudPlatform"
      description = "Deploy an enterprise developer platform on Google Cloud"
      maintainers = ["gtsorbo", "erictune", "yliaog", "sleighton2022", "apeabody"]
      topics      = join(",", [local.common_topics.e2e, local.common_topics.ops])
    },
    {
      name            = "terraform-genai-rag"
      short_name      = "genai-rag"
      org             = "GoogleCloudPlatform"
      description     = "Deploys a Generative AI RAG solution"
      maintainers     = ["davenportjw", "bradmiro"]
      groups          = ["dee-platform-ops", "dee-data-ai", local.jss_common_group]
      enable_periodic = true
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name        = "terraform-google-artifact-registry"
      org         = "GoogleCloudPlatform"
      description = "Create and manage Artifact Registry repositories"
      maintainers = ["prabhu34"]
      topics      = join(",", [local.common_topics.containers, local.common_topics.devtools])
    },
    {
      name        = "terraform-google-bigtable"
      org         = "GoogleCloudPlatform"
      description = "Create and manage Google Bigtable resources"
      maintainers = ["hariprabhaam"]
      topics      = local.common_topics.da
    },
    {
      name        = "terraform-google-secure-web-proxy"
      org         = "GoogleCloudPlatform"
      description = "Create and manage Secure Web Proxy on GCP for secured egress web traffic"
      maintainers = ["maitreya-source"]
      topics      = join(",", [local.common_topics.security, local.common_topics.net])
    },
    {
      name            = "terraform-cloud-client-api"
      short_name      = "cloud-client-api"
      org             = "GoogleCloudPlatform"
      description     = "Deploys an example application that uses Cloud Client APIs"
      maintainers     = ["glasnt", "iennae"]
      groups          = ["team-egg", local.jss_common_group]
      enable_periodic = true
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name        = "kms-solutions"
      org         = "GoogleCloudPlatform"
      description = "Store Cloud KMS scripts, artifacts, code samples, and more."
      maintainers = ["tdbhacks", "erlanderlo", "g-swap", "nb-goog"]
      lint_env = {
        ENABLE_BPMETADATA = "1"
      }
    },
    {
      name        = "terraform-dataanalytics-eventdriven"
      short_name  = "dataanalytics-eventdriven"
      org         = "GoogleCloudPlatform"
      description = "Uses click-to-deploy to demonstrate how to load data from Cloud Storage to BigQuery using an event-driven load function."
      groups      = [local.jss_common_group]
      maintainers = ["fellipeamedeiros", "sylvioneto"]
    },
    {
      name        = "terraform-google-apphub"
      org         = "GoogleCloudPlatform"
      description = "Creates and manages AppHub resources"
      maintainers = ["q2w"]
      admins      = ["bharathkkb"]
    },
  ]
}
