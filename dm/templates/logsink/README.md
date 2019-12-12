# Logsink

This template creates a logsink (logging sink). The logsink destination can
exist prior to creating the logsink or can be created by the logsink template.
If the resources are created by the logsink, the logsink uniqueWriter service
account will be granted the appropriate permissions to the destination
resource.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Create one of the following:
  - [GCS bucket](https://cloud.google.com/storage/docs/json_api/v1/buckets)
  - [PubSub topic](https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.topics)
  - [BigQuery dataset](https://cloud.google.com/bigquery/docs/reference/rest/v2/datasets)
- Grant the [logging.configWriter or logging.admin](https://cloud.google.com/logging/docs/access-control)
  IAM role to the project service account
- Grant the [`pubsub.admin`](https://cloud.google.com/pubsub/docs/access-control)
  IAM role to the project service account if creating a pubsub logging sink
  destination
- Grant the [`storage.admin`](https://cloud.google.com/storage/docs/access-control/iam-roles)
  IAM role to the project service account if creating a bucket logging sink
  destination
- Grant the [`bigquery.admin`](https://cloud.google.com/bigquery/docs/access-control)
  IAM role to the project service account if creating bq logging sink
  destination

#### If you are going to create bucket, pubsub or BigQuery destinations in current project:  

- Grant the [resourcemanager.projectIamAdmin or owner](https://cloud.google.com/iam/docs/understanding-roles) IAM role on the project to the *DM Service Account* to grant roles within the project
- Grant the [roles/resourcemanager.folderIamAdmin owner](https://cloud.google.com/iam/docs/understanding-roles) IAM role on the folder to the *DM Service Account* to grant roles within the folder
- Grant the [roles/iam.securityAdmin or owner](https://cloud.google.com/iam/docs/understanding-roles) IAM role on the organization to the *DM Service Account* to grant roles within the organization and all nested resources
- Grant the [logging.configWriter or logging.admin](https://cloud.google.com/logging/docs/access-control) IAM role on the project to the *DM Service Account* to grant roles within the project

## If you specify destination project and are going to create bucket, pubsub or BigQuery destinations:

- Grant the [resourcemanager.projectIamAdmin or owner](https://cloud.google.com/iam/docs/understanding-roles) IAM role on the project to the *DM Service Account* to grant roles within the project
- Grant the [roles/resourcemanager.folderIamAdmin owner](https://cloud.google.com/iam/docs/understanding-roles) IAM role on the folder to the *DM Service Account* to grant roles within the folder
- Grant the [roles/iam.securityAdmin or owner](https://cloud.google.com/iam/docs/understanding-roles) IAM role on the organization to the *DM Service Account* to grant roles within the organization and all nested resources
- Grant the [logging.configWriter or logging.admin](https://cloud.google.com/logging/docs/access-control) IAM role on the project to the *DM Service Account* to grant roles within the project
  IAM role to the project service account
- Grant the [`pubsub.admin`](https://cloud.google.com/pubsub/docs/access-control)
  IAM role to the project service account if creating a pubsub logging sink
  destination
- Grant the [`storage.admin`](https://cloud.google.com/storage/docs/access-control/iam-roles)
  IAM role to the project service account if creating a bucket logging sink
  destination
- Grant the [`bigquery.admin`](https://cloud.google.com/bigquery/docs/access-control)
  IAM role to the project service account if creating bq logging sink
  destination

## Deployment

### Resources

- [logging.v2.sink](https://cloud.google.com/logging/docs/reference/v2/rest/v2/projects.sinks)
- [pubsub.v1.topic](https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.topics)
- [storage.v1.bucket](https://cloud.google.com/storage/docs/creating-buckets)
- [bigquery.v2.dataset](https://cloud.google.com/bigquery/docs/reference/rest/v2/datasets)

### Properties

See `properties` section in the schema file(s):

- [Logsink](logsink.py.schema)

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
   case, [examples/logsink.yaml](examples/logsink.yaml):

```shell
    cp templates/logsink/examples/logsink.yaml my_logsink.yaml
```

4. Change the values in the config file to match your specific GCP setup (for
   properties, refer to the schema files listed above):

```shell
    vim my_logsink.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace \<YOUR\_DEPLOYMENT\_NAME\> with the relevant
   deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
    --config my_logsink.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Organization logging entries exported to PubSub](examples/org_logsink_pubsub_destination.yaml)
- [Billing account logging entries exported to Storage](examples/billingaccount_logsink_bucket_destination.yaml)
- [Folder logging entries exported to BigQuery](examples/folder_logsink_bq_destination.yaml)
- [Project logging entries exported to Storage](examples/project_logsink_bucket_destination.yaml)
