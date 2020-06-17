Shared VPC Network
==================================================

# NAME

  shared-vpc

# SYNOPSIS

  Config Connector YAML files to create a VPC network inside a
  host project to be consumed from within a service project.

## Consumption

  Download the package using [kpt](https://googlecontainertools.github.io/kpt/):
  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/projects/kpt/shared-vpc .
  ```

## Requirements

  A working cluster with Config Connector installed.

  The "cnrm-system" service account, which must have:
  - `roles/resourcemanager.projectCreator` in the target organization if service and
host projects do not yet exist, or the `owner` role in the projects if they already exist.
  - `roles/compute.xpnAdmin` in the target organization
  - `roles/billing.user` in the target billing account
  - Cloud Billing and Cloud Resource Manager APIs enabled in the project managed by Config Connector

## Usage
  Set the ID for billing account, host project, service project, and organization:
  ```
  kpt cfg set . billing-account VALUE
  kpt cfg set . host-project VALUE
  kpt cfg set . service-project VALUE
  kpt cfg set . org-id VALUE
  ```
  Set the default namespace setter to reflect the namespace you will apply the solution YAMLs to. This may be the namespace you set [here](https://cloud.google.com/config-connector/docs/how-to/setting-default-namespace).
  ```
  kpt cfg set . default-namespace VALUE
  ```
  where `VALUE` is the name of the namespace you found to be applicable above.

  _Optionally,_ you can also change the name of the VPC network, from the default value of `sharedvpcnetwork`.

  Once your configuration is complete, simply apply:
  ```
  kubectl apply -f .
  ```

  You can check the applied resources by running the following command:
  ```
  kubectl get -f .
  ```

  **Note:** To see the applied resources for a given namespace, run
  `kubectl get gcp --namespace <namespace>`, where `<namespace>` is replaced by
  the corresponding namespace in the `0-namespace.yaml` file. You'll need to use
  type `gcpservice` to check the status of Service resources defined in
  `service.yaml`.

  If you want to clean up the resources, run;
  ```
  kubectl delete -f .
  ```

  **Note:** If `computesharedvpchostproject` can't be deleted with
  the error message `Cannot disable project as a shared VPC host because it has
  active service projects.` but `computesharedvpcserviceproject` is already
  deleted, you'll need to [manually detach](
  https://cloud.google.com/vpc/docs/deprovisioning-shared-vpc#detach_service_projects)
  the service project (specificed using kpt setter `service-project`) from the
  host project (specified using kpt setter `host-project`). More details about
  the root cause can be found in [this GitHub issue](
  https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/167).


# License

  Apache 2.0 - See [LICENSE](/LICENSE) for more information.

