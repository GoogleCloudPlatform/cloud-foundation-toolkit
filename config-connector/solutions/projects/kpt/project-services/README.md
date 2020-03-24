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
  If you are enabling services for a project other than the one you have
installed Config Connector in, set the `project-id` and `namespace`:

  ```
  kpt cfg set . project-id VALUE
  kpt cfg set . namespace project-annotated
  ```
  
  Before applying a service, set the service name. For example, to enable
[Compute Engine](https://cloud.google.com/compute/docs):

  ```
  kpt cfg set . service-name compute.googleapis.com
  ```

  Note: the package-default value will enable
[Firebase](https://firebase.google.com/docs).

  Once your information is in the configs, simply apply.

  ```
  kubectl apply -f .
  ```

  To enable multiple services, copy the `service.yaml` into either a separate
file or the same file seperated by a yaml seperator and manually change its
`metadata.name`.


# License

  Apache 2.0 - See [LICENSE](/LICENSE) for more information.

