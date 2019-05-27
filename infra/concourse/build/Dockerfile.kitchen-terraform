# Copyright 2018 Google LLC
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

ARG BUILD_TERRAFORM_IMAGE
ARG BUILD_RUBY_VERSION
FROM $BUILD_TERRAFORM_IMAGE as cft-terraform

FROM ruby:$BUILD_RUBY_VERSION-alpine

RUN apk add --no-cache \
    bash \
    coreutils \
    curl \
    git \
    g++ \
    jq \
    make \
    musl-dev \
    openssh \
    python \
    python3

SHELL ["/bin/bash", "-c"]

ENV APP_BASE_DIR="/cft"

RUN cd /tmp && \
    wget https://releases.hashicorp.com/packer/1.4.1/packer_1.4.1_linux_amd64.zip && \
    unzip packer_1.4.1_linux_amd64.zip && \
    rm -rf packer_1.4.1_linux_amd64.zip && \
    mv packer /bin/

COPY --from=cft-terraform $APP_BASE_DIR $APP_BASE_DIR

ENV HOME="$APP_BASE_DIR/home"
ENV PATH $APP_BASE_DIR/bin:$APP_BASE_DIR/google-cloud-sdk/bin:$PATH
ENV GOOGLE_APPLICATION_CREDENTIALS="$CREDENTIALS_PATH" \
    CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE="$CREDENTIALS_PATH"

# Fix base64 inconsistency
SHELL ["/bin/bash", "-c"]
RUN echo 'base64() { if [[ $@ == "--decode" ]]; then command base64 -d | more; else command base64 "$@"; fi; }' >> $APP_BASE_DIR/home/.bashrc

RUN terraform --version && \
    gcloud --version && \
    ruby --version && \
    bundle --version && \
    packer --version

WORKDIR /opt/kitchen
ADD ./build/data/Gemfile .
ADD ./build/data/Gemfile.lock .
ADD ./build/data/requirements.txt .
RUN bundle install && pip3 install -r requirements.txt

WORKDIR $APP_BASE_DIR/workdir

RUN gcloud components install beta --quiet
RUN gcloud components install alpha --quiet

# Authenticate gcloud with service account credentials key to allow gsutil authentication
ADD ./build/scripts/gcloud_auth.sh $HOME/entrypoint_scripts/
RUN chmod +x $HOME/entrypoint_scripts/gcloud_auth.sh
ENTRYPOINT ["/cft/home/entrypoint_scripts/gcloud_auth.sh"]
