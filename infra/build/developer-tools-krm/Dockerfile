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

RUN apk update && apk add --no-cache openrc docker-cli docker screen nodejs-current npm yamllint util-linux-misc libc6-compat ansible-lint py3-jsonschema pipx

# Required to download and install golangci-lint
ARG GOLANGCI_VERSION
ENV GOLANGCI_VERSION ${GOLANGCI_VERSION}

ADD ./build/install_golangci_lint.sh /build/
RUN /build/install_golangci_lint.sh ${GOLANGCI_VERSION}

# Required to download and install asmcli
ARG ASMCLI_VERSION
ENV ASMCLI_VERSION ${ASMCLI_VERSION}

ADD ./build/install_asmcli.sh /build/
RUN /build/install_asmcli.sh ${ASMCLI_VERSION}

# Required to download and install gator
ARG GATOR_VERSION
ENV GATOR_VERSION ${GATOR_VERSION}

ADD ./build/install_gator.sh /build/
RUN /build/install_gator.sh ${GATOR_VERSION}

# Required to download and install OPA
ARG OPA_VERSION
ENV OPA_VERSION ${OPA_VERSION}

ADD ./build/install_opa.sh /build/
RUN /build/install_opa.sh ${OPA_VERSION}

# Required to download and install kind
ARG KIND_VERSION
ENV KIND_VERSION ${KIND_VERSION}

ADD ./build/install_kind.sh /build/
RUN /build/install_kind.sh ${KIND_VERSION}

# Required to download and install crane
ARG GCRANE_VERSION
ENV GCRANE_VERSION ${GCRANE_VERSION}

ADD ./build/install_gcrane.sh /build/
RUN /build/install_gcrane.sh ${GCRANE_VERSION}

RUN rm -rf /build

# Add dind helper for prow
ADD ./build/scripts/runner.sh /usr/local/bin/
ADD ./build/scripts/prow_entrypoint.sh /usr/local/bin/

# cookiecutter is used for terraform-google-module-template tests
# https://github.com/terraform-google-modules/terraform-google-module-template
RUN pipx install cookiecutter
