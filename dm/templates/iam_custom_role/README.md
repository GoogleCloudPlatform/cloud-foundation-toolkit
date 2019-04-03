# Custom IAM Role

This template creates a custom IAM role for an organization or a project.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [iam.roleAdmin, iam.organizationRoleAdmin or owner](https://cloud.google.com/iam/docs/understanding-custom-roles#required_permissions_and_roles) IAM role to the project service account

## Deployment

### Resources

- [Creating custom IAM roles](https://cloud.google.com/iam/docs/creating-custom-roles)
- [gcp-types/iam-v1:organizations.roles](https://cloud.google.com/iam/reference/rest/v1/organizations.roles/create)
- [gcp-types/iam-v1:projects.roles](https://cloud.google.com/iam/reference/rest/v1/projects.roles/create)

### Properties

See `properties` section in the schema file(s):

-  [Organization](organization_custom_role.py.schema)
-  [Project](project_custom_role.py.schema)


### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples):

```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [community/cloud-foundation](../../) directory:

```shell
    cd community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment; in this case, [examples/iam\_custom\_role.yaml](examples/iam_custom_role.yaml):

```shell
    cp templates/iam_custom_role/examples/iam_custom_role.yaml my_iamcustomrole.yaml
```

4. Change the values in the config file to match your specific GCP setup (for properties, refer to the schema files listed above):

```shell
    vim my_iamcustomrole.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
    --config my_iamcustomrole.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Custom IAM role](examples/iam_custom_role.yaml)
