# Terraform Modules
The Cloud Foundation Toolkit includes over **41** Terraform modules.

* [address](https://github.com/terraform-google-modules/terraform-google-address) - A Terraform module for managing Google Cloud IP addresses.
* [bastion-host](https://github.com/terraform-google-modules/terraform-google-bastion-host) - This module will generate a bastion host vm compatible with OS Login and IAP Tunneling that can be used to access internal VMs.
  * [bastion-group](https://github.com/terraform-google-modules/terraform-google-bastion-host/tree/master/modules/bastion-group)
  * [iap-tunneling](https://github.com/terraform-google-modules/terraform-google-bastion-host/tree/master/modules/iap-tunneling)
* [bigquery](https://github.com/terraform-google-modules/terraform-google-bigquery) - This module allows you to create opinionated Google Cloud Platform BigQuery datasets and tables.
  * [authorization](https://github.com/terraform-google-modules/terraform-google-bigquery/tree/master/modules/authorization)
  * [udf](https://github.com/terraform-google-modules/terraform-google-bigquery/tree/master/modules/udf)
* [bootstrap](https://github.com/terraform-google-modules/terraform-google-bootstrap) - A module for bootstrapping Terraform usage in a new GCP organization.
  * [cloudbuild](https://github.com/terraform-google-modules/terraform-google-bootstrap/tree/master/modules/cloudbuild)
* [cloud-datastore](https://github.com/terraform-google-modules/terraform-google-cloud-datastore) - A Terraform module to help you to manage Google Cloud Datastore.
* [cloud-dns](https://github.com/terraform-google-modules/terraform-google-cloud-dns) - This module makes it easy to create and manage Google Cloud DNS public or private zones, and their records. https://registry.terraform.io/modules/terraform-google-modules/cloud-dns/google/
* [cloud-nat](https://github.com/terraform-google-modules/terraform-google-cloud-nat) - This module handles opinionated Google Cloud Platform Cloud NAT creation and configuration.
* [cloud-router](https://github.com/terraform-google-modules/terraform-google-cloud-router) - Manage a Cloud Router on GCP
  * [interconnect_attachment](https://github.com/terraform-google-modules/terraform-google-cloud-router/tree/master/modules/interconnect_attachment)
  * [interface](https://github.com/terraform-google-modules/terraform-google-cloud-router/tree/master/modules/interface)
* [cloud-storage](https://github.com/terraform-google-modules/terraform-google-cloud-storage) - This module makes it easy to create one or more GCS buckets, and assign basic permissions on them to arbitrary users.
  * [simple_bucket](https://github.com/terraform-google-modules/terraform-google-cloud-storage/tree/master/modules/simple_bucket)
* [container-vm](https://github.com/terraform-google-modules/terraform-google-container-vm) - This module simplifies deploying containers on GCE instances.
  * [cos-coredns](https://github.com/terraform-google-modules/terraform-google-container-vm/tree/master/modules/cos-coredns)
  * [cos-generic](https://github.com/terraform-google-modules/terraform-google-container-vm/tree/master/modules/cos-generic)
  * [cos-mysql](https://github.com/terraform-google-modules/terraform-google-container-vm/tree/master/modules/cos-mysql)
* [dataflow](https://github.com/terraform-google-modules/terraform-google-dataflow) - This module handles opiniated Dataflow job configuration and deployments.
  * [dataflow_bucket](https://github.com/terraform-google-modules/terraform-google-dataflow/tree/master/modules/dataflow_bucket)
* [datalab](https://github.com/terraform-google-modules/terraform-google-datalab) - This module will create DataLab instances with support for GPU instances. 
  * [iap_firewall](https://github.com/terraform-google-modules/terraform-google-datalab/tree/master/modules/iap_firewall)
  * [instance](https://github.com/terraform-google-modules/terraform-google-datalab/tree/master/modules/instance)
  * [template_files](https://github.com/terraform-google-modules/terraform-google-datalab/tree/master/modules/template_files)
* [endpoints-dns](https://github.com/terraform-google-modules/terraform-google-endpoints-dns) - 
* [event-function](https://github.com/terraform-google-modules/terraform-google-event-function) - Terraform module for responding to logging events with a function
  * [event-folder-log-entry](https://github.com/terraform-google-modules/terraform-google-event-function/tree/master/modules/event-folder-log-entry)
  * [event-project-log-entry](https://github.com/terraform-google-modules/terraform-google-event-function/tree/master/modules/event-project-log-entry)
  * [repository-function](https://github.com/terraform-google-modules/terraform-google-event-function/tree/master/modules/repository-function)
* [folders](https://github.com/terraform-google-modules/terraform-google-folders) - This module helps create several folders under the same parent
* [forseti](https://github.com/terraform-google-modules/terraform-google-forseti) - A Terraform module for installing Forseti on GCP
  * [client](https://github.com/terraform-google-modules/terraform-google-forseti/tree/master/modules/client)
  * [client_config](https://github.com/terraform-google-modules/terraform-google-forseti/tree/master/modules/client_config)
  * [client_gcs](https://github.com/terraform-google-modules/terraform-google-forseti/tree/master/modules/client_gcs)
  * [client_iam](https://github.com/terraform-google-modules/terraform-google-forseti/tree/master/modules/client_iam)
  * [cloudsql](https://github.com/terraform-google-modules/terraform-google-forseti/tree/master/modules/cloudsql)
  * [on_gke](https://github.com/terraform-google-modules/terraform-google-forseti/tree/master/modules/on_gke)
  * [real_time_enforcer](https://github.com/terraform-google-modules/terraform-google-forseti/tree/master/modules/real_time_enforcer)
  * [real_time_enforcer_organization_sink](https://github.com/terraform-google-modules/terraform-google-forseti/tree/master/modules/real_time_enforcer_organization_sink)
  * [real_time_enforcer_project_sink](https://github.com/terraform-google-modules/terraform-google-forseti/tree/master/modules/real_time_enforcer_project_sink)
  * [real_time_enforcer_roles](https://github.com/terraform-google-modules/terraform-google-forseti/tree/master/modules/real_time_enforcer_roles)
  * [rules](https://github.com/terraform-google-modules/terraform-google-forseti/tree/master/modules/rules)
  * [server](https://github.com/terraform-google-modules/terraform-google-forseti/tree/master/modules/server)
  * [server_config](https://github.com/terraform-google-modules/terraform-google-forseti/tree/master/modules/server_config)
  * [server_gcs](https://github.com/terraform-google-modules/terraform-google-forseti/tree/master/modules/server_gcs)
  * [server_iam](https://github.com/terraform-google-modules/terraform-google-forseti/tree/master/modules/server_iam)
* [gcloud](https://github.com/terraform-google-modules/terraform-google-gcloud) - A module for executing gcloud commands within Terraform.
* [gke-gitlab](https://github.com/terraform-google-modules/terraform-google-gke-gitlab) - Installs GitLab on Kubernetes Engine
* [gsuite-export](https://github.com/terraform-google-modules/terraform-google-gsuite-export) - 
* [healthcare](https://github.com/terraform-google-modules/terraform-google-healthcare) - This module handles opinionated Google Cloud Platform Healthcare datasets and stores.
* [iam](https://github.com/terraform-google-modules/terraform-google-iam) - This Terraform module makes it easier to non-destructively manage multiple IAM roles for resources on Google Cloud Platform.
  * [audit_config](https://github.com/terraform-google-modules/terraform-google-iam/tree/master/modules/audit_config)
  * [billing_accounts_iam](https://github.com/terraform-google-modules/terraform-google-iam/tree/master/modules/billing_accounts_iam)
  * [custom_role_iam](https://github.com/terraform-google-modules/terraform-google-iam/tree/master/modules/custom_role_iam)
  * [folders_iam](https://github.com/terraform-google-modules/terraform-google-iam/tree/master/modules/folders_iam)
  * [helper](https://github.com/terraform-google-modules/terraform-google-iam/tree/master/modules/helper)
  * [kms_crypto_keys_iam](https://github.com/terraform-google-modules/terraform-google-iam/tree/master/modules/kms_crypto_keys_iam)
  * [kms_key_rings_iam](https://github.com/terraform-google-modules/terraform-google-iam/tree/master/modules/kms_key_rings_iam)
  * [member_iam](https://github.com/terraform-google-modules/terraform-google-iam/tree/master/modules/member_iam)
  * [organizations_iam](https://github.com/terraform-google-modules/terraform-google-iam/tree/master/modules/organizations_iam)
  * [projects_iam](https://github.com/terraform-google-modules/terraform-google-iam/tree/master/modules/projects_iam)
  * [pubsub_subscriptions_iam](https://github.com/terraform-google-modules/terraform-google-iam/tree/master/modules/pubsub_subscriptions_iam)
  * [pubsub_topics_iam](https://github.com/terraform-google-modules/terraform-google-iam/tree/master/modules/pubsub_topics_iam)
  * [service_accounts_iam](https://github.com/terraform-google-modules/terraform-google-iam/tree/master/modules/service_accounts_iam)
  * [storage_buckets_iam](https://github.com/terraform-google-modules/terraform-google-iam/tree/master/modules/storage_buckets_iam)
  * [subnets_iam](https://github.com/terraform-google-modules/terraform-google-iam/tree/master/modules/subnets_iam)
* [jenkins](https://github.com/terraform-google-modules/terraform-google-jenkins) - 
  * [artifact_storage](https://github.com/terraform-google-modules/terraform-google-jenkins/tree/master/modules/artifact_storage)
* [kms](https://github.com/terraform-google-modules/terraform-google-kms) - Simple Cloud KMS module that allows managing a keyring, zero or more keys in the keyring, and IAM role bindings on individual keys.
* [kubernetes-engine](https://github.com/terraform-google-modules/terraform-google-kubernetes-engine) - A Terraform module for configuring GKE clusters.
  * [acm](https://github.com/terraform-google-modules/terraform-google-kubernetes-engine/tree/master/modules/acm)
  * [auth](https://github.com/terraform-google-modules/terraform-google-kubernetes-engine/tree/master/modules/auth)
  * [beta-private-cluster](https://github.com/terraform-google-modules/terraform-google-kubernetes-engine/tree/master/modules/beta-private-cluster)
  * [beta-private-cluster-update-variant](https://github.com/terraform-google-modules/terraform-google-kubernetes-engine/tree/master/modules/beta-private-cluster-update-variant)
  * [beta-public-cluster](https://github.com/terraform-google-modules/terraform-google-kubernetes-engine/tree/master/modules/beta-public-cluster)
  * [config-sync](https://github.com/terraform-google-modules/terraform-google-kubernetes-engine/tree/master/modules/config-sync)
  * [k8s-operator-crd-support](https://github.com/terraform-google-modules/terraform-google-kubernetes-engine/tree/master/modules/k8s-operator-crd-support)
  * [private-cluster](https://github.com/terraform-google-modules/terraform-google-kubernetes-engine/tree/master/modules/private-cluster)
  * [private-cluster-update-variant](https://github.com/terraform-google-modules/terraform-google-kubernetes-engine/tree/master/modules/private-cluster-update-variant)
  * [safer-cluster](https://github.com/terraform-google-modules/terraform-google-kubernetes-engine/tree/master/modules/safer-cluster)
  * [safer-cluster-update-variant](https://github.com/terraform-google-modules/terraform-google-kubernetes-engine/tree/master/modules/safer-cluster-update-variant)
  * [services](https://github.com/terraform-google-modules/terraform-google-kubernetes-engine/tree/master/modules/services)
  * [workload-identity](https://github.com/terraform-google-modules/terraform-google-kubernetes-engine/tree/master/modules/workload-identity)
* [log-export](https://github.com/terraform-google-modules/terraform-google-log-export) - This module allows you to create log exports at the project, folder, or organization level.
  * [bigquery](https://github.com/terraform-google-modules/terraform-google-log-export/tree/master/modules/bigquery)
  * [pubsub](https://github.com/terraform-google-modules/terraform-google-log-export/tree/master/modules/pubsub)
  * [storage](https://github.com/terraform-google-modules/terraform-google-log-export/tree/master/modules/storage)
* [memorystore](https://github.com/terraform-google-modules/terraform-google-memorystore) - 
* [network](https://github.com/terraform-google-modules/terraform-google-network) - A Terraform module that makes it easy to set up a new VPC Network in GCP.
  * [fabric-net-firewall](https://github.com/terraform-google-modules/terraform-google-network/tree/master/modules/fabric-net-firewall)
  * [fabric-net-svpc-access](https://github.com/terraform-google-modules/terraform-google-network/tree/master/modules/fabric-net-svpc-access)
  * [network-peering](https://github.com/terraform-google-modules/terraform-google-network/tree/master/modules/network-peering)
  * [routes](https://github.com/terraform-google-modules/terraform-google-network/tree/master/modules/routes)
  * [routes-beta](https://github.com/terraform-google-modules/terraform-google-network/tree/master/modules/routes-beta)
  * [subnets](https://github.com/terraform-google-modules/terraform-google-network/tree/master/modules/subnets)
  * [subnets-beta](https://github.com/terraform-google-modules/terraform-google-network/tree/master/modules/subnets-beta)
  * [vpc](https://github.com/terraform-google-modules/terraform-google-network/tree/master/modules/vpc)
* [org-policy](https://github.com/terraform-google-modules/terraform-google-org-policy) - A Terraform module for managing GCP org policies.
  * [bucket_policy_only](https://github.com/terraform-google-modules/terraform-google-org-policy/tree/master/modules/bucket_policy_only)
  * [domain_restricted_sharing](https://github.com/terraform-google-modules/terraform-google-org-policy/tree/master/modules/domain_restricted_sharing)
  * [restrict_vm_external_ips](https://github.com/terraform-google-modules/terraform-google-org-policy/tree/master/modules/restrict_vm_external_ips)
  * [skip_default_network](https://github.com/terraform-google-modules/terraform-google-org-policy/tree/master/modules/skip_default_network)
* [project-factory](https://github.com/terraform-google-modules/terraform-google-project-factory) - Opinionated Google Cloud Platform project creation and configuration with Shared VPC, IAM, APIs, etc.
  * [app_engine](https://github.com/terraform-google-modules/terraform-google-project-factory/tree/master/modules/app_engine)
  * [budget](https://github.com/terraform-google-modules/terraform-google-project-factory/tree/master/modules/budget)
  * [core_project_factory](https://github.com/terraform-google-modules/terraform-google-project-factory/tree/master/modules/core_project_factory)
  * [fabric-project](https://github.com/terraform-google-modules/terraform-google-project-factory/tree/master/modules/fabric-project)
  * [gsuite_enabled](https://github.com/terraform-google-modules/terraform-google-project-factory/tree/master/modules/gsuite_enabled)
  * [gsuite_group](https://github.com/terraform-google-modules/terraform-google-project-factory/tree/master/modules/gsuite_group)
  * [project_services](https://github.com/terraform-google-modules/terraform-google-project-factory/tree/master/modules/project_services)
  * [shared_vpc](https://github.com/terraform-google-modules/terraform-google-project-factory/tree/master/modules/shared_vpc)
* [pubsub](https://github.com/terraform-google-modules/terraform-google-pubsub) - This module makes it easy to create Google Cloud Pub/Sub topic and subscriptions associated with the topic.
  * [cloudiot](https://github.com/terraform-google-modules/terraform-google-pubsub/tree/master/modules/cloudiot)
* [sap](https://github.com/terraform-google-modules/terraform-google-sap) - This module is a collection of multiple opinionated submodules to deploy SAP Products.
  * [netweaver](https://github.com/terraform-google-modules/terraform-google-sap/tree/master/modules/netweaver)
  * [sap_hana](https://github.com/terraform-google-modules/terraform-google-sap/tree/master/modules/sap_hana)
  * [sap_hana_ha](https://github.com/terraform-google-modules/terraform-google-sap/tree/master/modules/sap_hana_ha)
  * [sap_hana_python](https://github.com/terraform-google-modules/terraform-google-sap/tree/master/modules/sap_hana/sap_hana_python)
* [scheduled-function](https://github.com/terraform-google-modules/terraform-google-scheduled-function) - This modules makes it easy to set up a scheduled job to trigger events/run functions.
  * [project_cleanup](https://github.com/terraform-google-modules/terraform-google-scheduled-function/tree/master/modules/project_cleanup)
* [secret](https://github.com/terraform-google-modules/terraform-google-secret) - 
  * [gcs-object](https://github.com/terraform-google-modules/terraform-google-secret/tree/master/modules/gcs-object)
  * [secret-infrastructure](https://github.com/terraform-google-modules/terraform-google-secret/tree/master/modules/secret-infrastructure)
* [service-accounts](https://github.com/terraform-google-modules/terraform-google-service-accounts) - This module allows easy creation of one or more service accounts, and granting them basic roles.
* [slo](https://github.com/terraform-google-modules/terraform-google-slo) - Create SLOs on GCP from custom Stackdriver metrics. Capability to export SLOs to GCP services and other systems.
  * [slo](https://github.com/terraform-google-modules/terraform-google-slo/tree/master/modules/slo)
  * [slo-pipeline](https://github.com/terraform-google-modules/terraform-google-slo/tree/master/modules/slo-pipeline)
* [startup-scripts](https://github.com/terraform-google-modules/terraform-google-startup-scripts) - A library of useful startup scripts to embed in VMs created by Terraform
* [utils](https://github.com/terraform-google-modules/terraform-google-utils) - This module provides a way to get the shortnames for a given GCP region.
* [vault](https://github.com/terraform-google-modules/terraform-google-vault) - Modular deployment of Vault on Google Compute Engine with Terraform
* [vm](https://github.com/terraform-google-modules/terraform-google-vm) - This is a collection of opinionated submodules that can be used to provision VMs in GCP.
  * [compute_instance](https://github.com/terraform-google-modules/terraform-google-vm/tree/master/modules/compute_instance)
  * [instance_template](https://github.com/terraform-google-modules/terraform-google-vm/tree/master/modules/instance_template)
  * [mig](https://github.com/terraform-google-modules/terraform-google-vm/tree/master/modules/mig)
  * [mig_with_percent](https://github.com/terraform-google-modules/terraform-google-vm/tree/master/modules/mig_with_percent)
  * [preemptible_and_regular_instance_templates](https://github.com/terraform-google-modules/terraform-google-vm/tree/master/modules/preemptible_and_regular_instance_templates)
  * [umig](https://github.com/terraform-google-modules/terraform-google-vm/tree/master/modules/umig)
* [vpc-service-controls](https://github.com/terraform-google-modules/terraform-google-vpc-service-controls) - This module handles opinionated VPC Service Controls and Access Context Manager configuration and deployments.
  * [access_level](https://github.com/terraform-google-modules/terraform-google-vpc-service-controls/tree/master/modules/access_level)
  * [bridge_service_perimeter](https://github.com/terraform-google-modules/terraform-google-vpc-service-controls/tree/master/modules/bridge_service_perimeter)
  * [regular_service_perimeter](https://github.com/terraform-google-modules/terraform-google-vpc-service-controls/tree/master/modules/regular_service_perimeter)
* [vpn](https://github.com/terraform-google-modules/terraform-google-vpn) - A Terraform Module for setting up Google Cloud VPN
  * [vpn_ha](https://github.com/terraform-google-modules/terraform-google-vpn/tree/master/modules/vpn_ha)
