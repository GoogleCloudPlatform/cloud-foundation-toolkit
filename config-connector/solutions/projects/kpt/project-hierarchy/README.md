Project Hierarchy
==================================================

# NAME

  project-hierarchy

# SYNOPSIS

  Config Connector compatible YAML files to create
  a folder in an organization, and a project
  beneath it. Make sure your Config Connector
  service account has the folder creator and
  project creator role for your organization.

  In order to use, replace the
  `${BILLING_ACCOUNT_ID?}` and `${ORG_ID?}` values
  with your billing account and organization id
  numbers. Make sure your Config Connector service
  account has the folder creator role for your
  organization.

  To be nested beneath it, the project still needs
  a folder number. This can only be found after
  creating the folder. You can do so with
  ```
  kubectl apply -f folder.yaml
  ```

  Now the folder number can be found and set by
  executing `set_folder_number.sh`. Create your
  project under the new folder by applying the
  updated project YAML.
  ```
  kubectl apply -f project.yaml
  ```
=======
  Ensure your Config Connector service account has permissions for folder and project

