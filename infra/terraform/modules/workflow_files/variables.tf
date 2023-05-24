/**
 * Copyright 2023 Google LLC
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
    name              = string
    lint_env          = optional(map(string))
    disable_lint_yaml = optional(bool)
    enable_periodic   = optional(bool)
  }))
}

variable "repo_list" {
  description = "List of Repos"
  type        = map(any)
}
