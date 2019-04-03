# Google Cloud Key Management Service (KMS)

This template creates a Google Cloud KMS KeyRing and Keys.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [cloudkms.admin](https://cloud.google.com/kms/docs/iam) IAM role to
  the Deployment Manager service account

## Deployment

### Resources

- [gcp-types/cloudkms-v1](https://cloud.google.com/kms/docs/reference/rest/)
- [KMS Object heirarchy](https://cloud.google.com/kms/docs/object-hierarchy)
- [KMS Key Version States](https://cloud.google.com/kms/docs/key-states)
- [KMS Object Lifetime](https://cloud.google.com/kms/docs/object-hierarchy#lifetime)

### Properties

See the `properties` section in the schema file(s):

- [Cloud KMS](kms.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples)

   ```shell
   git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
   ```

2. Go to the [community/cloud-foundation](../../) directory

   ```shell
   cd community/cloud-foundation
   ```

3. Copy the example DM config to be used as a model for the deployment,
   in this case [examples/kms.yaml](examples/kms.yaml)

   ```shell
   cp templates/kms/examples/kms.yaml my_kms.yaml
   ```

4. Change the values in the config file to match your specific GCP setup.
   Refer to the properties in the schema files described above.

   ```shell
   vim my_kms.yaml  # <== Replace all <FIXME:..> placeholders in this file
   ```

5. Create your deployment as described below, replacing <YOUR_DEPLOYMENT_NAME>
   with your with your own deployment name

   ```shell
   gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_kms.yaml
   ```

> **Note**: Once created, this deployment cannot be deleted.
> Refer to `KMS Object Lifetime` in [Resources](#Resources) section

## Examples

- [KMS KeyRing with Encryption Keys](examples/kms.yaml)
- [KMS KeyRing with Signing Keys](examples/kms_signkey.yaml)
