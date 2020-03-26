Service Account
==================================================

# NAME

  service-account

# SYNOPSIS

  Config Connector compatible YAML files to create a service account in your desired project, and grant it a specific role (default to iam.serviceAccountKeyAdmin) on the service account resource level.

  You can also create any number of service accounts and grant them any number of roles by adding new or updating the config YAMLs and then applying them.

# CONSUMPTION

  Using [kpt](https://googlecontainertools.github.io/kpt/):

  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/iam/kpt/service-account service-account
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

  Change the role([predefined GCP IAM roles](https://cloud.google.com/iam/docs/understanding-roles#predefined_roles)) and apply the YAMLs. Please make sure to delete the iam policy created previously:

  ```
  kubectl delete -f .
  kpt cfg set . role VALUE
  kubectl apply -f .
  ```

**Grant Multiple Roles**

  Change the KCC resource name and the role each time before you grant a new
  role to the service account, also please make sure you delete the iam policy created previously:

  ```
  kubectl delete -f .
  kpt cfg set . iampolicy-name VALUE_1
  kpt cfg set . role VALUE_2
  kubectl apply -f iampolicy.yaml
  ```
**Alternative Usage**

You can modify the bindings object in the iampolicy.yaml to grant multiple roles to multiple identities at the same time. Here is an example:

```
  bindings:
    - role: roles/iam.serviceAccountKeyAdmin
      members:
        - serviceAccount:service-account-example-1@${PROJECT_ID}.iam.gserviceaccount.com
    - role: roles/iam.serviceAccountTokenCreator
      members:
        - serviceAccount:service-account-example-2@${PROJECT_ID}.iam.gserviceaccount.com
        - serviceAccount:service-account-example-1@${PROJECT_ID}.iam.gserviceaccount.com
```
In the above example, _service-account-example-1_ will have both of the iam.serviceAccountKeyAdmin role and roles/iam.serviceAccountTokenCreator, while _service-account-example-2_ will only have iam.serviceAccountTokenCreator role.


# LICENSE

  Apache 2.0 - See [LICENSE](/LICENSE) for more information.
