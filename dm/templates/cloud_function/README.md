# Cloud Function

This template creates a Cloud Function.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Enable the [Cloud Build API](https://cloud.google.com/cloud-build/docs/api/reference/rest/)
- Enable the [Cloud Functions API](https://cloud.google.com/functions/docs/reference/rest/)
- Make sure that your account has the Project Editor access level, or had been granted the [roles/deploymentmanager.editor](https://cloud.google.com/deployment-manager/docs/access-control#predefined_roles) IAM role
- Make sure that the [Google APIs service account](https://cloud.google.com/deployment-manager/docs/access-control#access_control_for_deployment_manager) has **default** permissions, or had been explicitly granted the [roles/cloudfunctions.developer](https://cloud.google.com/functions/docs/reference/iam/roles#standard-roles) IAM role
- Make sure that the [Cloud Functions service account](https://cloud.google.com/functions/docs/concepts/iam#cloud_functions_service_account)
has **default** permissions, or had been granted the [CloudFunctions.ServiceAgent](https://cloud.google.com/functions/docs/concepts/iam#cloud_functions_service_account) IAM role

## Deployment

### Resources

- [cloudfunctions.v1beta2.function](https://cloud.google.com/functions/docs/reference/rest/v1beta2/projects.locations.functions)
- [storage.v1.bucket](https://cloud.google.com/storage/docs/json_api/v1/buckets)

### Properties

See the `properties` section in the schema file(s):
- [Cloud Function](cloud_function.py.schema)

### Usage

1. Clone the [Deployment Manager Samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples):

```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [community/cloud-foundation](../../) directory:

```shell
    cd community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment; in this case, [examples/cloud\_function.yaml](examples/cloud_function.yaml):

```shell
    cp templates/cloud_function/examples/cloud_function.yaml my_cloud_function.yaml
```

4. Change the values in the config file to match your specific GCP setup (for properties, refer to the schema files listed above):

```shell
    vim my_cloud_function.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
    --config my_cloud_function.yaml
```

`Note:` The Cloud Function HTTP trigger has no built-in authentication. Any user who has the link to the HTTP trigger can call it.

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

`Note:` To upload local source code to a GS bucket, the corresponding feature must be enabled first. This can be achieved by importing `upload.py` into your config file:

```yaml
    - path: templates/cloud_function/upload.py
      name: upload.py
```

`Note:` For Cloud Functions created from the local source code, deployment deletion will not delete the bucket to which that source code was uploaded during the build. Also, it will not clean up the [Cloud Build](https://cloud.google.com/cloud-build/) history.

## Examples

- [Cloud Function](examples/cloud_function.yaml)
- [Cloud Function with local source code upload](examples/cloud_function_upload.yaml)
