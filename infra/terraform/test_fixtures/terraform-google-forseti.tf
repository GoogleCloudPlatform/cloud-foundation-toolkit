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
  version = "~> 0.6.0"

  project_id      = "${module.forseti-host-project.project_id}"
  network_name    = "network-01"
  shared_vpc_host = "true"

  subnets = [
    {
      subnet_name   = "us-central1-01"
      subnet_ip     = "10.128.0.0/20"
      subnet_region = "us-central1"
    },
  ]

  secondary_ranges {
    us-central1-01 = []
  }

  providers {
    "google"      = "google.phoogle"
    "google-beta" = "google-beta.phoogle"
  }
}

// Create a Forseti service project.
//
// Note: This project is both attached to the Forseti host project and has a
//       network defined within this project. We use both the host project shared
//       VPC and a network defined within this project in different test cases.
//
// Note: Because this isn't using the Project Factory we're using the default
//       network and firewall rules. If this project is rebuilt we should switch
//       to the Project Factory.
resource "google_project" "forseti" {
  provider = "google.phoogle"

  name            = "ci-forseti"
  project_id      = "ci-forseti"
  folder_id       = "${google_folder.phoogle_cloud_foundation_cicd.name}"
  billing_account = "${module.variables.phoogle_billing_account}"
}

// Associate the forseti host project and forseti service project.
resource "google_compute_shared_vpc_service_project" "shared_vpc_attachment" {
  provider = "google.phoogle"

  host_project    = "${module.forseti-host-project.project_id}"
  service_project = "${google_project.forseti.project_id}"
}

resource "google_project_service" "forseti" {
  provider = "google.phoogle"

  count   = "${length(local.forseti_required_apis)}"
  project = "${google_project.forseti.id}"
  service = "${element(local.forseti_required_apis, count.index)}"
}

resource "google_service_account" "forseti" {
  provider = "google.phoogle"

  project      = "${google_project.forseti.id}"
  account_id   = "ci-forseti"
  display_name = "ci-forseti"
}

resource "google_organization_iam_member" "forseti" {
  provider = "google.phoogle"

  count = "${length(local.forseti_org_required_roles)}"

  org_id = "${var.phoogle_org_id}"
  role   = "${element(local.forseti_org_required_roles, count.index)}"
  member = "serviceAccount:${google_service_account.forseti.email}"
}

// Grant the forseti service account the rights to create GCE instances within
// the host project network.
resource "google_project_iam_member" "forseti-host" {
  provider = "google.phoogle"

  count = "${length(local.forseti_host_project_required_roles)}"

  project = "${module.forseti-host-project.project_id}"
  role    = "${element(local.forseti_host_project_required_roles, count.index)}"
  member  = "serviceAccount:${google_service_account.forseti.email}"
}

// Grant the forseti service account rights over the Forseti service project.
resource "google_project_iam_member" "forseti" {
  provider = "google.phoogle"

  count = "${length(local.forseti_project_required_roles)}"

  project = "${google_project.forseti.id}"
  role    = "${element(local.forseti_project_required_roles, count.index)}"
  member  = "serviceAccount:${google_service_account.forseti.email}"
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
  member  = "serviceAccount:${google_service_account.forseti.email}"
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

  service_account_id = "${google_service_account.forseti.id}"
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
    phoogle_project_id       = "${google_project.forseti.id}"
    phoogle_network_project  = "${module.forseti-host-project.project_id}"
    phoogle_network          = "${module.forseti-host-network-01.network_name}"
    phoogle_region           = "${module.forseti-host-network-01.subnets_regions[0]}"
    phoogle_subnetwork       = "${module.forseti-host-network-01.subnets_names[0]}"
    phoogle_enforcer_project = "${module.forseti-enforcer-project.project_id}"
    phoogle_sa               = "${base64decode(google_service_account_key.forseti.private_key)}"
  }
}
