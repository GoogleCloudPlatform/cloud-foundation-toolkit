# Google Cloud Storage Bucket

This template creates a Google Cloud Storage bucket.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [storage.admin](https://cloud.google.com/storage/docs/access-control/iam-roles) IAM role to the Deployment Manager service account

## Deployment

### Resources

- [storage.v1.bucket](https://cloud.google.com/storage/docs/creating-buckets)

### Properties

See the `properties` section in the schema file(s):

- [gcs_bucket](gcs_bucket.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples)

   ```shell
   git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
   ```

2. Go to the [community/cloud-foundation](../../) directory

   ```shell
   cd community/cloud-foundation
   ```

3. Copy the example DM config to be used as a model for the deployment, in this case [examples/gcs\_bucket.yaml](examples/gcs_bucket.yaml)

   ```shell
   cp templates/gcs_bucket/examples/gcs_bucket.yaml my_gcs_bucket.yaml
   ```

4. Change the values in the config file to match your specific GCP setup.
   Refer to the properties in the schema files described above.

   ```shell
   vim my_gcs_bucket.yaml  # <== Replace the <FIXME:..> placeholders in the file
   ```

5. Create your deployment as described below, replacing <YOUR_DEPLOYMENT_NAME>
   with your with your own deployment name

   ```shell
   gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
       --config my_gcs_bucket.yaml
   ```

6. In case you need to delete your deployment:

   ```shell
   gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
   ```

## Examples

- [Storage Bucket](examples/gcs_bucket.yaml)
- [Storage Bucket with LifeCycle Enabled](examples/gcs_bucket_lifecycle.yaml)
- [Storage Bucket with IAM Bindings](examples/gcs_bucket_iam_bindings.yaml)
