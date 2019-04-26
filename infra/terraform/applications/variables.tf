module "variables" { source = "../variables" }

variable "concourse_subdomain" {
  default = {
    primary = "concourse"
  }
}

variable "lets_encrypt_email" {
  default = {
    primary = "ctrott@google.com"
  }
}

variable "phoogle_credentials_path" {
  description = "Path to credentials file for phoogle.net organization."
}
