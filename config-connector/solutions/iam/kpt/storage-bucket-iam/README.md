Storage Bucket IAM
==================================================
# NAME
  storage-bucket-iam
# SYNOPSIS
  Config Connector compatible yaml to enable permissions for a storage bucket.
# CONSUMPTION
  Download the package using [kpt](https://googlecontainertools.github.io/kpt/).
  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/iam/kpt/storage-bucket-iam storage-bucket-iam
  ```
# REQUIREMENTS
- A working Config Connector installation.
- The cnrm-system service account must have `roles/storage.admin`.
- A storage bucket managed by IAM.

  Note: Using uniform bucket-level access control is recommended for this package.
# SETTERS
|    NAME     |           VALUE            |     SET BY      |      DESCRIPTION       | COUNT |
|-------------|----------------------------|-----------------|------------------------|-------|
| bucket-name | ${BUCKET_NAME?}            | PLACEHOLDER     | name of storage bucket | 1     |
| iam-member  | ${IAM_MEMBER?}             | PLACEHOLDER     | member to grant role   | 1     |
| role        | roles/storage.objectViewer | package-default | IAM role to grant      | 1     |
# USAGE
  Set the name of the bucket you want to configure permissions for.
  ```
  kpt cfg set . bucket-name your-bucket
  ```
  Set the IAM member to grant a role to.
  ```
  kpt cfg set . iam-member user:name@example.com
  ```
  Optionally, set the [storage 
  role](https://cloud.google.com/iam/docs/understanding-roles#storage-roles) (defaults to
  `roles/storage.objectViewer`) that you want to apply and the IAM member the role will apply to.
  ```
  kpt cfg set . role roles/storage.admin
  ```
  Once the configuration is satisfactory, apply:
  ```
  kubectl apply -f .
  ```
# LICENSE
  Apache 2.0 - See [LICENSE](/LICENSE) for more information.
