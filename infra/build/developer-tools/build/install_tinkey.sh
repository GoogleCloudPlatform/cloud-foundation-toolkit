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

cd /build

TINKEY_VERSION=$1

gsutil cp "gs://tinkey/tinkey-${TINKEY_VERSION}.tar.gz" .
tar -xzvf "tinkey-${TINKEY_VERSION}.tar.gz"

install -o 0 -g 0 -m 0755 tinkey_deploy.jar /usr/bin/
install -o 0 -g 0 -m 0755 tinkey /usr/bin/

rm "tinkey-${TINKEY_VERSION}.tar.gz"
rm tinkey
rm tinkey.bat
rm -f tinkey_deploy.jar
