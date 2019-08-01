locals {
  forseti_org_required_roles = [
    "roles/resourcemanager.organizationAdmin",
    "roles/iam.organizationRoleAdmin",         // Permissions to manage/test real-time-enforcer roles
    "roles/logging.configWriter",              // Permissions to create stackdriver log exports
  ]

  forseti_host_project_required_roles = [
    "roles/compute.admin",
  ]

  forseti_project_required_roles = [
    "roles/compute.instanceAdmin",
    "roles/compute.networkAdmin",
    "roles/compute.securityAdmin",
    "roles/iam.serviceAccountAdmin",
    "roles/iam.serviceAccountUser",
    "roles/serviceusage.serviceUsageAdmin",
    "roles/storage.admin",
    "roles/cloudsql.admin",
    "roles/pubsub.admin",
  ]

  forseti_enforcer_project_required_roles = [
    "roles/storage.admin", // Permissions to create GCS buckets that the enforcer will manage
  ]

  forseti_required_apis = [
    "compute.googleapis.com",
    "serviceusage.googleapis.com",
    "cloudresourcemanager.googleapis.com",
  ]
}

// Define a host project for the Forseti shared VPC test suite.
module "forseti-host-project" {
  source  = "terraform-google-modules/project-factory/google"
  version = "2.4.1"

  providers {
    "google"      = "google.phoogle"
    "google-beta" = "google-beta.phoogle"
  }

  billing_account = "${module.variables.phoogle_billing_account}"
  name            = "ci-forseti-host"
  org_id          = "${var.phoogle_org_id}"

  activate_apis    = "${local.forseti_required_apis}"
  credentials_path = "${var.phoogle_credentials_path}"
  folder_id        = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  project_id       = "ci-forseti-host-a1b2"
}

// Define a shared VPC network within the Forseti host project.
module "forseti-host-network-01" {
  source  = "terraform-google-modules/network/google"
  version = "0.8.0"

  providers {
    "google"      = "google.phoogle"
    "google-beta" = "google-beta.phoogle"
  }

  project_id      = "${module.forseti-host-project.project_id}"
  network_name    = "forseti-network"
  shared_vpc_host = "true"

  secondary_ranges {
    forseti-subnetwork = []
  }

  subnets = [
    {
      subnet_name   = "forseti-subnetwork"
      subnet_ip     = "10.128.0.0/20"
      subnet_region = "us-central1"
    },
  ]
}

resource "google_compute_router" "forseti_host" {
  provider = "google.phoogle"

  name    = "forseti-host"
  network = "${module.forseti-host-network-01.network_self_link}"

  bgp {
    asn = "64514"
  }

  region  = "us-central1"
  project = "${module.forseti-host-project.project_id}"
}

resource "google_compute_router_nat" "forseti_host" {
  provider = "google.phoogle"

  name                               = "forseti-host"
  router                             = "${google_compute_router.forseti_host.name}"
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"

  subnetwork {
    name                    = "${module.forseti-host-network-01.subnets_self_links[0]}"
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }

  project = "${module.forseti-host-project.project_id}"
  region  = "${google_compute_router.forseti_host.region}"
}

module "forseti-service-project" {
  source  = "terraform-google-modules/project-factory/google"
  version = "2.4.1"

  providers {
    "google"      = "google.phoogle"
    "google-beta" = "google-beta.phoogle"
  }

  billing_account = "${module.variables.phoogle_billing_account}"
  name            = "ci-forseti"
  org_id          = "${var.phoogle_org_id}"

  activate_apis      = "${local.forseti_required_apis}"
  credentials_path   = "${var.phoogle_credentials_path}"
  folder_id          = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  random_project_id  = "true"
  shared_vpc         = "ci-forseti-host-a1b2"
  shared_vpc_subnets = ["projects/ci-forseti-host-a1b2/regions/us-central1/subnetworks/forseti-subnetwork"]
}

module "forseti-service-network" {
  source  = "terraform-google-modules/network/google"
  version = "0.8.0"

  providers {
    "google"      = "google.phoogle"
    "google-beta" = "google-beta.phoogle"
  }

  network_name = "forseti-network"
  project_id   = "${module.forseti-service-project.project_id}"

  secondary_ranges = {
    forseti-subnetwork = []
  }

  subnets = [
    {
      subnet_name   = "forseti-subnetwork"
      subnet_ip     = "10.129.0.0/20"
      subnet_region = "us-central1"
    },
  ]
}

resource "google_compute_router" "forseti_service" {
  provider = "google.phoogle"

  name    = "forseti-service"
  network = "${module.forseti-service-network.network_self_link}"

  bgp {
    asn = "64514"
  }

  region  = "us-central1"
  project = "${module.forseti-service-project.project_id}"
}

resource "google_compute_router_nat" "forseti_service" {
  provider = "google.phoogle"

  name                               = "forseti-service"
  router                             = "${google_compute_router.forseti_service.name}"
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"

  subnetwork {
    name                    = "${module.forseti-service-network.subnets_self_links[0]}"
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }

  project = "${module.forseti-service-project.project_id}"
  region  = "${google_compute_router.forseti_service.region}"
}

resource "google_organization_iam_member" "forseti" {
  provider = "google.phoogle"

  count = "${length(local.forseti_org_required_roles)}"

  org_id = "${var.phoogle_org_id}"
  role   = "${element(local.forseti_org_required_roles, count.index)}"
  member = "serviceAccount:${module.forseti-service-project.service_account_email}"
}

// Grant the forseti service account the rights to create GCE instances within
// the host project network.
resource "google_project_iam_member" "forseti-host" {
  provider = "google.phoogle"

  count = "${length(local.forseti_host_project_required_roles)}"

  project = "${module.forseti-host-project.project_id}"
  role    = "${element(local.forseti_host_project_required_roles, count.index)}"
  member  = "serviceAccount:${module.forseti-service-project.service_account_email}"
}

// Grant the forseti service account rights over the Forseti service project.
resource "google_project_iam_member" "forseti" {
  provider = "google.phoogle"

  count = "${length(local.forseti_project_required_roles)}"

  project = "${module.forseti-service-project.project_id}"
  role    = "${element(local.forseti_project_required_roles, count.index)}"
  member  = "serviceAccount:${module.forseti-service-project.service_account_email}"
}

// Define a project for the Forseti real time enforcer.
//
// This project holds resources will be managed by the Forseti real time enforcer.
// At present this is limited to GCS buckets but may be expanded as the real time
// enforcer manages additional resources.
module "forseti-enforcer-project" {
  source  = "terraform-google-modules/project-factory/google"
  version = "~> 2.0"

  name              = "ci-forseti-enforcer"
  random_project_id = "true"
  org_id            = "${var.phoogle_org_id}"
  folder_id         = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  billing_account   = "${module.variables.phoogle_billing_account}"
  credentials_path  = "${var.phoogle_credentials_path}"

  activate_apis = [
    "compute.googleapis.com",
    "storage-api.googleapis.com",
    "storage-component.googleapis.com",
  ]

  providers {
    "google"      = "google.phoogle"
    "google-beta" = "google-beta.phoogle"
  }
}

// Grant the ci-forseti service account permissions to create log exports and
// buckets in the enforcer project.
resource "google_project_iam_member" "forseti-enforcer" {
  provider = "google.phoogle"

  count = "${length(local.forseti_enforcer_project_required_roles)}"

  project = "${module.forseti-enforcer-project.project_id}"
  role    = "${element(local.forseti_enforcer_project_required_roles, count.index)}"
  member  = "serviceAccount:${module.forseti-service-project.service_account_email}"
}

// Create the Forseti real time enforcer roles. We're using a vendored copy of the
// real_time_enforcer_roles module because this code isn't in a stable branch that we
// can access. When the real_time_enforcer_roles module is released in v1.4.0 we should
// switch to that:
//
// ```hcl
// module "real_time_enforcer_roles" {
//   source  = "terraform-google-modules/forseti/google"
//   version = "~> 1.4"
//   # ...
// }
//
// ```
module "real_time_enforcer_roles" {
  source = "../../modules/real_time_enforcer_roles"

  org_id = "${var.phoogle_org_id}"

  providers {
    "google"      = "google.phoogle"
    "google-beta" = "google-beta.phoogle"
  }
}

resource "google_service_account_key" "forseti" {
  provider = "google.phoogle"

  service_account_id = "${module.forseti-service-project.service_account_id}"
}

resource "random_id" "forseti_github_webhook_token" {
  byte_length = 20
}

data "template_file" "forseti_github_webhook_url" {
  template = "https://concourse.infra.cft.tips/api/v1/teams/cft/pipelines/$${pipeline}/resources/pull-request/check/webhook?webhook_token=$${webhook_token}"

  vars {
    pipeline      = "terraform-google-forseti"
    webhook_token = "${random_id.forseti_github_webhook_token.hex}"
  }
}

resource "kubernetes_secret" "forseti" {
  metadata {
    namespace = "concourse-cft"
    name      = "forseti"
  }

  data {
    github_webhook_token     = "${random_id.forseti_github_webhook_token.hex}"
    phoogle_project_id       = "${module.forseti-service-project.project_id}"
    phoogle_network_project  = "${module.forseti-host-project.project_id}"
    phoogle_network          = "${module.forseti-host-network-01.network_name}"
    phoogle_region           = "${module.forseti-host-network-01.subnets_regions[0]}"
    phoogle_subnetwork       = "${module.forseti-host-network-01.subnets_names[0]}"
    phoogle_enforcer_project = "${module.forseti-enforcer-project.project_id}"
    phoogle_sa               = "${base64decode(google_service_account_key.forseti.private_key)}"
  }
}
