Shared VPC Network
==================================================

# NAME

  shared-vpc

# SYNOPSIS

  Config Connector YAML files to create a VPC network inside a host project,
  which can be consumed from within a service project.

## Consumption

  Using [kpt](https://googlecontainertools.github.io/kpt/):
  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/config-connector/solutions/projects/kpt/shared-vpc .
  ```

## Requirements

  A working cluster with Config Connector installed.

  The cnrm-system service account, which must have
`roles/resourcemanager.projectCreator` in the target organization if service and
host projects do not yet exist, or the `owner` role in the projects if they already exist.

## Usage
  Set the ID for billing account, host project, and service project:
  ```
  kpt cfg set . billing-account VALUE
  kpt cfg set . host-project VALUE
  kpt cfg set . service-project VALUE
  ```

  You can also change the name of the VPC network, from the default value of `sharedvpcnetwork`.

  Once your configuration is complete, simply apply:
  ```
  kubectl apply -f .
  ```

  Note: services that have been applied will have type `gcpservice`


# License

  Apache 2.0 - See [LICENSE](/LICENSE) for more information.

