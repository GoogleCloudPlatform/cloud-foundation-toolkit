Subnet
==================================================

# NAME

  subnet

# SYNOPSIS

  Config Connector compatible YAML files to create a subnet in your desired project, and grant a specific member a role (default to `roles/compute.networkUser`) for accessing the subnet that just created.

# CONSUMPTION

  Fetch the [kpt](https://googlecontainertools.github.io/kpt/) package of the solution:

  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/iam/kpt/subnet subnet
  ```

# REQUIREMENTS

  -   A working Config Connector instance using the "cnrm-system" service
      account with either both `roles/compute.networkAdmin` and
      `roles/iam.securityAdmin` roles or `roles/owner` in the project
      managed by Config Connector.
  -   Compute Engine API enabled in the project where Config Connector is
      installed
  -   Compute Engine API enabled in the project managed by Config Connector if
      it is a different project

# USAGE
  Replace `${IAM_MEMBER?}` with the GCP identity to grant access to:
  ```
  kpt cfg set . iam-member user:name@example.com
  ```

  _Optionally_, you can change the following fields before you apply the YAMLs:
  - the name of the compute network
  ```
  kpt cfg set . compute-network-name VALUE
  ```
  
  - the name of the subnet
  ```
  kpt cfg set . subnet-name new-subnet-name
  ```
  
  - the region of the subnet
  ```
  kpt cfg set . subnet-region us-west1
  ```

  - the role granted to the GCP identity.
  (you can find all of the subnet related IAM roles
  [here](https://cloud.google.com/iam/docs/understanding-roles#compute-engine-roles)):

  ```
  kpt cfg set . role roles/compute.networkViewer
  ```

  Apply the YAMLs:

  ```
  kubectl apply -f .
  ```

# LICENSE

  Apache 2.0 - See [LICENSE](/LICENSE) for more information.
