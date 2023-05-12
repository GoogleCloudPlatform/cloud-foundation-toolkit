# Developer Tools Image

See design decisions and discussion in [Common developer-tools container with
credentials handling (DevEx 2.0)][design]

## Building and Releasing

To build a local Docker image from the `Dockerfile`, run `build-image-developer-tools`.

To release a local Docker image to the registry, update the value of `DOCKER_TAG_VERSION_DEVELOPER_TOOLS`
in the `Makefile` following Semantic Versioning and run `release-image-developer-tools`.

Review the `Makefile` to identify other variable inputs to the build and release workflow.

## Environment Variables

The following environment variables are inputs to the running container.
Environment variables also implement feature flags to enable or disable
behavior.  These variables are considered a public API and interface into the
running container instance.

| Environment Variable | Description |
| --- | --- |
| SERVICE_ACCOUNT_JSON | The JSON string content used to initialize credentials inside the container |
| CFT_DISABLE_INIT_CREDENTIALS | Disables automatic initialization of credentials via ~root/.bashrc if the value has length '

[design]: https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/issues/255
