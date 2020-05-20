Project IAM
==================================================

# NAME
  project-iam
# SYNOPSIS
  Config Connector compatible YAML files to grant a role for a member in a project.
# CONSUMPTION
  Download the package using [kpt](https://googlecontainertools.github.io/kpt/):
  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/iam/kpt/project-iam project-iam
  ```
# REQUIREMENTS
  *   A working Config Connector cluster using "cnrm-system" service account
      that has the `roles/resourcemanager.projectIamAdmin` role in your desired
      project (it doesn't need to be the project managed by Config Connector).
  *   The project managed by Config Connector has Cloud Resource Manager API
      enabled.
# SETTERS
|    NAME    |        VALUE         |     SET BY      |       DESCRIPTION        | COUNT |
|------------|----------------------|-----------------|--------------------------|-------|
| member     | ${IAM_MEMBER?}       | PLACEHOLDER     | IAM member to grant role | 1     |
| project-id | ${PROJECT_ID?}       | PLACEHOLDER     | ID of project            | 1     |
| role       | roles/logging.viewer | package-default | IAM role to grant        | 1     |
# USAGE
Setters marked as `PLACEHOLDER` are required. Set them using kpt:
```
kpt cfg set . member user:name@example.com
kpt cfg set . project-id your-project
```
_Optionally_ set the role to grant in the same manner.

Once the configuration is satisfactory, apply the YAML:
```
kubectl apply -f .
```
# LICENSE
  Apache 2.0 - See [LICENSE](/LICENSE) for more information.
