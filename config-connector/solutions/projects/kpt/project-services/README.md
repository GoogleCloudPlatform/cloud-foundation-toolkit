Project Services
==================================================
# NAME
  project-services
# SYNOPSIS
  Config Connector compatible YAML files to enable services on a project.
# CONSUMPTION
  Download the package using [kpt](https://googlecontainertools.github.io/kpt/).
  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/projects/kpt/project-services project-services
  ```
# REQUIREMENTS
  A working cluster with Config Connector installed.

  The "cnrm-system" service account must have
`roles/serviceusage.serviceUsageAdmin` or `roles/owner` for the desired project.
# SETTERS
|     NAME     |             VALUE              |     SET BY      |   DESCRIPTION    | COUNT |
|--------------|--------------------------------|-----------------|------------------|-------|
| project-id   | ${PROJECT_ID?}                 | PLACEHOLDER     | ID of project    | 1     |
| service-name | firebasehosting.googleapis.com | package-default | API service name | 1     |
# USAGE

  Set project-id to the ID of the project to enable services for.
  ```
  kpt cfg set . project-id your-project-id
  ```
  _Optionally_, change the service name before applying the service. For example, to enable
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
  Note: services that have been applied will have type `gcpservice` and be in the `project-annotated` namespace.
# LICENSE
  Apache 2.0 - See [LICENSE](/LICENSE) for more information.
