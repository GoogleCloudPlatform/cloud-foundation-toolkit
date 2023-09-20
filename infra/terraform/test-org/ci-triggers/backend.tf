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

terraform {
  backend "gcs" {
    bucket = "cft-infra-test-tfstate"
    prefix = "state/ci-triggers"
  }
}

data "terraform_remote_state" "org" {
  backend = "gcs"
  config = {
    bucket = "cft-infra-test-tfstate"
    prefix = "state/org"
  }
}

data "terraform_remote_state" "tf-validator" {
  backend = "gcs"
  config = {
    bucket = "cft-infra-test-tfstate"
    prefix = "state/tf-validator"
  }
}

data "terraform_remote_state" "sfb-bootstrap" {
  backend = "gcs"
  config = {
    bucket = "bkt-b-tfstate-1d93"
    prefix = "terraform/bootstrap/state"
  }
}
