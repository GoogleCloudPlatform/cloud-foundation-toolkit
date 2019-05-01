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

output "forseti-rt-enforcer-viewer-role-id" {
  description = "The forseti real time enforcer viewer Role ID."
  value       = "${google_organization_iam_custom_role.forseti-enforcer-viewer.role_id}"
}

output "forseti-rt-enforcer-writer-role-id" {
  description = "The forseti real time enforcer writer Role ID."
  value       = "${google_organization_iam_custom_role.forseti-enforcer-writer.role_id}"
}
