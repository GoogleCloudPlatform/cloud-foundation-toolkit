# BigQuery

This template creates a BigQuery dataset and table.

## Prerequisites
- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [`roles/bigquery.dataEditor`, `roles/bigquery.dataOwner` or `roles/bigquery.admin`](https://cloud.google.com/bigquery/docs/access-control) IAM roles to the project service account

## Deployment

### Resources

- [bigquery.v2.dataset](https://cloud.google.com/bigquery/docs/reference/rest/v2/datasets)
- [bigquery.v2.tables](https://cloud.google.com/bigquery/docs/reference/rest/v2/tables)

### Properties

See the `properties` section in the schema file(s)

- [BigQuery Dataset](bigquery_dataset.py.schema)
- [BigQuery Tables](bigquery_table.py.schema)

### Usage

1. Clone the [Deployment Manager Samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples)

```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [community/cloud-foundation](../../) directory

```shell
    cd community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment, in this case [examples/bigquery.yaml](examples/bigquery.yaml)

```shell
    cp templates/bigquery/examples/bigquery.yaml my_bigquery.yaml
```

4. Change the values in the config file to match your specific GCP setup.
   Refer to the properties in the schema files described above.

```shell
    vim my_bigquery.yaml  # <== change values to match your GCP setup
```

5. Create your deployment as described below, replacing <YOUR_DEPLOYMENT_NAME>
   with your with your own deployment name

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_bigquery.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Bigquery Dataset and Table](examples/bigquery.yaml)
