MySQL Private
==================================================
# NAME
  mysql-private
# SYNOPSIS
  Config Connector compatible YAML files for creating a MySQL instance on a private network
# REQUIREMENTS
  A working Config Connector installation managing a project with Cloud SQL Admin API, Service Networking API, and Cloud Resource Manager API enabled.
# CONSUMPTION
  Download the package using [kpt](https://googlecontainertools.github.io/kpt/).
  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/config-connector/solutions/sql/kpt/mysql-private mysql-private
  ```
# USAGE
## SETTERS
|     NAME      |         VALUE          |     SET BY      |      DESCRIPTION       | COUNT |
|---------------|------------------------|-----------------|------------------------|-------|
| database-name | mysql-private-database | package-default | name of SQL database   | 1     |
| instance-name | mysql-private-instance | package-default | name of SQL instance   | 3     |
| password      | ${PASSWORD?}           | PLACEHOLDER     | password of SQL user   | 1     |
| region        | us-central1            | package-default | region of SQL instance | 1     |
| user-name     | mysql-private-user     | package-default | name of SQL user       | 1     |

  Configure setters using kpt as follows:
  ```
  kpt cfg set . NAME VALUE
  ```
  Setting placeholder values is required, changing package-defaults is optional.

  For this package to work properly, the following resources must be in a ready state before the SQLInstance YAML is applied:
  - `ComputeNetwork`
  - `ComputeAddress`
  - `ServiceNetworkingConnection`
  
  To ensure this is the case, use the following:
  ```
  kubectl apply -f network
  kubectl wait --for=condition=Ready -f network 
  kubectl apply -f sql
  ```

# LICENSE
  Apache 2.0 - See [LICENSE](/LICENSE) for more information.

