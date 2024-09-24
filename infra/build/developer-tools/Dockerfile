# Copyright 2019-2023 Google LLC
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
# Download and verify the integrity of the download first
ARG ALPINE_VERSION
FROM alpine:$ALPINE_VERSION

# Required to override base path
ARG KITCHEN_TEST_BASE_PATH="test/integration"
ENV KITCHEN_TEST_BASE_PATH ${KITCHEN_TEST_BASE_PATH}

# $WORKSPACE is intended for assets that persist across multiple build steps in a pipeline.
# It's also where the project git repository is located.
# https://cloud.google.com/cloud-build/docs/build-config
ENV WORKSPACE="/workspace"
ENV PATH /usr/local/google-cloud-sdk/bin:$PATH

# Install Distro Dependencies
# diffutils contains the full version of diff needed for the -exclude argument.
# That argument is needed for check_documentation in task_helper_functions.sh
# rsync is needed for check_documentation in task_helper_functions.sh
# install gcompat, because protoc needs a real glibc or compatible layer
RUN apk upgrade --no-cache && \
    apk add --no-cache \
    bash \
    parallel \
    coreutils \
    curl \
    findutils \
    git \
    grep \
    g++ \
    jq \
    make \
    python3-dev \
    musl-dev \
    openssh \
    openssl \
    python3 \
    py-pip \
    ca-certificates \
    diffutils \
    rsync \
    gcompat \
    gettext \
    libintl \
    py3-crcmod \
    py3-flake8 \
    py3-jinja2 \
    py3-google-auth \
    py3-google-api-python-client \
    py3-pygithub \
    py3-requests \
    ruby-dev \
    go \
    kubectl

# Add all build scripts
ADD ./build/*.sh /build/

# Add cloud sdk first, gsutil required by some installers
ARG CLOUD_SDK_VERSION
RUN /build/install_cloud_sdk.sh ${CLOUD_SDK_VERSION}

# Optimized to minimize layers
ARG GSUITE_PROVIDER_VERSION
ARG TERRAFORM_DOCS_VERSION
ARG BATS_VERSION
ARG BATS_SUPPORT_VERSION
ARG BATS_ASSERT_VERSION
ARG BATS_MOCK_VERSION
ARG TERRAGRUNT_VERSION
ARG KUSTOMIZE_VERSION
ARG PROTOC_VERSION
ARG PROTOC_GEN_GO_VERSION
ARG PROTOC_GEN_GO_GRPC_VERSION
ARG PROTOC_GEN_GO_INJECT_TAG
ARG TFLINT_VERSION
ARG TINKEY_VERSION
RUN /build/install_shellcheck.sh && \
    /build/install_hadolint.sh && \
    /build/install_addlicense.sh && \
    /build/install_gsuite_terraform_provider.sh ${GSUITE_PROVIDER_VERSION} && \
    /build/install_terraform_docs.sh ${TERRAFORM_DOCS_VERSION} && \
    /build/install_bats.sh ${BATS_VERSION} ${BATS_SUPPORT_VERSION} ${BATS_ASSERT_VERSION} ${BATS_MOCK_VERSION} && \
    /build/install_terragrunt.sh ${TERRAGRUNT_VERSION} && \
    /build/install_kustomize.sh ${KUSTOMIZE_VERSION} && \
    /build/install_protoc.sh ${PROTOC_VERSION} ${PROTOC_GEN_GO_VERSION} ${PROTOC_GEN_GO_GRPC_VERSION} ${PROTOC_GEN_GO_INJECT_TAG} && \
    /build/install_tflint.sh ${TFLINT_VERSION} && \
    /build/install_tinkey.sh ${TINKEY_VERSION}

# Optimized to independently cache layers for commonly updated tools:
# https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/.github/workflows/update-tooling.yml
ARG TERRAFORM_VERSION
RUN /build/install_terraform.sh ${TERRAFORM_VERSION}

ARG KPT_VERSION
RUN /build/install_kpt.sh ${KPT_VERSION}

ARG CFT_CLI_VERSION
RUN /build/install_cft_cli.sh ${CFT_CLI_VERSION}

ARG MODULE_SWAPPER_VERSION
RUN /build/install_module-swapper.sh ${MODULE_SWAPPER_VERSION}

ARG TFLINT_BP_PLUGIN_VERSION
RUN /build/install_tflint_plugin.sh ${TFLINT_BP_PLUGIN_VERSION}

WORKDIR /opt/kitchen
ADD ./build/data/Gemfile .
ADD ./build/data/Gemfile.lock .
RUN gem install bundler:2.5.10 && bundle install --retry=10

COPY ./build/verify_boilerplate/ /build/verify_boilerplate/
RUN /build/install_verify_boilerplate.sh

# NOTE: Most developers are expected to source this file when entering the
# container, then run init_credentials
ADD ./build/scripts/task_helper_functions.sh /usr/local/bin/
# Provide task helper functions with wrapper scripts for Terraform module lint
# checks
ADD ./build/scripts/task_wrapper_scripts/* /usr/local/bin/
# Execute a command with SA credentials configured.
ADD ./build/scripts/execute_with_credentials.sh /usr/local/bin/
# Execute the module's integration tests and cleanup.
ADD ./build/scripts/test_integration.sh /usr/local/bin/
# Execute the module's lint tests and cleanup.
ADD ./build/scripts/test_lint.sh /usr/local/bin/
# Helper script for dynamically generating env variable export statements out of terraform outputs.
# It is used to pipe outputs from test/setup into the kitchen_do command.
ADD ./build/scripts/export_tf_outputs.py /usr/local/bin/
# Helper script for posting linting feedback to GitHub PR
# Requires GitHub PAT with repo scope, set as an env var GITHUB_PAT_TOKEN
ADD ./build/scripts/gh_lint_comment.py /usr/local/bin/

# Intended to help developers iterate quickly
ADD ./build/home/bash_history /root/.bash_history
ADD ./build/home/bashrc /root/.bashrc

# TFLint configs
ADD ./build/home/.tflint.example.hcl /root/tflint/.tflint.example.hcl
ADD ./build/home/.tflint.module.hcl /root/tflint/.tflint.module.hcl

# Set TF cache dir
ENV TF_PLUGIN_CACHE_DIR /workspace/test/integration/tmp/.terraform
RUN mkdir -p ${TF_PLUGIN_CACHE_DIR}

# Add GO_PATH TO PATH
ENV GO_PATH="/root/go"
ENV PATH=$GO_PATH/bin:$PATH

WORKDIR $WORKSPACE
RUN terraform --version && \
    terraform-docs --version && \
    gcloud --version && \
    ruby --version && \
    bundle --version && \
    kubectl version --client=true && \
    terragrunt -version && \
    kustomize version && \
    addlicense -help && \
    cft version && \
    protoc --version && \
    go version

# Cleanup intermediate build artifacts
RUN rm -rf /build
CMD ["/bin/bash"]
