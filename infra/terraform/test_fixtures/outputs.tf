output "github_webhook_urls" {
  description = "Webhook payload URLs to configure with module Github repositories"
  value = {
    terraform-google-container-vm = "${data.template_file.container_vm_github_webhook_url.rendered}"
    terraform-google-kubernetes-engine = "${data.template_file.kubernetes_engine_github_webhook_url.rendered}"
    terraform-google-project-factory = "${data.template_file.project_factory_github_webhook_url.rendered}"
  }
}
