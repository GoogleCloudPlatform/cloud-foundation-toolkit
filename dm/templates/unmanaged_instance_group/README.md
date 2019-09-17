# Unmanaged Instance Group

This template creates a unmanaged instance group.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, setup billing, enable requisite APIs](../project/README.md)
- Grant the [compute.admin](https://cloud.google.com/compute/docs/access/iam) IAM role to the [Deployment Manager service account](https://cloud.google.com/deployment-manager/docs/access-control#access_control_for_deployment_manager)

## Deployment

### Resources

- [compute.v1.instanceGroups](https://cloud.google.com/compute/docs/reference/latest/instanceGroups)

### Properties

See the `properties` section in the schema file(s):

- [Unmanaged Instance Group](unmanaged_instance_group.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit)

```shell
    git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit
```

2. Go to the [dm](../../) directory

```shell
    cd dm
```

3. Copy the example DM config to be used as a model for the deployment, in this
   case [examples/unmanaged\_instance\_group\_add\_instance.yaml](examples/unmanaged_instance_group_add_instance.yaml)

```shell
    cp templates/unmanaged_instance_group/examples/unmanaged_instance_group_add_instance.yaml \
        my_unmanaged_instance_group_add_instance.yaml
```

4. Change the values in the config file to match your specific GCP setup.
   Refer to the properties in the schema files described above.

```shell
    vim my_unmanaged_instance_group_add_instance.yaml # <== change values to match your GCP setup
```

5. Create your deployment as described below, replacing
   \<YOUR\_DEPLOYMENT\_NAME\> with your with your own deployment name

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_unmanaged_instance_group_add_instance.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Unmanaged Instance Group](examples/unmanaged_instance_group.yaml)
