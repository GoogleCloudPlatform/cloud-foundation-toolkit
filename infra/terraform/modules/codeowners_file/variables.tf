/**
 * Copyright 2022 Google LLC
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

/******************************************
  Required variables
*******************************************/

variable "org" {
  description = "GitHub Org"
  type        = string
}

variable "owner" {
  description = "Primary owner"
  type        = string
  nullable    = false
}

variable "repos_map" {
  description = "Map of Repos"
  type = map(object({
    name         = string
    org          = string
    owners       = optional(list(string), [])
    groups       = optional(list(string), [])
  }))
}
