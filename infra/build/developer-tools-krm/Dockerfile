# Copyright 2021 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ARG BASE_IMAGE_VERSION
FROM cft/developer-tools:$BASE_IMAGE_VERSION

RUN apk update && apk add --no-cache docker-cli nodejs-current nodejs-npm

# Additional go tooling
RUN GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.40.1
