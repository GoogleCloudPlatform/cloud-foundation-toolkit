/**
 * Copyright 2019 Google LLC
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
  user_project_override = true
  billing_project       = local.ci_project_id
  alias                 = "override"
}

resource "google_access_context_manager_access_policy" "access_policy" {
  provider = google.override
  parent   = "organizations/${local.org_id}"
  title    = "default policy"
}

// For ingress policy
resource "google_access_context_manager_service_perimeter" "default_perimeter" {
  parent = "accesspolicies/${google_access_context_manager_access_policy.access_policy.name}"
  name   = "accesspolicies/${google_access_context_manager_access_policy.access_policy.name}/serviceperimeters/default_perimeter"
  title  = "Default Perimeter"
  lifecycle {
    ignore_changes = [status[0].resources]
  }
}

// For AlloyDB PSC
resource "google_access_context_manager_service_perimeter_ingress_policy" "deafult_ingress_policy" {
  perimeter = google_access_context_manager_service_perimeter.default_perimeter.name
  ingress_from {
    identity_type = "ANY_SERVICE_ACCOUNT"
    sources {
      access_level = "*"
    }
  }
  ingress_to {
    resources = ["*"]
    operations {
      service_name = "networkconnectivity.googleapis.com"
      method_selectors {
        method = "*"
      }
    }
  }
  lifecycle {
    create_before_destroy = true
  }
}
