terraform {
  backend "gcs" {
    bucket = "cloud-foundation-cicd-tfstate"
    prefix = "networking"
  }
}

module "network" {

  source  = "terraform-google-modules/network/google"
  version = "0.4.0"

  project_id   = "${module.variables.project_id}"

  network_name = "${local.network_name}"

  subnets = [
    {
      subnet_name   = "${local.subnet_name}"
      subnet_ip     = "${var.subnet_range[terraform.workspace]}"
      subnet_region = "${module.variables.region[terraform.workspace]}"
    }
  ]

  secondary_ranges = {
    "${local.subnet_name}" = [
      {
        range_name    = "${local.subnet_range_pods_name}"
        ip_cidr_range = "${var.subnet_range_pods[terraform.workspace]}"
      },
      {
        range_name    = "${local.subnet_range_services_name}"
        ip_cidr_range = "${var.subnet_range_services[terraform.workspace]}"
      },
    ]
  }
}
