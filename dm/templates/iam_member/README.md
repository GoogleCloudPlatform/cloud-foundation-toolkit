# IAM Member

This template grants IAM roles for a projects, folders and organizations.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)

### Grant the appropriate IAM permissions depending on your usecase
Grant the [owner](https://cloud.google.com/iam/docs/understanding-roles) IAM role on the project to the *DM Service Account* to grant roles within the project. This allows DM to set IAM on the Project or on the resource level.

For more restrictive permissions grant the appropriate resource level admin permission:

- Grant the [resourcemanager.projectIamAdmin](https://cloud.google.com/iam/docs/understanding-roles) IAM role on the project to the *DM Service Account* to grant roles within the project
- Grant the [roles/resourcemanager.folderIamAdmin](https://cloud.google.com/iam/docs/understanding-roles) IAM role on the folder to the *DM Service Account* to grant roles within the folder
- Grant the [roles/iam.securityAdmin](https://cloud.google.com/iam/docs/understanding-roles) IAM role on the organization to the *DM Service Account* to grant roles within the organization and all nested resources
- Etc.

## Development

### Resources

Resources are created based on the input properties:
- [cloudresourcemanager-v1:virtual.projects.iamMemberBinding](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/google/resource-snippets/cloudresourcemanager-v1/policies.jinja)
    - This virtual endpoint implements projects.getIamPolicy and projects.setIamPolicy internally with proper concurancy handling.
- [cloudresourcemanager-v2:virtual.folders.iamMemberBinding](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/google/resource-snippets/cloudresourcemanager-v2/policies.jinja)
- [cloudresourcemanager-v1:virtual.organizations.iamMemberBinding](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/google/resource-snippets/cloudresourcemanager-v1/policies.jinja)
- storage-v1:virtual.buckets.iamMemberBinding
- cloudfunctions-v1:virtual.projects.locations.functions.iamMemberBinding

### Properties

See `properties` section in the schema file(s):

-  [IAM Member](iam_member.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit):

```shell
    git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit
```

2. Go to the [dm](../../) directory:

```shell
    cd dm
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
