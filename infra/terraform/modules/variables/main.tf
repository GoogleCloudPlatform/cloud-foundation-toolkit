variable "name_prefix" {
  default     = "cicd"
  description = "Common prefix for naming resources such as networks and k8s clusters."
}

variable "project_id" {
  default     = "cloud-foundation-cicd"
  description = "ID of project where all CICD resources will be launched."
}

variable "region" {
  default = {
    primary = "us-west1"
  }
  description = "GCP region to launch resources in. Keys should correspond to Terraform workspaces."
}

variable "phoogle_billing_account" {
  default = "01E8A0-35F760-5CF02A"
}
