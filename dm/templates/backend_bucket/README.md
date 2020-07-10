# Backend Bucket

This template creates a backend bucket.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [compute.admin](https://cloud.google.com/compute/docs/access/iam)
IAM role to the Deployment Manager service account

## Deployment

### Resources

- [compute.v1.backendBucket](https://cloud.google.com/compute/docs/reference/rest/v1/backendBuckets)

### Properties

See the `properties` section in the schema file(s):
- [Backend Bucket](backend_bucket.py.schema)

### Usage

1. Clone the [Deployment Manager Samples repository](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit):

```shell
    git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit
```

2. Go to the [dm](../../) directory:

```shell
    cd dm
```

3. Copy the example DM config to be used as a model for the deployment; in this
case, [examples/backend\_bucket.yaml](examples/backend_bucket.yaml):

```shell
    cp templates/backend_bucket/examples/backend_bucket.yaml \
       my_backend_bucket.yaml
```

4. Change the values in the config file to match your specific GCP setup (for
properties, refer to the schema files listed above):

```shell
    vim my_backend_bucket.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant
deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_backend_bucket.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Backend Bucket](examples/backend_bucket.yaml)
