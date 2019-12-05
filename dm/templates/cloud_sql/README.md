# Cloud SQL

This template creates a Cloud SQL instance with databases and users.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Enable the [Cloud SQL API](https://cloud.google.com/sql/docs/mysql/admin-api/)
- Enable the [Cloud SQL Admin API](https://cloud.google.com/sql/docs/mysql/admin-api/)
- Grant the [roles/cloudsql.admin](https://cloud.google.com/sql/docs/mysql/project-access-control)
  IAM role to the Deployment Manager service account

## Deployment

### Resources

- [sqladmin.v1beta4.instance](https://cloud.google.com/sql/docs/mysql/admin-api/v1beta4/instances)
- [sqladmin.v1beta4.database](https://cloud.google.com/sql/docs/mysql/admin-api/v1beta4/databases)
- [sqladmin.v1beta4.user](https://cloud.google.com/sql/docs/mysql/admin-api/v1beta4/users)

### Properties

See the `properties` section in the schema file(s):

- [Cloud SQL](cloud_sql.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit):

    ```shell
    git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit
    ```

2. Go to the [dm](../../) directory:

    ```shell
    cd dm
    ```

3. Copy the example DM config to be used as a model for the deployment; in this
   case, [examples/cloud\_sql.yaml](examples/cloud_sql.yaml):

    ```shell
    cp templates/cloud_sql/examples/cloud_sql.yaml my_cloud_sql.yaml
    ```

4. Change the values in the config file to match your specific GCP setup (for
   properties, refer to the schema files listed above):

    ```shell
    vim my_cloud_sql.yaml  # <== change values to match your GCP setup
    ```

5. Create your deployment (replace \<YOUR\_DEPLOYMENT\_NAME\> with the relevant
   deployment name):

    ```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_cloud_sql.yaml
    ```

   To deploy with CFT:

    ```shell
    cft apply my_cloud_sql.yaml
    ```

6. In case you need to delete your deployment:

    ```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
    ```

   To delete deployment with CFT:

    ```shell
    cft delete my_cloud_sql.yaml
    ```

`Notes:` After a Cloud SQL instance is deleted, its name cannot be reused for
up to 7 days.

## Examples

- [Cloud SQL](examples/cloud_sql.yaml)
- [Cloud SQL with Read Replica](examples/cloud_sql_read_replica.yaml)
- [Cloud SQL Postgres](examples/cloud_sql_postgres.yaml)
- [Cloud SQL MSSQL Server](examples/cloud_sql_mssql.yaml)
- [Cloud SQL Private Networking](examples/cloud_sql_private_network.yaml)
