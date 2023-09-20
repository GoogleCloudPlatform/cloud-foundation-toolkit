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

PROTOC_VERSION=$1
PROTOC_GEN_GO_VERSION=$2
PROTOC_GEN_GO_GRPC_VERSION=$3
PROTOC_GEN_GO_INJECT_TAG=$4

mkdir -p /build/install_protoc
cd /build/install_protoc

curl -LO "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip"
unzip -q "protoc-${PROTOC_VERSION}-linux-x86_64.zip" -d $HOME/.local
chmod 755 $HOME/.local/bin/protoc
cp $HOME/.local/bin/protoc /usr/local/bin/
cp -R $HOME/.local/include/ /usr/local/include/

go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GEN_GO_VERSION}
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${PROTOC_GEN_GO_GRPC_VERSION}
go install github.com/favadi/protoc-go-inject-tag@v${PROTOC_GEN_GO_INJECT_TAG}

rm -rf /build/install_protoc
