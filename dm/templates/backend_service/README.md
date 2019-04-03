# Backend Service

This template creates a backend service.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [compute.admin](https://cloud.google.com/compute/docs/access/iam)
IAM role to the Deployment Manager service account

## Deployment

### Resources

- [compute.v1.backendService](https://cloud.google.com/compute/docs/reference/rest/v1/backendServices)
- [compute.v1.regionalBackendService](https://cloud.google.com/compute/docs/reference/latest/regionBackendServices)

### Properties

See the `properties` section in the schema file(s):
- [Backend Service](backend_service.py.schema)

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
case, [examples/backend\_service\_regional.yaml](examples/backend_service_regional.yaml):

```shell
    cp templates/backend_service/examples/backend_service_regional.yaml \
       my_backend_service.yaml
```

4. Change the values in the config file to match your specific GCP setup (for
properties, refer to the schema files listed above):

```shell
    vim my_backend_service.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant
deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_backend_service.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Regional Backend Service](examples/backend_service_regional.yaml)
- [Global Backend Service](examples/backend_service_global.yaml)
