resource "google_compute_global_address" "postgres" {

  provider = "google-beta"

  name          = "${module.variables.name_prefix}-postgres-${terraform.workspace}"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = "${data.terraform_remote_state.networking.network_self_link}"
}

resource "google_service_networking_connection" "postgres" {

  provider = "google-beta"

  network = "${data.terraform_remote_state.networking.network_self_link}"
  service = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [
    "${google_compute_global_address.postgres.name}",
  ]
}
