IAM Project
==================================================

# NAME

  iam-project

# SYNOPSIS

  Config Connector compatible YAML files to create
  a project in a folder, binding an IAM member
  as project owner.

# CONSUMPTION

  Using [kpt](https://googlecontainertools.github.io/kpt/):

  Run `kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/projects/kpt/project-owner project-owner`.

# REQUIREMENTS

  A working Config Connector cluster using a
  service account with the following roles in
  the organization:
  - roles/resourcemanager.folderViewer
  - roles/resourcemanager.projectCreator
  - roles/iam.securityAdmin

# USAGE

  Replace the
  `${BILLING_ACCOUNT_ID?}` value:

  From within this directory, run
  ```
  kpt cfg set . billing-account VALUE
  ```
  replacing `VALUE` with your billing account
  ID.

  Replace the `${FOLDER_ID}` the same way, using:
  ```
  kpt cfg set . folder-id VALUE
  ```
  where VALUE is the numeric folder ID of the folder to create the new project under.

  You will need to reset the project ID,
  since a project with the given ID already exists.
  ```
  kpt cfg set . project-id VALUE
  ```

  To change the IAM member owning the project.
  ```
  kpt cfg set . iam-name VALUE
  ```
  where VALUE is the fully qualified IAM name of target member, e.g. "user:me@example.com".

  Now you can fully apply this solution.
  ```
  kubectl apply -f .
  ```

# LICENSE

Apache 2.0 - See [LICENSE](/LICENSE) for more information.
