# Organization Policy

This template creates an organization policy.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [resourcemanager.organizationAdmin](https://cloud.google.com/resource-manager/docs/access-control-org) IAM role to the project service account


## Deployment

### Resources

- [cloudresourcemanager.v1.project](https://cloud.google.com/resource-manager/reference/rest/v1/projects/setOrgPolicy)


### Properties

See the `properties` section in the schema file(s):

-  [Org Policy](org_policy.py.schema)


### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples):

```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [community/cloud-foundation](../../) directory:

```shell
    cd community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment; in this case, [examples/org_policy.yaml](examples/org_policy.yaml):

```shell
    cp templates/org_policy/examples/org_policy.yaml my_org_policy.yaml
```

4. Change the values in the config file to match your specific GCP setup (for properties, refer to the schema files listed above):

```shell
    vim my_org_policy.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_org_policy.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Org policy](examples/org_policy.yaml)
