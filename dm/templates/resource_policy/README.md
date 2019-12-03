# Resource Policy

This template creates a resource policy.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Enable the [Compute Engine API](https://cloud.google.com/compute/docs/reference/rest/v1/)
- Make sure that the [Google APIs service account](https://cloud.google.com/deployment-manager/docs/access-control#access_control_for_deployment_manager) has *compute.resourcePolicies.create* permissions

## Deployment

### Resources

- [gcp-types/compute-v1:resourcePolicies](https://cloud.google.com/compute/docs/reference/rest/v1/resourcePolicies/insert)

### Properties

See the `properties` section in the schema file(s):
- [Resource Policy](resource_policy.py.schema)

### Usage

1. Clone the [Deployment Manager Samples repository](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit):

```shell
    git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit
```

2. Go to the [dm](../../) directory:

```shell
    cd cloud-foundation-toolkit/dm
```

3. Copy the example DM config to be used as a model for the deployment; in this case, [examples/resource\_policy.yaml](examples/resource_policy.yaml):

```shell
    cp templates/resource_policy/examples/resource_policy.yaml my_resource_policy.yaml
```

4. Change the values in the config file to match your specific GCP setup (for properties, refer to the schema files listed above):

```shell
    vim my_resource_policy.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
    --config my_resource_policy.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Resource Policy](examples/resource_policy.yaml)
