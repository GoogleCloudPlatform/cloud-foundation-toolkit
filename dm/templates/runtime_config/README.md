# Google Cloud Runtime Configurator

This template creates a Runtime Configurator with the associated resources.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Install gcloud **beta** components:

  ```shell
  gcloud components update
  gcloud components install beta
  ```

- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Enable the [Cloud Runtime Configurator API](https://console.developers.google.com/apis/api/runtimeconfig.googleapis.com)
- Grant the [Cloud RuntimeConfig Admin](https://cloud.google.com/deployment-manager/runtime-configurator/access-control)
  IAM role to the Deployment Manager service account

## Deployment

### Resources

- [v1beta1.projects.configs](https://cloud.google.com/deployment-manager/runtime-configurator/create-and-delete-runtimeconfig-resources)

### Properties

See the `properties` section in the schema file(s):

- [Runtime Config Schema](runtime_config.py.schema)
- [Variable Schema](variable.py.schema)
- [Waiter Schema](waiter.py.schema)

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
   in this case [examples/runtime\_config.yaml](examples/runtime_config.yaml)

   ```shell
   cp templates/runtime_config/examples/runtime_config.yaml my_runtime_config.yaml
   ```

4. Change the values in the config file to match your specific GCP setup.
   Refer to the properties in the schema files described above.

   ```shell
   vim my_runtime_config.yaml  # <== change values to match your GCP setup
   ```

5. Create your deployment as described below, replacing `<YOUR_DEPLOYMENT_NAME>`
   with your with your own deployment name

   ```shell
   gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
       --config my_runtime_config.yaml
   ```

6. In case you need to delete your deployment:

   ```shell
   gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
   ```

## Examples

- [Cloud Runtime Configurator with Variables and Waiters](examples/runtime_config.yaml)
