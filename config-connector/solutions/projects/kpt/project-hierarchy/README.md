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
  with your billing account and organization ID.
  You will need to replace the project ID, since a
  project with the given ID already exists.


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

  Now the folder number can be found with the
  following.
  ```
  kubectl describe -f folder.yaml | grep Name:\ *folders\/ | sed "s/.*folders\///"
  ```
  You can set the folder number using the
  following kpt command:
  ```
  kpt cfg set . folder-number <number found above> --set-by "README-instructions"
  ```


  Now you can create your project under the new
  folder by applying the updated project YAML.
  ```
  kubectl apply -f project.yaml
  ```
## License

Apache 2.0 - See [LICENSE](LICENSE) for more information.
