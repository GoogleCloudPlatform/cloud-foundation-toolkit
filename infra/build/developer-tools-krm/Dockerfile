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

RUN apk update && apk add --no-cache openrc docker-cli docker screen nodejs-current npm

# Additional go tooling
RUN GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.40.1

# Required to download and install asmcli
ARG ASMCLI_VERSION
ENV ASMCLI_VERSION ${ASMCLI_VERSION}

ADD ./build/install_asmcli.sh /build/
RUN /build/install_asmcli.sh ${ASMCLI_VERSION}

# Add dind helper for prow
ADD ./build/scripts/runner.sh /usr/local/bin/
ADD ./build/scripts/prow_entrypoint.sh /usr/local/bin/
