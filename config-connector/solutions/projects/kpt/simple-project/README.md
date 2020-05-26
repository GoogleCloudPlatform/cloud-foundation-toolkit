Simple Project
==================================================

# NAME

  simple-project

# SYNOPSIS

  Config Connector compatible YAML files to create a project in an organization.

# CONSUMPTION
  Using kpt:
  ```
  BASE=https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
  kpt pkg get $BASE/config-connector/solutions/projects/kpt/simple-project simple-project
  ```

# REQUIREMENTS
  A working cluster with Config Connector installed.
  
  Cloud Resource Manager API and Cloud Billing API enabled in the project where Config Connector is installed.
  
  The "cnrm-system" service account must have `roles/resourcemanager.projectCreator` in your organization and `roles/billing.user` for your billing account.
  
# USAGE
  In order to use, replace the `${PROJECT_ID?}`, `${BILLING_ACCOUNT_ID?}` and
  `${ORG_ID?}` values with a unique new project ID, your billing account and
  your organization id. You can do this with kpt setters:
  ```
  kpt cfg set . project-id VALUE
  kpt cfg set . billing-account VALUE 
  kpt cfg set . org-id VALUE 
  ```

  Note: Updating the project-id will set both the project's ID and name to the
  same value, if you want a different value for project name, edit
  `project.yaml` and replace spec.name with your preferred project name.
  
  Once your information is in the configs, simply apply.

  ```
  kubectl apply -f .
  ```
  
# License

  Apache 2.0 - See [LICENSE](/LICENSE) for more information.

