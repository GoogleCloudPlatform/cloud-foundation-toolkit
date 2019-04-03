# Cloud Router

This template creates a Cloud Router.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Create a [network](../network/README.md)
- Grant the [compute.networkAdmin](https://cloud.google.com/compute/docs/access/iam) IAM role to the project service account

## Deployment

### Resources

- [compute.v1.router](https://cloud.google.com/compute/docs/reference/rest/v1/routers)

### Properties

See the `properties` section in the schema file(s):
- [Cloud Router](cloud_router.py.schema)

### Usage

1. Clone the [Deployment Manager Samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples):

```
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [deploymentmanager-samples/community/cloud-foundation](../../) directory:

```
    cd deploymentmanager-samples/community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment; in this case, [examples/cloud_router.yaml](examples/cloud_router.yaml):

```
    cp templates/cloud_router/examples/cloud_router.yaml my_cloud_router.yaml
```

4. Change the values in the config file to match your specific GCP setup (for properties, refer to the schema files listed above):

```
    vim my_cloud_router.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

```
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
    --config my_cloud_router.yaml
```

6. In case you need to delete your deployment:

```
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Cloud Router](examples/cloud_router.yaml)
