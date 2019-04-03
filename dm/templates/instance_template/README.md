# Instance Template

This template creates an instance template.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, setup billing, enable requisite APIs](../project/README.md)
- Grant the [compute.admin](https://cloud.google.com/compute/docs/access/iam) IAM
role to the [Deployment Manager service account](https://cloud.google.com/deployment-manager/docs/access-control#access_control_for_deployment_manager)

## Deployment

### Resources

- [compute.v1.instanceTemplate](https://cloud.google.com/compute/docs/reference/latest/instanceTemplates)

### Properties

See the `properties` section in the schema file(s):

- [Instance Template](instance_template.py.schema)

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
   case [examples/instance\_template.yaml](examples/instance_template.yaml)

```shell
    cp templates/instance_template/examples/instance_template.yaml \
        my_instance_template.yaml
```

4. Change the values in the config file to match your specific GCP setup.
   Refer to the properties in the schema files described above.

```shell
    vim my_instance_template.yaml  # <== change values to match your GCP setup
```

5. Create your deployment as described below, replacing <YOUR_DEPLOYMENT_NAME>
   with your with your own deployment name

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_instance_template.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Instance Template](examples/instance_template.yaml)
