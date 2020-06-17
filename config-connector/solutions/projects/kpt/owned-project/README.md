Owned Project
==================================================

# NAME

  owned-project

# SYNOPSIS

  Config Connector compatible YAML files to create
  a project in a folder, binding an IAM member
  as project owner.

# CONSUMPTION

  Fetch the package using [kpt](https://googlecontainertools.github.io/kpt/).

  `kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/projects/kpt/owned-project owned-project`

# REQUIREMENTS

  -   A working Config Connector cluster using the cnrm-system service account
      with the following roles in the target folder:
      -   `roles/resourcemanager.folderViewer`
      -   `roles/resourcemanager.projectCreator`
      -   `roles/iam.securityAdmin`
  -   The IAM member meets the requirements specified
      [here](https://cloud.google.com/resource-manager/reference/rest/v1/projects/setIamPolicy#top_of_page).

# USAGE

  Replace the
  `${BILLING_ACCOUNT_ID?}` value.
  From within this directory, run
  ```
  kpt cfg set . billing-account VALUE
  ```
  replacing `VALUE` with your billing account
  ID.

  Replace the `${FOLDER_ID?}`, `${IAM_MEMBER?}`, and `${PROJECT_ID?}` values the same way, using:
  ```
  kpt cfg set . folder-id VALUE
  kpt cfg set . iam-member VALUE
  kpt cfg set . project-id VALUE
  ```
  where the folder-id `VALUE` is the numeric folder ID of the folder to create the new project under, the iam-member `VALUE` is the fully qualified IAM name of target member, e.g. "user:me@example.com", and the project-id `VALUE` is the globally unique name you want your project to have.

  Now you can fully apply this solution.
  ```
  kubectl apply -f .
  ```

# LICENSE

Apache 2.0 - See [LICENSE](/LICENSE) for more information.
