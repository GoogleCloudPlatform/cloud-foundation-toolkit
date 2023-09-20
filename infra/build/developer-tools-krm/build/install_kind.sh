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

mkdir -p /build/install_kind
cd /build/install_kind

KIND_VERSION=$1

wget -nv "https://kind.sigs.k8s.io/dl/v${KIND_VERSION}/kind-linux-amd64"
install -o 0 -g 0 -m 0755 kind-linux-amd64 /usr/local/bin/kind

rm -rf /build/install_kind
