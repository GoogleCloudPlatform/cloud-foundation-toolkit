# Forseti Security

This template creates a new project (or reuses an existing one) and deploys
the [Forseti Security](https://forsetisecurity.org/) solution in it. The
solution consists of two instances (the client and the server), two service
accounts for them, a Cloud SQL instance, and the firewall rules with the IAM
policies for all of the above to function properly. See the Forseti Security
[manual](https://forsetisecurity.org/docs/v2.0/setup/manual.html)
installation page for specific configuration changes.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the following roles to Deployment Manager's service account at the
  organization level:
  - Billing Account Administrator
  - Org Admin
  - Storage Admin
  - Cloud SQL Admin
- Enable the [Cloud SQL Admin API](https://cloud.google.com/sql/docs/mysql/admin-api/)
- When reusing an existing project, enable the following APIs:
  - api-admin.googleapis.com
  - api-appengine.googleapis.com
  - api-bigquery-json.googleapis.com
  - api-cloudbilling.googleapis.com
  - api-cloudresourcemanager.googleapis.com
  - api-compute.googleapis.com
  - api-deploymentmanager.googleapis.com
  - api-iam.googleapis.com
  - api-sql-component.googleapis.com
  - api-sqladmin.googleapis.com
- Create a Cloud Storage bucket containing the `forseti_conf_server.yaml` file
  in the configs folder; see for [example](https://github.com/GoogleCloudPlatform/forseti-security/blob/stable/configs/server/forseti_conf_server.yaml.sample).

## Deployment

### Resources

- [compute.v1.instance](https://cloud.google.com/compute/docs/reference/rest/v1/instances)
- [compute.v1.network](https://cloud.google.com/compute/docs/reference/rest/v1/networks)
- [compute.v1.firewall](https://cloud.google.com/compute/docs/reference/rest/v1/firewalls)
- [cloudresourcemanager.v1.project](https://cloud.google.com/resource-manager/reference/rest/v1/projects)
- [sqladmin.v1beta4.instance](https://cloud.google.com/sql/docs/mysql/admin-api/v1beta4/databases)
- [sqladmin.v1beta4.database](https://cloud.google.com/sql/docs/mysql/admin-api/v1beta4/instances)
- [iam.v1.serviceAccount](https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts)

### Properties

See the `properties` section in the schema file(s):

- [Forseti](forseti.py.schema)

### Usage

1. Clone the [Deployment Manager Samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples):

    ```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
    ```

2. Go to the [community/cloud-foundation](../../) directory:

    ```shell
    cd community/cloud-foundation
    ```

3. Copy the example DM config to be used as a model for the deployment; in this
   case, [examples/forseti.yaml](examples/forseti.yaml):

    ```shell
    cp templates/forseti/examples/forseti.yaml my_forseti.yaml
    ```

4. Change the values in the config file to match your specific GCP setup (for
   properties, refer to the schema files listed above):

    ```shell
    vim my_forseti.yaml  # <== change values to match your GCP setup
    ```

5. Create your deployment (replace \<YOUR\_DEPLOYMENT\_NAME\> with the relevant
   deployment name):

    ```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_forseti.yaml
    ```

6. In case you need to delete your deployment:

    ```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
    ```

## Examples

- [Forseti](examples/forseti.yaml)
