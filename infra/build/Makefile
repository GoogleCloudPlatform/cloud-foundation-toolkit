# Copyright 2019 Google LLC
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

SHELL := /usr/bin/env bash # Make will use bash instead of sh
TERRAFORM_VERSION := 0.12.17
CLOUD_SDK_VERSION := 271.0.0
GSUITE_PROVIDER_VERSION := 0.1.22
TERRAFORM_DOCS_VERSION := 0.6.0
RUBY_VERSION := 2.6.3
BATS_VERSION := 0.4.0
BATS_SUPPORT_VERSION := 0.3.0
BATS_ASSERT_VERSION := 2.0.0
BATS_MOCK_VERSION := 1.0.1
PACKER_VERSION := 1.4.4

DOCKER_TAG_VERSION_DEVELOPER_TOOLS := 0.7.0
DOCKER_TAG_MAJOR_VERSION_DEVELOPER_TOOLS = $(firstword $(subst ., , $(DOCKER_TAG_VERSION_DEVELOPER_TOOLS)))
DOCKER_IMAGE_DEVELOPER_TOOLS := cft/developer-tools
DOCKER_IMAGE_DEVELOPER_TOOLS_JENKINS := cft/developer-tools-jenkins
REGISTRY_URL := gcr.io/cloud-foundation-cicd


.PHONY: build-image-developer-tools
build-image-developer-tools:
	docker build \
		--build-arg CLOUD_SDK_VERSION=${CLOUD_SDK_VERSION} \
		--build-arg GSUITE_PROVIDER_VERSION=${GSUITE_PROVIDER_VERSION} \
		--build-arg TERRAFORM_VERSION=${TERRAFORM_VERSION} \
		--build-arg TERRAFORM_DOCS_VERSION=${TERRAFORM_DOCS_VERSION} \
		--build-arg RUBY_VERSION=${RUBY_VERSION} \
		--build-arg BATS_VERSION=${BATS_VERSION} \
		--build-arg BATS_SUPPORT_VERSION=${BATS_SUPPORT_VERSION} \
		--build-arg BATS_ASSERT_VERSION=${BATS_ASSERT_VERSION} \
		--build-arg BATS_MOCK_VERSION=${BATS_MOCK_VERSION} \
		-t ${DOCKER_IMAGE_DEVELOPER_TOOLS}:${DOCKER_TAG_VERSION_DEVELOPER_TOOLS} "./developer-tools"
	docker build \
		--build-arg BASE_IMAGE_VERSION=${DOCKER_TAG_VERSION_DEVELOPER_TOOLS} \
		--build-arg PACKER_VERSION=${PACKER_VERSION} \
		-t ${DOCKER_IMAGE_DEVELOPER_TOOLS_JENKINS}:${DOCKER_TAG_VERSION_DEVELOPER_TOOLS} "./developer-tools-jenkins"

.PHONY: release-image-developer-tools
release-image-developer-tools:
	docker tag ${DOCKER_IMAGE_DEVELOPER_TOOLS}:${DOCKER_TAG_VERSION_DEVELOPER_TOOLS} \
  ${REGISTRY_URL}/${DOCKER_IMAGE_DEVELOPER_TOOLS}:${DOCKER_TAG_VERSION_DEVELOPER_TOOLS}
	docker push ${REGISTRY_URL}/${DOCKER_IMAGE_DEVELOPER_TOOLS}:${DOCKER_TAG_VERSION_DEVELOPER_TOOLS}

	docker tag ${DOCKER_IMAGE_DEVELOPER_TOOLS}:${DOCKER_TAG_VERSION_DEVELOPER_TOOLS} \
	${REGISTRY_URL}/${DOCKER_IMAGE_DEVELOPER_TOOLS}:${DOCKER_TAG_MAJOR_VERSION_DEVELOPER_TOOLS}
	docker push ${REGISTRY_URL}/${DOCKER_IMAGE_DEVELOPER_TOOLS}:${DOCKER_TAG_MAJOR_VERSION_DEVELOPER_TOOLS}

	docker tag ${DOCKER_IMAGE_DEVELOPER_TOOLS_JENKINS}:${DOCKER_TAG_VERSION_DEVELOPER_TOOLS} \
  ${REGISTRY_URL}/${DOCKER_IMAGE_DEVELOPER_TOOLS_JENKINS}:${DOCKER_TAG_VERSION_DEVELOPER_TOOLS}
	docker push ${REGISTRY_URL}/${DOCKER_IMAGE_DEVELOPER_TOOLS_JENKINS}:${DOCKER_TAG_VERSION_DEVELOPER_TOOLS}

	docker tag ${DOCKER_IMAGE_DEVELOPER_TOOLS_JENKINS}:${DOCKER_TAG_VERSION_DEVELOPER_TOOLS} \
  ${REGISTRY_URL}/${DOCKER_IMAGE_DEVELOPER_TOOLS_JENKINS}:${DOCKER_TAG_MAJOR_VERSION_DEVELOPER_TOOLS}
	docker push ${REGISTRY_URL}/${DOCKER_IMAGE_DEVELOPER_TOOLS_JENKINS}:${DOCKER_TAG_MAJOR_VERSION_DEVELOPER_TOOLS}

# The following tasks are provided as examples and aren't invoked from this
# repository directly. These three make targets are the only targets that
# should be included in every CFT module's `Makefile`
# Enter docker container for local development
.PHONY: docker_run
docker_run:
	docker run --rm -it \
		-e SERVICE_ACCOUNT_JSON \
		-v $(CURDIR):/workspace \
		$(REGISTRY_URL)/${DOCKER_IMAGE_DEVELOPER_TOOLS}:${DOCKER_TAG_VERSION_DEVELOPER_TOOLS} \
		/bin/bash

# Execute integration tests within the docker container
.PHONY: docker_test_integration
docker_test_integration:
	docker run --rm -it \
		-e SERVICE_ACCOUNT_JSON \
		-v $(CURDIR):/workspace \
		$(REGISTRY_URL)/${DOCKER_IMAGE_DEVELOPER_TOOLS}:${DOCKER_TAG_VERSION_DEVELOPER_TOOLS} \
		/usr/local/bin/test_integration.sh

# Execute lint tests within the docker container
.PHONY: docker_test_lint
docker_test_lint:
	docker run --rm -it \
		-v $(CURDIR):/workspace \
		$(REGISTRY_URL)/${DOCKER_IMAGE_DEVELOPER_TOOLS}:${DOCKER_TAG_VERSION_DEVELOPER_TOOLS} \
		/usr/local/bin/test_lint.sh

# Generate documentation
.PHONY: generate_docs
generate_docs:
	docker run --rm -it \
		-v $(CURDIR):/workspace \
		$(REGISTRY_URL)/${DOCKER_IMAGE_DEVELOPER_TOOLS}:${DOCKER_TAG_VERSION_DEVELOPER_TOOLS} \
		/bin/bash -c 'source /usr/local/bin/task_helper_functions.sh && generate_docs'
