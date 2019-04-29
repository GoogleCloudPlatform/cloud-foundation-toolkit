output "project_id" {
  description = "Project ID of the created seed project"
  value       = "${module.project_factory.project_id}"
}

output "service_account" {
  description = "Email of the service account to be used for test project creation"
  value       = "${module.project_factory.service_account_email}"
}
