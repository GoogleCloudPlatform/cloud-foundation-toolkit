MySQL Private
==================================================
# NAME
  mysql-private
# SYNOPSIS
  Config Connector compatible YAML files for creating a MySQL instance on a private network
# CONSUMPTION
  Download the package using [kpt](https://googlecontainertools.github.io/kpt/).
  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/sql/kpt/mysql-private mysql-private
  ```

# REQUIREMENTS
  -   A working Config Connector instance using the "cnrm-system" service
      account with either both `roles/cloudsql.admin` and
      `roles/compute.networkAdmin` roles or `roles/owner` in the project managed
      by Config Connector
  -   The following APIs enabled in the project where Config Connector is
      installed:
      -   Cloud SQL Admin API
      -   Compute Engine API
  -   The following APIs enabled in the project managed by Config Connector:
      -   Cloud SQL Admin API
      -   Compute Engine API
      -   Service Networking API

# SETTERS
|     NAME      |         VALUE          |     SET BY      |          DESCRIPTION          | COUNT |
|---------------|------------------------|-----------------|-------------------------------|-------|
| database-name | mysql-private-database | package-default | name of SQL database          | 1     |
| instance-name | mysql-private-instance | package-default | name of SQL instance          | 3     |
| password      | ${PASSWORD?}           | PLACEHOLDER     | SQL password (base64 encoded) | 1     |
| region        | us-central1            | package-default | region of SQL instance        | 1     |
| username      | ${USERNAME?}           | PLACEHOLDER     | name of SQL user              | 1     |

# USAGE

  Configure setters using kpt as follows:
  ```
  kpt cfg set . NAME VALUE
  ```
  Setting placeholder values is required, changing package-defaults is optional.

  Set `username` to the SQL username that you will use to access the database.
  ```
  kpt cfg set . username your-username
  ```
  _Optionally_ set `database-name`, `instance-name`, and `region` in the same
manner. Note that if your instance is deleted the name you used will be
reserved for 7 days. You will need to use a new name in order to re-create the
instance.

  `password` should be set to a [base64 encoded](https://kubernetes.io/docs/concepts/configuration/secret/#creating-a-secret-manually) value.
  ```
  kpt cfg set . password $(echo -n 'your-password' | base64)
  ```
  Due to the bug in Config Connector ([more details](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/148)), the following resources must be in a ready state before the SQLInstance YAML is applied:
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
