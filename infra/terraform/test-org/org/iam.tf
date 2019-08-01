
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

module "iam_binding" {
  source = "terraform-google-modules/iam/google"

  folders = [local.folders["ci"]]

  bindings = {
    "roles/resourcemanager.projectCreator" = [
      "group:cft-ci-robots@test.infra.cft.tips",
    ]

    "roles/resourcemanager.folderViewer" = [
      "group:cft-ci-robots@test.infra.cft.tips",
    ]
  }
}
