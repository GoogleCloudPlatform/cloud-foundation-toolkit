#! /bin/bash
# Copyright 2023 Google LLC
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

mkdir -p /build/install_cft_cli
cd /build/install_cft_cli

CFT_CLI_VERSION=$1

if ! wget -nv "https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/releases/download/cli%2Fv${CFT_CLI_VERSION}/cft-linux-amd64"; then
  wget -nv "https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/releases/download/v${CFT_CLI_VERSION}/cft-linux-amd64"
fi

install -o 0 -g 0 -m 0755 cft-linux-amd64 /usr/local/bin/cft

rm -rf /build/install_cft_cli
