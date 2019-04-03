# Managed Instance Group

This template creates a managed instance group.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, setup billing, enable requisite APIs](../project/README.md)
- Grant the [compute.admin](https://cloud.google.com/compute/docs/access/iam) IAM role to the [Deployment Manager service account](https://cloud.google.com/deployment-manager/docs/access-control#access_control_for_deployment_manager)

## Deployment

### Resources

- [compute.v1.instance](https://cloud.google.com/compute/docs/reference/latest/instances)
- [compute.v1.autoscaler](https://cloud.google.com/compute/docs/reference/latest/autoscalers)
- [compute.v1.regionalAutoscaler](https://cloud.google.com/compute/docs/reference/latest/regionAutoscalers)
- [compute.v1.instanceTemplate](https://cloud.google.com/compute/docs/reference/latest/instanceTemplates)
- [compute.v1.instanceGroupManager](https://cloud.google.com/compute/docs/reference/latest/instanceGroupManagers)
- [compute.v1.regionalInstanceGroupManager](https://cloud.google.com/compute/docs/reference/latest/regionInstanceGroupManagers)

### Properties

See the `properties` section in the schema file(s):

- [Managed Instance Group](managed_instance_group.py.schema)

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
   case [examples/managed\_instance\_group.yaml](examples/managed_instance_group.yaml)

```shell
    cp templates/managed_instance_group/examples/managed_instance_group.yaml \
        my_managed_instance_group.yaml
```

4. Change the values in the config file to match your specific GCP setup.
   Refer to the properties in the schema files described above.

```shell
    vim my_managed_instance_group.yaml # <== change values to match your GCP setup
```

5. Create your deployment as described below, replacing
   \<YOUR\_DEPLOYMENT\_NAME\> with your with your own deployment name

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_managed_instance_group.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Managed Instance Group](examples/managed_instance_group.yaml)
- [Managed Instance Group with Health Check](examples/managed_instance_group_healthcheck.yaml)
