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

mkdir -p /build/install_shellcheck
cd /build/install_shellcheck

wget -nv "https://github.com/koalaman/shellcheck/releases/download/v0.6.0/shellcheck-v0.6.0.linux.x86_64.tar.xz"
tar -xf shellcheck-v0.6.0.linux.x86_64.tar.xz
install -o 0 -g 0 -m 0755 shellcheck-v0.6.0/shellcheck /usr/local/bin/shellcheck

rm -rf /build/install_shellcheck
