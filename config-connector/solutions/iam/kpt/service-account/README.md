Service Account
==================================================

# NAME

  service-account

# SYNOPSIS

  Config Connector compatible YAML files to create a service account in your desired project, and grant a specifi member a role (default to iam.serviceAccountKeyAdmin) to access the service account that just created.

# CONSUMPTION

  Using [kpt](https://googlecontainertools.github.io/kpt/):

  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/iam/kpt/service-account service-account
  ```

# REQUIREMENTS

  A working Config Connector cluster using the cnrm-system service account.

# USAGE

  Replace the `${SERVICE_ACCOUNT?}` with a service account name you want to create:
  ```
  kpt cfg set . service-account-name VALUE
  ```

  Replace `${IAM_MEMBER?}` with the GCP identity to grant access to:
  ```
  kpt cfg set . iam-member user:name@example.com
  ```

  _Optionally_, you can also change the role granted to the GCP identity in the previous step.
  (you can find all of the service account related IAM roles
  [here](https://cloud.google.com/iam/docs/understanding-roles#service-accounts-roles)):

  ```
  kpt cfg set . role VALUE
  ```

  Apply the YAMLs:

  ```
  kubectl apply -f .
  ```

# LICENSE

  Apache 2.0 - See [LICENSE](/LICENSE) for more information.
