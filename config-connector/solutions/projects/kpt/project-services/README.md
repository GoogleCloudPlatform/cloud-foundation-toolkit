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
  The cnrm-system service account must have `roles/serviceusage.serviceUsageAdmin` or `roles/owner`.


## Usage
  First, set the project-id:

  ```
  kpt cfg set . project-id VALUE 
  ```

  Before applying a service, set the service name. For example, to enable [Compute Engine](https://cloud.google.com/compute/docs):

  ```
  kpt cfg set . service-name compute.googleapis.com
  ```

  Note: the package-default value will enable [Firebase](https://firebase.google.com/docs).

  Once your information is in the configs, simply apply.

  ```
  kubectl apply -f .
  ```

  To enable multiple services, copy the yaml for one service into a separate file and manually change the `metadata.name`.

  To see services that have been enabled, run `kubectl get gcpservices`.


# License

  Apache 2.0 - See [LICENSE](/LICENSE) for more information.

