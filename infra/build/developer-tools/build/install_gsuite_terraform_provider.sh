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

mkdir -p /build/install_gsuite_terraform_provider
cd /build/install_gsuite_terraform_provider

GSUITE_PROVIDER_VERSION=$1

wget -nv "https://github.com/DeviaVir/terraform-provider-gsuite/releases/download/v${GSUITE_PROVIDER_VERSION}/terraform-provider-gsuite_${GSUITE_PROVIDER_VERSION}_linux_amd64.tgz"
tar -xzf "terraform-provider-gsuite_${GSUITE_PROVIDER_VERSION}_linux_amd64.tgz"
install -o 0 -g 0 -m 0755 -d ~/.terraform.d/plugins/
install -o 0 -g 0 -m 0755 "terraform-provider-gsuite_v${GSUITE_PROVIDER_VERSION}" ~/.terraform.d/plugins/

rm -rf /build/install_gsuite_terraform_provider
