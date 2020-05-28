PostgreSQL Public
==================================================
# NAME
  postgres-public
# SYNOPSIS
  Config Connector compatible yaml files to configure a public PostgreSQL database
# CONSUMPTION
  Download the package using [kpt](https://googlecontainertools.github.io/kpt/).
  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/sql/kpt/postgres-public postgres-public
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
|          NAME           |              VALUE              |     SET BY      |          DESCRIPTION          | COUNT |
|-------------------------|---------------------------------|-----------------|-------------------------------|-------|
| authorized-network      | postgres-public-solution-sample | package-default | name of authorized network    | 1     |
| authorized-network-cidr | 130.211.0.0/28                  | package-default | authorized network CIDR range | 1     |
| instance-name           | postgres-ha-solution            | package-default | name of SQL instance          | 3     |
| password                | ${PASSWORD?}                    | PLACEHOLDER     | password for SQL user         | 1     |
# USAGE
  Configure setters using kpt as follows:
  ```
  kpt cfg set . NAME VALUE
  ```
  Setting placeholder values is required, changing package-defaults is optional.

  `password` should be set to a [base64
encoded](https://kubernetes.io/docs/concepts/configuration/secret/#creating-a-secret-manually)
value.
  ```
  kpt cfg set . password $(echo -n 'password' | base64)
  ```
  _Optionally_ set `authorized-network`, `authorized-network-cidr`, and `instance-name` in the same manner.

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

