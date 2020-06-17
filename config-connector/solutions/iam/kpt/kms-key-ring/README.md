KMS Key Ring
==================================================

# NAME
  kms-key-ring
# SYNOPSIS
  Config Connector compatible yaml files for creating a kms key ring and applying a role to it.
# CONSUMPTION
  Download the package using [kpt](https://googlecontainertools.github.io/kpt/).
  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/iam/kpt/kms-key-ring kms-key-ring
  ```
# REQUIREMENTS
  -   A working Config Connector instance using the "cnrm-system" service
      account with either `roles/cloudkms.admin` or `roles/owner` in the project
      managed by Config Connector.
  -   Cloud Key Management Service (KMS) API enabled in the project where Config
      Connector is installed
  -   Cloud Key Management Service (KMS) API enabled in the project managed by
      Config Connector if it is a different project

# SETTERS
|    NAME    |        VALUE         |     SET BY      |     DESCRIPTION      | COUNT |
|------------|----------------------|-----------------|----------------------|-------|
| iam-member | ${IAM_MEMBER?}       | PLACEHOLDER     | member to grant role | 1     |
| location   | us-central1          | package-default | location of key ring | 1     |
| ring-name  | allowed-ring         | package-default | name of key ring     | 2     |
| role       | roles/cloudkms.admin | package-default | IAM role to grant    | 1     |
# USAGE
  Set the IAM member that you would like to apply a role to.
  ```
  kpt cfg set . iam-member user:name@example.com
  ```
  _Optionally_ set the name of the KMS keyring (defaults to `allowed-ring`).
  ```
  kpt cfg set . ring-name your-ring-name
  ```
  _Optionally_ set the [IAM role](https://cloud.google.com/iam/docs/understanding-roles#cloud-kms-roles) to grant (defaults to `roles/cloudkms.admin`).
  ```
  kpt cfg set . role roles/cloudkms.importer
  ```
  _Optionally_ set the location of the ring (defaults to `us-central1`)
  ```
  kpt cfg set . location us-west1
  ```
# LICENSE
  Apache 2.0 - See [LICENSE](/LICENSE) for more information.
