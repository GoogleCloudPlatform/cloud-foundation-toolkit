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

mkdir -p /build/install_kustomize
cd /build/install_kustomize

KUSTOMIZE_VERSION=$1

wget -nv https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v${KUSTOMIZE_VERSION}/kustomize_v${KUSTOMIZE_VERSION}_linux_amd64.tar.gz
tar -xzf kustomize_v${KUSTOMIZE_VERSION}_linux_amd64.tar.gz
install -o 0 -g 0 -m 0755 kustomize /usr/local/bin/kustomize

rm -rf /build/install_kustomize
