#! /bin/bash
# Copyright 2019 Google LLC
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
set -u

CLOUD_SDK_VERSION=$1

cd /build

wget "https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz"
tar -C /usr/local -xzf "google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz"
rm "google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz"

# TODO: Cargo-culted the symlink from a previous method. Would be nice to know
# why this is necessary
ln -s /lib /lib64

gcloud config set core/disable_usage_reporting true
gcloud config set component_manager/disable_update_check true
gcloud config set survey/disable_prompts true
gcloud components install beta --quiet
gcloud components install alpha --quiet

gcloud --version
gsutil version -l
