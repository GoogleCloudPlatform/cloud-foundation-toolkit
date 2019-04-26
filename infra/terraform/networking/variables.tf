module "variables" { source = "../variables" }

locals {
  network_name = "${module.variables.name_prefix}-${terraform.workspace}"
  subnet_name = "${module.variables.name_prefix}"
  subnet_range_pods_name = "${module.variables.name_prefix}-pods"
  subnet_range_services_name = "${module.variables.name_prefix}-services"
}

variable "subnet_range" {
  default = {
    primary = "10.100.0.0/16"
  }
}

variable "subnet_range_pods" {
  default = {
    primary = "10.101.0.0/16"
  }
}

variable "subnet_range_services" {
  default = {
    primary = "10.102.0.0/16"
  }
}
