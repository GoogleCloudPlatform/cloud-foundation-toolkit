SHELL := /usr/bin/env bash # Make will use bash instead of sh

BUILD_TERRAFORM_VERSION := 0.12.0
BUILD_CLOUD_SDK_VERSION := 239.0.0
BUILD_PROVIDER_GOOGLE_VERSION := 2.7.0
BUILD_PROVIDER_GSUITE_VERSION := 0.1.19
BUILD_RUBY_VERSION := 2.6.3
# Make sure you update DOCKER_TAG_VERSION_TERRAFORM or DOCKER_TAG_VERSION_KITCHEN_TERRAFORM independently:
# If you make changes to the Docker.terraform file, update DOCKER_TAG_VERSION_TERRAFORM
# If you make changes to the Docker.kitchen-terraform file, update DOCKER_TAG_VERSION_KITCHEN_TERRAFORM
# Also make sure to update the version appropriately as described below
# Removing software components or upgrading a component to a backwards incompatible release should constitute a major release.
# Adding a component or upgrading a component to a backwards compatible release should constitute a minor release.
# Fixing bugs or making trivial changes should be considered a patch release.

DOCKER_TAG_VERSION_TERRAFORM := 2.0.0
DOCKER_TAG_VERSION_KITCHEN_TERRAFORM := 2.0.0

REGISTRY_URL := gcr.io/cloud-foundation-cicd

DOCKER_IMAGE_LINT := cft/lint
DOCKER_TAG_LINT := 2.3.0

DOCKER_IMAGE_UNIT := cft/unit
DOCKER_TAG_UNIT := latest

DOCKER_IMAGE_TERRAFORM := cft/terraform
DOCKER_IMAGE_KITCHEN_TERRAFORM := cft/kitchen-terraform

.PHONY: build-image-lint
build-image-lint:
	docker build -f build/Dockerfile.lint \
		-t ${DOCKER_IMAGE_LINT}:${DOCKER_TAG_LINT} .

.PHONY: build-image-unit
build-image-unit:
	docker build -f build/Dockerfile.unit \
		-t ${DOCKER_IMAGE_UNIT}:${DOCKER_TAG_UNIT} .

.PHONY: build-image-terraform
build-image-terraform:
	docker build -f build/Dockerfile.terraform \
		--build-arg BUILD_TERRAFORM_VERSION=${BUILD_TERRAFORM_VERSION} \
		--build-arg BUILD_CLOUD_SDK_VERSION=${BUILD_CLOUD_SDK_VERSION} \
		--build-arg BUILD_PROVIDER_GOOGLE_VERSION=${BUILD_PROVIDER_GOOGLE_VERSION} \
		--build-arg BUILD_PROVIDER_GSUITE_VERSION=${BUILD_PROVIDER_GSUITE_VERSION} \
		-t ${DOCKER_IMAGE_TERRAFORM}:${DOCKER_TAG_VERSION_TERRAFORM} .

.PHONY: build-image-kitchen-terraform
build-image-kitchen-terraform:
	docker build -f build/Dockerfile.kitchen-terraform \
		--build-arg BUILD_TERRAFORM_IMAGE=${REGISTRY_URL}/${DOCKER_IMAGE_TERRAFORM}:${DOCKER_TAG_VERSION_TERRAFORM} \
		--build-arg BUILD_RUBY_VERSION=${BUILD_RUBY_VERSION} \
		-t ${DOCKER_IMAGE_KITCHEN_TERRAFORM}:${DOCKER_TAG_VERSION_KITCHEN_TERRAFORM} .

.PHONY: release-image-lint
release-image-lint:
	docker tag ${DOCKER_IMAGE_LINT}:${DOCKER_TAG_LINT} \
		${REGISTRY_URL}/${DOCKER_IMAGE_LINT}:${DOCKER_TAG_LINT}
	docker push ${REGISTRY_URL}/${DOCKER_IMAGE_LINT}:${DOCKER_TAG_LINT}

.PHONY: release-image-unit
release-image-unit:
	docker tag ${DOCKER_IMAGE_UNIT}:${DOCKER_TAG_UNIT} \
		${REGISTRY_URL}/${DOCKER_IMAGE_UNIT}:${DOCKER_TAG_UNIT}
	docker push ${REGISTRY_URL}/${DOCKER_IMAGE_UNIT}:${DOCKER_TAG_UNIT}

.PHONY: release-image-terraform
release-image-terraform:
	docker tag ${DOCKER_IMAGE_TERRAFORM}:${DOCKER_TAG_VERSION_TERRAFORM} \
		${REGISTRY_URL}/${DOCKER_IMAGE_TERRAFORM}:${DOCKER_TAG_VERSION_TERRAFORM}
	docker push ${REGISTRY_URL}/${DOCKER_IMAGE_TERRAFORM}:${DOCKER_TAG_VERSION_TERRAFORM}

.PHONY: release-image-kitchen-terraform
release-image-kitchen-terraform:
	docker tag ${DOCKER_IMAGE_KITCHEN_TERRAFORM}:${DOCKER_TAG_VERSION_KITCHEN_TERRAFORM} \
		${REGISTRY_URL}/${DOCKER_IMAGE_KITCHEN_TERRAFORM}:${DOCKER_TAG_VERSION_KITCHEN_TERRAFORM}
	docker push ${REGISTRY_URL}/${DOCKER_IMAGE_KITCHEN_TERRAFORM}:${DOCKER_TAG_VERSION_KITCHEN_TERRAFORM}
