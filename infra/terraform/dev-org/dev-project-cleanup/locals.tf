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

locals {
  org_id          = "684124036889"
  billing_account = "01B29E-828C10-4D001A"
  cleanup_folder  = "690527133635" // cft-infra-dev
  exclude_labels  = { "cft-ci" = "permanent" }
  region          = "us-central1"
  app_location    = "us-central"
}
