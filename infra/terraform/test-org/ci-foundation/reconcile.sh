#! /bin/bash
# Copyright 2021 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

echo "Reconciling gcp-org in production"
BUILD_ID=$(gcloud beta builds triggers run gcp-org---terraform-apply --branch=production --project "${FOUNDATION_CICD_PROJECT_ID}" --format="value(metadata.build.id)")
gcloud beta builds log "$BUILD_ID" --project "${FOUNDATION_CICD_PROJECT_ID}" --stream

repos=("gcp-environments" "gcp-networks" "gcp-projects")
envs=("development" "non-production" "production")

for repo in "${repos[@]}"; do
    TRIGGER_NAME="${repo}---terraform-apply"
    for env in "${envs[@]}"; do
        echo "Reconciling ${repo} in ${env}"
        BUILD_ID=$(gcloud beta builds triggers run "${TRIGGER_NAME}" --branch="${env}" --project "${FOUNDATION_CICD_PROJECT_ID}" --format="value(metadata.build.id)")
        gcloud beta builds log "$BUILD_ID" --project "${FOUNDATION_CICD_PROJECT_ID}" --stream
    done
done
