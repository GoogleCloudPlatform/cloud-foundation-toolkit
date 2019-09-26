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

BATS_MOCK_VERSION=$1

cd /build
wget "https://github.com/jasonkarns/bats-mock/archive/v${BATS_MOCK_VERSION}.zip"
unzip "v${BATS_MOCK_VERSION}.zip"
cp -r "bats-mock-${BATS_MOCK_VERSION}" /usr/local/bats-mock
rm -rf "v${BATS_MOCK_VERSION}.zip" "bats-mock-${BATS_MOCK_VERSION}"
