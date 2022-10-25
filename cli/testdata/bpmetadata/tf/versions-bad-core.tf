terraform {
  required_version_not_good = ">= 0.13.0"

  provider_meta "google" {
    module_name = "blueprints/terraform/terraform-google-kubernetes-engine:hub/v23.1.0"
  }
}
