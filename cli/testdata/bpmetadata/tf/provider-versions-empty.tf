terraform {
  required_version = ">= 0.13.0"

  required_providers {
  }

  provider_meta "google" {
    module_name = "blueprints/terraform/terraform-google-kubernetes-engine:hub/v23.1.0"
  }
  provider_meta "google-beta" {
    module_name = "blueprints/terraform/terraform-google-kubernetes-engine:hub/v23.1.0"
  }
}
