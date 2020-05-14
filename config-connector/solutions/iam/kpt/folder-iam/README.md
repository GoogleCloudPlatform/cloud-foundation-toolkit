Folder IAM
==================================================

# NAME

  folder-iam

# SYNOPSIS

  Config Connector compatible YAML files to grant a specific member a role (default to `roles/resourcemanager.folderEditor`) to an existing folder.

# CONSUMPTION

  Download the package using [kpt](https://googlecontainertools.github.io/kpt/):

  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/iam/kpt/folder-iam folder-iam
  ```

# REQUIREMENTS

  -   A working Config Connector instance using the "cnrm-system" service
      account with the following role in the desired folder:
      -   roles/resourcemanager.folderIamAdmin
  -   Cloud Resource Manager API enabled in the project where Config Connector
      is installed

# USAGE

  Replace the `${FOLDER_ID?}` with a folder ID you want to add member to:
  ```
  kpt cfg set . folder-id VALUE
  ```

  Replace the `${IAM_MEMBER?}` with a GCP identity to grant role to:
  ```
  kpt cfg set . iam-member VALUE
  ```

  _Optionally_, you can also change the role granted to the member. (you can find all of the folder related IAM roles
  [here](https://cloud.google.com/iam/docs/understanding-roles#resource-manager-roles)):

  ```
  kpt cfg set . role roles/resourcemanager.folderViewer
  ```

  Apply the YAMLs:

  ```
  kubectl apply -f .
  ```

# LICENSE

  Apache 2.0 - See [LICENSE](/LICENSE) for more information.
