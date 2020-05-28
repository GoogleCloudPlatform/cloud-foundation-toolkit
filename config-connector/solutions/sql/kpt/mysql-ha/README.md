MySQL High Availability
==================================================
# NAME
  mysql-ha
# SYNOPSIS
  Config Connector compatible YAMLs for creating a high availability MySQL cluster
# CONSUMPTION
  Download the package using [kpt](https://googlecontainertools.github.io/kpt/):
  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/sql/kpt/mysql-ha mysql-ha
  ```
# REQUIREMENTS
  -   A working Config Connector instance using the "cnrm-system" service
      account with either `roles/cloudsql.admin` or `roles/owner` in the project
      managed by Config Connector
  -   Cloud SQL Admin API enabled in the project where Config Connector is
      installed
  -   Cloud SQL Admin API enabled in the project managed by Config Connector if
      it is a different project

# SETTERS
|       NAME        |        VALUE        |     SET BY      |          DESCRIPTION          | COUNT |
|-------------------|---------------------|-----------------|-------------------------------|-------|
| instance-name     | mysql-ha-solution   | package-default | name of SQL instance          | 14    |
| test-pw           | ${PASSWORD_1?}      | PLACEHOLDER     | password for SQL user "test"<br>(base64 encoded)  | 1     |
| test2-pw          | ${PASSWORD_2?}      | PLACEHOLDER     | password for SQL user "test2"<br>(base64 encoded) | 1     |
| test3-pw          | ${PASSWORD_3?}      | PLACEHOLDER     | password for SQL user "test3"<br>(base64 encoded) | 1     |

# USAGE
  Configure setters using kpt as follows:
  ```
  kpt cfg set . NAME VALUE
  ```
  Setting placeholder values is required, changing package-defaults is optional.

  Set `test-pw`, `test2-pw`, and `test3-pw` to the [base64
  encoded](https://kubernetes.io/docs/concepts/configuration/secret/#creating-a-secret-manually) passwords for user `test`,
  user `test2`, and user `test3`:
  ```
  kpt cfg set . test-pw $(echo -n 'first-password' | base64)
  kpt cfg set . test2-pw $(echo -n 'second-password' | base64)
  kpt cfg set . test3-pw $(echo -n 'third-password' | base64)
  ```
  _Optionally,_ set `instance-name` in the same manner.

  **Note:** If your SQL Instance is deleted, the name you used will be reserved
  for **7 days**. In order to re-apply this solution, you need to run
  `kpt cfg set . instance-name new-instance-name` to change to a new
  instance name that hasn't been used in the last 7 days.
 
  Once the configuration is satisfactory, apply:
  ```
  kubectl apply -f .
  ```
  **Note:** It will take up to ~40 mins for all the resources to be `Ready`.
  
# LICENSE
  Apache 2.0 - See [LICENSE](/LICENSE) for more information.

