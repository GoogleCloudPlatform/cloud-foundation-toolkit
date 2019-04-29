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

variable "suffix" {
  description = "An optional suffix to append to the role IDs. Role ID suffixes are only needed if the roles are being recreated within 40 days of their deletion."
  default     = ""
}

variable "org_id" {
  description = "The organization ID where the custom Forseti roles will be created."
}

variable "prevent_destroy" {
  description = "Prevent the destruction of the IAM roles. IAM roles have a soft delete that prevent names from being reused, so this should only be set to false if the roles are being permanently torn down."
  default     = "true"
}
