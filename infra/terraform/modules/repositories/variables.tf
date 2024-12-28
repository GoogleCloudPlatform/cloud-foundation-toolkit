/**
 * Copyright 2022-2024 Google LLC
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

variable "repos_map" {
  description = "Map of Repos"
  type = map(object({
    name         = string
    short_name   = optional(string)
    org          = string
    description  = optional(string)
    admins       = optional(list(string), [])
    maintainers  = optional(list(string), [])
    homepage_url = optional(string, null)
    module       = optional(bool, true)
    topics       = optional(string)
    groups       = optional(list(string), [])
  }))
}

variable "ci_teams" {
  description = "Repo collaborator teams for approved CI"
  type        = list(string)
  default     = []
}
