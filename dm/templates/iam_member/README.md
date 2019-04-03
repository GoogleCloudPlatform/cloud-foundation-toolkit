# IAM Member

This template grants IAM roles for a project.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [resourcemanager.projectIamAdmin or owner](https://cloud.google.com/iam/docs/understanding-roles) IAM role to the project service account

## Development

### Resources

- [virtual.projects.iamMemberBinding](https://github.com/GoogleCloudPlatform/deploymentmanager-samples/blob/master/google/resource-snippets/cloudresourcemanager-v1/policies.jinja)
    - This virtual endpoint implements projects.getIamPolicy and projects.setIamPolicy internally with proper concurancy handling.

### Properties

See `properties` section in the schema file(s):

-  [IAM Member](iam_member.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples):

```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [community/cloud-foundation](../../) directory:

```shell
    cd community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment; in this case, [examples/iam_member.yaml](examples/iam_member.yaml):

```shell
    cp templates/iam_member/examples/iam_member.yaml my_iammember.yaml
```

4. Change the values in the config file to match your specific GCP setup (for properties, refer to the schema files listed above):

```shell
    vim my_iammember.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
    --config my_iammember.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [IAM member](examples/iam_member.yaml)
