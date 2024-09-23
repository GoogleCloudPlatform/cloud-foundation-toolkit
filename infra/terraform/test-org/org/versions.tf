/**
 * Copyright 2019-2023 Google LLC
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

terraform {
  required_version = ">= 1.4.4"
  required_providers {
    external = {
      source  = "hashicorp/external"
      version = ">= 1.2, < 3"
    }
    google = {
      source  = "hashicorp/google"
      version = ">= 3.19, < 7"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = ">= 3.19, < 7"
    }
    null = {
      source  = "hashicorp/null"
      version = ">= 2.1, < 4"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 2.3.1, < 4"
    }
    github = {
      source  = "integrations/github"
      version = "~> 6.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.13, < 3"
    }
  }
}
