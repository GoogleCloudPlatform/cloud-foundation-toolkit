#! /bin/bash
# Copyright 2022 Google LLC
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

mkdir -p /build/install_terraform
cd /build/install_terraform

TERRAFORM_VERSION=$1

wget -nv "https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip"
unzip -q terraform_${TERRAFORM_VERSION}_linux_amd64.zip
install -o 0 -g 0 -m 0755 terraform /usr/local/bin/terraform

rm -rf /build/install_terraform
