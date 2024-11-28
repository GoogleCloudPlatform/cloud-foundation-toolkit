locals {
  int_required_roles = [
    "roles/cloudsql.admin",
    "roles/compute.networkAdmin",
    "roles/iam.serviceAccountAdmin",
    "roles/resourcemanager.projectIamAdmin",
    "roles/storage.admin",
    "roles/workflows.admin",
    "roles/cloudscheduler.admin",
    "roles/iam.serviceAccountUser"
  ]
  int_required_project_roles = [
    "roles/owner",
    "roles/storage.admin"
  ]
}
