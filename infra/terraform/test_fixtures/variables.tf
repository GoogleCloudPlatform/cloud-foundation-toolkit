module "variables" { source = "../variables" }

variable "phoogle_credentials_path" {
  description = "Path to credentials file for phoogle.net organization."
}

variable "phoogle_org_id" {
  default = "826592752744"
}

variable "core_group" {
  description = "The CFT core group"
  default     = "cloud-foundation-core@google.com"
}
