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

mkdir -p /build/install_asmcli
cd /build/install_asmcli

ASMCLI_VERSION=$1

curl https://storage.googleapis.com/csm-artifacts/asm/asmcli_${ASMCLI_VERSION} > asmcli
install -o 0 -g 0 -m 0755 asmcli /usr/local/bin/asmcli

rm -rf /build/install_asmcli
