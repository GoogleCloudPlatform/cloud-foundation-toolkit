variable "domain" {
  description = "The GCP domain where G Suite resources will be managed."
  default     = "phoogle.net"
}

variable "credentials_path" {
  description = "Path to the credentials file for phoogle.net organization."
}

variable "core_group" {
  description = "The Cloud Foundation Team core group"
  default     = "cloud-foundation-core@google.com"
}

variable "impersonated_user_email" {
  description = <<EOD
A G Suite user account to impersonate when managing G Suite resources.
Prefer using your own account when possible for auditability.
EOD
}

variable "org_id" {
  description = "The phoogle.net organization ID"
  default     = "826592752744"
}
