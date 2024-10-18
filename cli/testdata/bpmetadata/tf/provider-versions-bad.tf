terraform {
  required_version = ">= 0.13.0"

  required_providers {
    google = {
      version = ">= 4.4.0, < 7"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
    }
  }

  provider_meta "google" {
    module_name = "blueprints/terraform/terraform-google-kubernetes-engine:hub/v23.1.0"
  }
  provider_meta "google-beta" {
    module_name = "blueprints/terraform/terraform-google-kubernetes-engine:hub/v23.1.0"
  }
}
