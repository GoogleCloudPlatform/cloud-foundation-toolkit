KMS Crypto Key
==================================================
# NAME
  kms-crypto-key
# SYNOPSIS
  Config Connector compatible yaml files to create a kms key ring, a kms crypto key,
  and apply an IAM role to the crypto key.
# CONSUMPTION
  Download the package using [kpt](https://googlecontainertools.github.io/kpt/).
  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/iam/kpt/kms-crypto-key kms-crypto-key
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
|    NAME    |         VALUE         |     SET BY      |       DESCRIPTION        | COUNT |
|------------|-----------------------|-----------------|--------------------------|-------|
| iam-member | ${IAM_MEMBER?}        | PLACEHOLDER     | IAM member to grant role | 1     |
| key-name   | allowed-key           | package-default | name of key              | 2     |
| location   | us-central1           | package-default | location of ring         | 1     |
| ring-name  | allowed-ring          | package-default | name of ring             | 2     |
| role       | roles/cloudkms.signer | package-default | IAM role to grant        | 1     |
# USAGE
  Set the IAM member to apply a role to:
  ```
  kpt cfg set . iam-member user:name@example.com
  ```
  _Optionally_ set the role to apply:
  ```
  kpt cfg set . role roles/cloudkms.cryptoKeyDecrypter
  ```
  _Optionally_ set the crypto key name, key ring name, and key ring location:
  ```
  kpt cfg set . key-name your-key
  kpt cfg set . ring-name your-ring
  kpt cfg set . location us-west1
  ```
  Once the values are satisfactory, apply:
  ```
  kubectl apply -f .
  ```
# LICENSE
  Apache 2.0 - See [LICENSE](/LICENSE) for more information.

