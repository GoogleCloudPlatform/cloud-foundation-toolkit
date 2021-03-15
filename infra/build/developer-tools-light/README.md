# Developer Tools Image Light

This image is light version of developer-tools with Terraform and Cloud SDK.

## Building and Releasing

To build a local Docker image from the `Dockerfile`, run `build-image-developer-tools`.

To release a local Docker image to the registry, update the value of `DOCKER_TAG_VERSION_DEVELOPER_TOOLS`
in the `Makefile` following Semantic Versioning and run `release-image-developer-tools`.

Review the `Makefile` to identify other variable inputs to the build and release workflow.

