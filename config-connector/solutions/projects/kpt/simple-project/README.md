Simple Project
==================================================

# NAME

  simple-project

# SYNOPSIS

  Config Connector compatible YAML files to create
  a project in an organization. Make sure your Config Connector
  service account has the project creator role for your organization.

  In order to use, replace the `${BILLING_ACCOUNT_ID?}` and `${ORG_ID?}` values
  with your billing account and organization id numbers. You can do this with kpt setters or manually. To use kpt setters:
  
  ```
  kpt cfg set . billing-account <insert billing account here>
  kpt cfg set . org-id <insert organization id here>
  ```
  
  Once your information is in the configs, simply apply.

  ```
  kubectl apply -f .
  ```
  
# License

  Apache 2.0 - See [LICENSE](/LICENSE) for more information.

