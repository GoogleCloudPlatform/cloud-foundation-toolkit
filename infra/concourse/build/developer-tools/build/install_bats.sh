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

BATS_VERSION=$1

# bats required envsubst missing in Alpine by default
apk add gettext libintl

cd /build
wget "https://github.com/sstephenson/bats/archive/v${BATS_VERSION}.zip"
unzip "v${BATS_VERSION}.zip"
cd "bats-${BATS_VERSION}"
./install.sh /usr/local
rm -rf "v${BATS_VERSION}" "bats-${BATS_VERSION}"

