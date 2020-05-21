Project Hierarchy
==================================================

# NAME

  project-hierarchy

# SYNOPSIS

  Config Connector compatible YAML files to create
  a folder in an organization, and a project
  beneath it.

# CONSUMPTION

  Using [kpt](https://googlecontainertools.github.io/kpt/):

  Run `kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/projects/kpt/project-hierarchy project-hierarchy`.

# REQUIREMENTS

  A working Config Connector cluster using a
  service account with the following roles in
  the organization:
  - `roles/resourcemanager.folderCreator`
  - `roles/resourcemanager.projectCreator`

# USAGE

  Replace the
  `${BILLING_ACCOUNT_ID?}` and `${ORG_ID?}` values:

  From within this directory, run
  ```
  kpt cfg set . billing-account VALUE
  ```
  and
  ```
  kpt cfg set . org-id VALUE
  ```
  replacing `VALUE` with your billing account
  and organization ID respectively.

  You will also need to reset the project ID,
  since a project with the given ID already exists.
  ```
  kpt cfg set . project-id VALUE
  ```


  Currently, to create a project under a
  folder, you must supply a numeric folder ID,
  which is only available after the folder is
  created. An issue outlining this shortfall in
  Config Connector functionality is filed on the
  project's GitHub,
  https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/104.


  To be nested beneath it, the project still needs
  a folder number. This can only be found after
  creating the folder. You can do so with
  ```
  kubectl apply -f folder.yaml
  ```

  Wait for GCP to generate the folder.
  ```
  kubectl wait --for=condition=Ready -f folder.yaml
  ```

  Now extract the folder number.
  ```
  FOLDER_NUMBER=$(kubectl describe -f folder.yaml | grep Name:\ *folders\/ | sed "s/.*folders\///")
  ```
  You can set the folder number using the
  following command:
  ```
  kpt cfg set . folder-number $FOLDER_NUMBER --set-by "README-instructions"
  ```


  Now you can fully apply this solution.
  ```
  kubectl apply -f .
  ```

# LICENSE

Apache 2.0 - See [LICENSE](/LICENSE) for more information.
