# Autoscaler

This template creates an autoscaler.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [compute.admin](https://cloud.google.com/compute/docs/access/iam)
IAM role to the Deployment Manager service account

## Deployment

### Resources

- [compute.v1.autoscaler](https://cloud.google.com/compute/docs/reference/latest/autoscalers)
- [compute.v1.regionalAutoscaler](https://cloud.google.com/compute/docs/reference/latest/regionAutoscalers)

### Properties

See the `properties` section in the schema file(s):

- [Autoscaler](autoscaler.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples)

```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [community/cloud-foundation](../../) directory

```shell
    cd community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment, in this
   case [examples/autoscaler\_zonal.yaml](examples/autoscaler_zonal.yaml)

```shell
    cp templates/autoscaler/examples/autoscaler_zonal.yaml my_autoscaler.yaml
```

4. Change the values in the config file to match your specific GCP setup.
   Refer to the properties in the schema files described above.

```shell
    vim my_autoscaler.yaml  # <== change values to match your GCP setup
```

5. Create your deployment as described below, replacing <YOUR_DEPLOYMENT_NAME>
   with your with your own deployment name

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_autoscaler.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Zonal autoscaler](examples/autoscaler_zonal.yaml)
- [Regional autoscaler](examples/autoscaler_regional.yaml)
