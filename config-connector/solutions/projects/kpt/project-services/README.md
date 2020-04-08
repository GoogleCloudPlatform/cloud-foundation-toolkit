Project Services
==================================================

# NAME

  project-services

# SYNOPSIS

  Config Connector compatible YAML files to enable services on a project.

## Consumption

  Using kpt:
  ```
  BASE=https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
  kpt pkg get $BASE/config-connector/solutions/projects/kpt/project-services project-services
  ```

## Requirements

  A working cluster with Config Connector installed.

  The cnrm-system service account must have
`roles/serviceusage.serviceUsageAdmin` or `roles/owner` for the desired project.


## Usage
  Set project-id to the ID of the project to enable services for:
  ```
  kpt cfg set . project-id VALUE
  ```


  Optionally, change the service name before applying the service. For example, to enable
[Compute Engine](https://cloud.google.com/compute/docs):
  ```
  kpt cfg set . service-name compute.googleapis.com
  ```

  The package-default value will enable
[Firebase](https://firebase.google.com/docs).

  Once your configuration is complete, simply apply:
  ```
  kubectl apply -f .
  ```

  Note: services that have been applied will have type `gcpservice`


# License

  Apache 2.0 - See [LICENSE](/LICENSE) for more information.

