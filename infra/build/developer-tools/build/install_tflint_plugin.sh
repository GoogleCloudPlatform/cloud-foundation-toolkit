#! /bin/bash
# Copyright 2024 Google LLC
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

mkdir -p /build/install_tflint_plugin
cd /build/install_tflint_plugin

TFLINT_BP_PLUGIN_VERSION=$1
wget -nv "https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/releases/download/tflint-ruleset-blueprint%2Fv${TFLINT_BP_PLUGIN_VERSION}/tflint-ruleset-blueprint_linux_amd64.zip"
unzip -q tflint-ruleset-blueprint_linux_amd64.zip
mkdir -p ~/.tflint.d/plugins
install -o 0 -g 0 -m 0755 tflint-ruleset-blueprint ~/.tflint.d/plugins
rm -rf /build/install_tflint_plugin
