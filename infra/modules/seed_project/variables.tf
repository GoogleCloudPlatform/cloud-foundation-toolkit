variable "username" {
  description = "Username of the user for whom the project is created"
}

variable "owner_emails" {
  description = "A list of identities to add as owners on the project in the member format described here: https://cloud.google.com/iam/docs/overview#iam_policy"
  type        = list
}

variable "seed_project_services" {
  description = "A list of APIs to enable on the project"

  default = [
    "admin.googleapis.com",
    "appengine.googleapis.com",
    "cloudbilling.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "compute.googleapis.com",
    "iam.googleapis.com",
    "iamcredentials.googleapis.com",
    "oslogin.googleapis.com",
    "serviceusage.googleapis.com",
  ]

  type = list
}

variable "seed_project_roles" {
  description = "A list of roles to grant to the project owners (including groups and service accounts)"

  default = [
    "roles/resourcemanager.projectCreator",
    "roles/compute.xpnAdmin",
    "roles/compute.networkAdmin",
    "roles/iam.serviceAccountAdmin",
    "roles/resourcemanager.projectIamAdmin",
  ]

  type = list
}

variable "org_id" {
  description = "Organization in which to place the project"
}

variable "seed_folder_id" {
  description = "The folder in which to place the created project"
}

variable "billing_account" {
  description = "Billing account to attach to the created project"
}
