PostgreSQL High Availability
==================================================
# NAME
  postgres-ha
# SYNOPSIS
  Config Connector compatible yaml files to configure a high availability PostgreSQL cluster
# CONSUMPTION
  Download the package using [kpt](https://googlecontainertools.github.io/kpt/).
  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/sql/kpt/postgres-ha postgres-ha
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
|       NAME        |         VALUE          |     SET BY      |          DESCRIPTION           | COUNT |
|-------------------|------------------------|-----------------|--------------------------------|-------|
| database-1-name   | postgres-ha-database-1 | package-default | name of first SQL database     | 1     |
| database-2-name   | postgres-ha-database-2 | package-default | name of second SQL database    | 1     |
| external-ip-range | 192.10.10.10/32        | package-default | ip range to allow to connect   | 4     |
| instance-name     | postgres-ha-solution   | package-default | name of main SQL instance      | 9     |
| password-1        | ${PASSWORD_1?}         | PLACEHOLDER     | password of  user              | 1     |
| password-2        | ${PASSWORD_2?}         | PLACEHOLDER     | password of  user              | 1     |
| password-3        | ${PASSWORD_3?}         | PLACEHOLDER     | password of  user              | 1     |
| region            | us-central1            | package-default | region of SQL instance         | 4     |
| username-1        | ${USERNAME_1?}         | PLACEHOLDER     | name of  user                  | 1     |
| username-2        | ${USERNAME_2?}         | PLACEHOLDER     | name of  user                  | 1     |
| username-3        | ${USERNAME_3?}         | PLACEHOLDER     | name of  user                  | 1     |
| zone              | us-central1-c          | package-default | zone of main instance          | 1     |
| zone-replica-1    | us-central1-a          | package-default | zone of first replica instance | 1     |
| zone-replica-2    | us-central1-b          | package-default | zone of second replica instance| 1     |
| zone-replica-3    | us-central1-c          | package-default | zone of third replica instance | 1     |
# USAGE
  Configure setters using kpt as follows:
  ```
  kpt cfg set . NAME VALUE
  ```
  Setting placeholder values is required, changing package-defaults is optional.

  Set `username-1`, `username-2', and `username-3` to the SQL usernames that you will use to access the database.
  ```
  kpt cfg set . username-1 first-username
  kpt cfg set . username-2 second-username
  kpt cfg set . username-3 third-username
  ```
  `password-1`, `password-2`, and `password-3` should be set to [base64
encoded](https://kubernetes.io/docs/concepts/configuration/secret/#creating-a-secret-manually)
values.
  ```
  kpt cfg set . password-1 $(echo -n 'first-password' | base64)
  kpt cfg set . password-2 $(echo -n 'second-password' | base64)
  kpt cfg set . password-3 $(echo -n 'third-password' | base64)
  ```
  _Optionally_ set `database-name`, `instance-name`, `region`, `zone`, and
`zone-replica` in the same manner.

  **Note:** If your SQL Instance is deleted, the name you used will be reserved
for **7 days**. In order to re-apply this solution, you need to run
`kpt cfg set . instance-name new-instance-name` to change to a new
instance name that hasn't been used in the last 7 days.
 
  Once the configuration is satisfactory, apply:
  ```
  kubectl apply -f .
  ```
# LICENSE
  Apache 2.0 - See [LICENSE](/LICENSE) for more information.

