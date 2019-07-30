output "github_webhook_urls" {
  description = "Webhook payload URLs to configure with module Github repositories"

  value = {
    terraform-google-address           = "${data.template_file.address_github_webhook_url.rendered}"
    terraform-google-cloud-nat         = "${data.template_file.cloud_nat_github_webhook_url.rendered}"
    terraform-google-container-vm      = "${data.template_file.container_vm_github_webhook_url.rendered}"
    terraform-google-event-function    = "${data.template_file.event_function_github_webhook_url.rendered}"
    terraform-google-forseti           = "${data.template_file.forseti_github_webhook_url.rendered}"
    terraform-google-iam               = "${data.template_file.iam_github_webhook_url.rendered}"
    terraform-google-kubernetes-engine = "${data.template_file.kubernetes_engine_github_webhook_url.rendered}"
    terraform-google-log-export        = "${data.template_file.log_export_github_webhook_url.rendered}"
    terraform-google-network           = "${data.template_file.network_github_webhook_url.rendered}"
    terraform-google-project-factory   = "${data.template_file.project_factory_github_webhook_url.rendered}"
    terraform-google-startup-scripts   = "${data.template_file.startup_scripts_github_webhook_url.rendered}"
    terraform-google-vm                = "${data.template_file.vm_github_webhook_url.rendered}"
    terraform-google-vpn               = "${data.template_file.vpn_github_webhook_url.rendered}"
  }
}

output "project_factory_cft_email" {
  description = "The shared CFT project factory service account email"
  value       = "${google_service_account.project_factory_cft.email}"
}
