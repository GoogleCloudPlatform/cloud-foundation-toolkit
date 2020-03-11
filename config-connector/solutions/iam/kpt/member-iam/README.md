Member IAM
==================================================

# NAME

  member-iam

# SYNOPSIS

  Config Connector compatible YAML files to create a service account in your desired project, and grant it a specific role (defaults to `compute.networkAdmin`) in the project.

  You can also create any number of service accounts and grant them any number of roles by altering and applying the altered config YAMLs.

# CONSUMPTION

  Using [kpt](https://googlecontainertools.github.io/kpt/):

  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/iam/kpt/member-iam member-iam
  ```

# REQUIREMENTS

  A working Config Connector cluster using a service account 
  (`cnrm-system@[PROJECT_ID].iam.gserviceaccount.com`) with the following 
  roles in your desired project (it doesn't need to be the project where you 
  installed Config Connector):

  - roles/resourcemanager.projectIamAdmin
  - roles/iam.serviceAccountAdmin

# USAGE

  Replace `${PROJECT_ID?}` with your desired project ID value from 
  within this directory:

  ```
  kpt cfg set . project-id VALUE
  ```

  Once the project ID is set in the configs, simply apply the YAMLs:

  ```
  kubectl apply -f .
  ```

## OPTIONAL WORKFLOW

  Optionally, you may want to create a different service account, grant a 
  different role to the service account, or grant multiple roles to the 
  service account, etc. Here are the instructions.

**Create a Different Service Account**

  Change the service account name and apply the YAMLs:

  ```
  kpt cfg set . service-account-name VALUE
  kubectl apply -f .
  ```

**Grant a Different Role**
  
  Change the role([predefined GCP IAM roles](https://cloud.google.com/iam/docs/understanding-roles#predefined_roles)) and apply the YAMLs:

  ```
  kpt cfg set . role VALUE
  kubectl apply -f .
  ```

**Grant Multiple Roles**

  Change the KCC resource name and the role each time before you grant a new
  role to the service account:

  ```
  kpt cfg set . policy-member-name VALUE_1
  kpt cfg set . role VALUE_2
  kubectl apply -f iampolicymember.yaml
  ```

# LICENSE

  Apache 2.0 - See [LICENSE](/LICENSE) for more information.