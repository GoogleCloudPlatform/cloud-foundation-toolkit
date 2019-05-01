module "variables" {
  source = "../variables"
}

variable "org_id" {
  description = "Numeric ID of the organzation to create the seed project in"
  default     = "826592752744"
}

variable "core_group" {
  description = "The CFT core group"
  default     = "cloud-foundation-core@google.com"
}
