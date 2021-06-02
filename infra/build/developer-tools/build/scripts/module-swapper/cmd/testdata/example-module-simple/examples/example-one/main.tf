/**
 * Copyright 2018 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

provider "google" {
  version = "~> 3.0"
}

# [START cloudrouter_create]
module "cloud_router" {
  source  = "terraform-google-modules/example-module-simple/google"
  version = "~> 0.4"

  name   = "my-router"
  region = "us-central1"

  bgp = {
    # The ASN (16550, 64512 - 65534, 4200000000 - 4294967294) can be any private ASN
    # not already used as a peer ASN in the same region and network or 16550 for Partner Interconnect.
    asn = "65001"
  }

  # project = "my-project-id"
  project = var.project
  # network = "my-network"
  network = var.network
}
# [END cloudrouter_create]

# Unrelated module
module "vpc" {
    source  = "terraform-google-modules/network/google"
    version = "~> 2.0.0"
    network_name = "example-vpc"

    routes = [
        {
            name                   = "egress-internet"
            description            = "route through IGW to access internet"
            destination_range      = "0.0.0.0/0"
            tags                   = "egress-inet"
            next_hop_internet      = "true"
        },
        {
            name                   = "app-proxy"
            description            = "route through proxy to reach app"
            destination_range      = "10.50.10.0/24"
            tags                   = "app-proxy"
            next_hop_instance      = "app-proxy-instance"
            next_hop_instance_zone = "us-west1-a"
        },
    ]
}
