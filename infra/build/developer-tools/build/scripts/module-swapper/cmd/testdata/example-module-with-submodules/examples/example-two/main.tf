# Unrelated module
module "test-unrelated-submodule" {
  source  = "terraform-google-modules/foo/google"
  version = "~> 3.2.0"

  project_id   = var.project_id # Replace this with your project ID in quotes
  network_name = "my-custom-mode-network"
  mtu          = 1460
}

module "test-module" {
  source  = "terraform-google-modules/example-module-with-submodules/google"
  version = "~> 3.2.0"

  project_id   = var.project_id # Replace this with your project ID in quotes
  network_name = "my-custom-mode-network"
  mtu          = 1460
}

module "test-submodule-1" {
  source  = "terraform-google-modules/example-module-with-submodules/google//modules/bar-module"
  version = "~> 3.2.0"

  project_id   = var.project_id # Replace this with your project ID in quotes
  network_name = "my-custom-mode-network"
  mtu          = 1460
}

module "test-submodule-2" {
  source  = "terraform-google-modules/example-module-with-submodules/google//modules/foo-module"
  version = "~> 3.2.0"

  project_id   = var.project_id # Replace this with your project ID in quotes
  network_name = "my-custom-mode-network"
  mtu          = 1460
}

# Unrelated submodule
module "test-unrelated-submodule" {
  source  = "terraform-google-modules/foo/google//modules/bar-module"
  version = "~> 3.2.0"

  project_id   = var.project_id # Replace this with your project ID in quotes
  network_name = "my-custom-mode-network"
  mtu          = 1460
}
