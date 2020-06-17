Service Account
==================================================

# NAME

  service-account

# SYNOPSIS

  Config Connector compatible YAML files to create a service account in your desired project, and grant a specific member a role (default to `roles/iam.serviceAccountKeyAdmin`) for accessing the service account that just created.

# CONSUMPTION

  Fetch the [kpt](https://googlecontainertools.github.io/kpt/) package of the solution:

  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/iam/kpt/service-account service-account
  ```

# REQUIREMENTS

  A working Config Connector instance using the "cnrm-system" service account
  with either `roles/iam.serviceAccountAdmin` or `roles/owner` in the project
  managed by Config Connector.

# USAGE
  Replace `${IAM_MEMBER?}` with the GCP identity to grant access to:
  ```
  kpt cfg set . iam-member user:name@example.com
  ```
  
  _Optionally_, you can change the following fields before you apply the YAMLs: 
  - the name of the service account:
  ```
  kpt cfg set . service-account-name VALUE
  ```
  - the role granted to the GCP identity.
  (you can find all of the service account related IAM roles
  [here](https://cloud.google.com/iam/docs/understanding-roles#service-accounts-roles)):

  ```
  kpt cfg set . role roles/iam.serviceAccountTokenCreator
  ```

  Apply the YAMLs:

  ```
  kubectl apply -f .
  ```

# LICENSE

  Apache 2.0 - See [LICENSE](/LICENSE) for more information.
