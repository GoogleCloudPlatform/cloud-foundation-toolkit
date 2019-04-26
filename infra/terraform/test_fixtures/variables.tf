module "variables" { source = "../variables" }

variable "phoogle_credentials_path" {
  description = "Path to credentials file for phoogle.net organization."
}

variable "phoogle_org_id" {
  default = "826592752744"
}
