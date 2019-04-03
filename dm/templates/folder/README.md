# Folder

This template creates a folder under an organization or under a parent folder.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [resourcemanager.folderAdmin or resourcemanager.folderCreator](https://cloud.google.com/resource-manager/docs/access-control-folders) IAM role to the project service account

## Deployment

### Resources

- [gcp-types/cloudresourcemanager-v2:folders](https://cloud.google.com/resource-manager/reference/rest/v2/folders/create)


### Properties

See `properties` section in the schema file(s):

-  [Folder](folder.py.schema)

### Usage


1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples):

```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [community/cloud-foundation](../../) directory:

```shell
    cd community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment; in this case, [examples/folder.yaml](examples/folder.yaml):

```shell
    cp templates/folder/examples/folder.yaml my_folder.yaml
```

4. Change the values in the config file to match your specific GCP setup (for properties, refer to the schema files listed above):

```shell
    vim my_folder.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
    --config my_folder.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Folder](examples/folder.yaml)
