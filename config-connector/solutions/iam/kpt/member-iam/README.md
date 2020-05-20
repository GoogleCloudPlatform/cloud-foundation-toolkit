Member IAM
==================================================

# NAME

  member-iam

# SYNOPSIS

  Config Connector compatible YAML files to create a service account in your desired project, and grant it a specific role (defaults to `compute.networkAdmin`) in the project.

# CONSUMPTION

  Using [kpt](https://googlecontainertools.github.io/kpt/):

  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/iam/kpt/member-iam member-iam
  ```

# REQUIREMENTS

  *   A working Config Connector cluster using "cnrm-system" service account
      that has the following roles in your desired project (it doesn't need to
      be the project managed by Config Connector):

      -   roles/resourcemanager.projectIamAdmin
      -   roles/iam.serviceAccountAdmin

  *   Cloud Resource Manager API enabled in the project where Config Connector
      is installed

# USAGE

  Replace `${PROJECT_ID?}` with your desired project ID value from 
  within this directory:

  ```
  kpt cfg set . project-id VALUE
  ```

  _Optionally_, you can also change the service account name and role
  (you can find all the predefined GCP IAM roles
  [here](https://cloud.google.com/iam/docs/understanding-roles#predefined_roles)):

  ```
  kpt cfg set . service-account-name VALUE
  kpt cfg set . role VALUE
  ```

  Once the fields are set in the configs, apply the YAMLs:

  ```
  kubectl apply -f .
  ```

  You can check the resources you just created:

  ```
  kubectl get iamserviceaccount --namespace member-iam-solution
  kubectl get iampolicymember --namespace member-iam-solution
  ```

# LICENSE

  Apache 2.0 - See [LICENSE](/LICENSE) for more information.