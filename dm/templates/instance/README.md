# Instance

This template creates a Compute Instance.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [compute.admin](https://cloud.google.com/compute/docs/access/iam) IAM role to the [Deployment Manager service account](https://cloud.google.com/deployment-manager/docs/access-control#access_control_for_deployment_manager)

## Deployment

### Resources

- [compute.v1.instance](https://cloud.google.com/compute/docs/reference/rest/v1/instances)

### Properties

See the `properties` section in the schema file(s):

- [Instance](instance.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples):

```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [community/cloud-foundation](../../) directory:

```shell
    cd community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment; in this case, [examples/instance.yaml](examples/instance.yaml):

```shell
    cp templates/instance/examples/instance.yaml my_instance.yaml
```

4. Change the values in the config file to match your specific GCP setup (for properties, refer to the schema files listed above):

```shell
    vim my_instance.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
    --config my_instance.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Compute Instance](examples/instance.yaml)
