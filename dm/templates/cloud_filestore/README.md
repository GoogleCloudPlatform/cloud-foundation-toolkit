# Cloud Filestore

This template creates a Cloud Filestore instance.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Enable the [Cloud Build API](https://cloud.google.com/cloud-build/docs/api/reference/rest/)
- Enable the [Cloud Filestore API](https://cloud.google.com/filestore/docs/reference/rest/)
- Make sure that your account has the Project Editor access level, or had been granted the [roles/deploymentmanager.editor](https://cloud.google.com/deployment-manager/docs/access-control#predefined_roles) IAM role
- Make sure that the [Google APIs service account](https://cloud.google.com/deployment-manager/docs/access-control#access_control_for_deployment_manager) has **default** permissions, or had been explicitly granted the [roles/file.editor](https://cloud.google.com/functions/docs/reference/iam/roles#standard-roles) IAM role

## Deployment

### Resources

- [gcp-types/file-v1beta1:instances](https://cloud.google.com/filestore/docs/reference/rest/v1beta1/projects.locations.instances/create)

### Properties

See the `properties` section in the schema file(s):
- [Cloud Filestore](cloud_filestore.py.schema)

### Usage

1. Clone the [Deployment Manager Samples repository](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit):

```shell
    git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit
```

2. Go to the [dm](../../) directory:

```shell
    cd cloud-foundation-toolkit/dm
```

3. Copy the example DM config to be used as a model for the deployment; in this case, [examples/cloud\_filestore.yaml](examples/cloud_filestore.yaml):

```shell
    cp templates/cloud_filestore/examples/cloud_filestore.yaml my_cloud_filestore.yaml
```

4. Change the values in the config file to match your specific GCP setup (for properties, refer to the schema files listed above):

```shell
    vim my_cloud_filestore.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
    --config my_cloud_filestore.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Cloud Filestore](examples/cloud_filestore.yaml)
