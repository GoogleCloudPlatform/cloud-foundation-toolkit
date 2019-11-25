# Cloud Run

This template creates a Cloud Run configuration.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
  IAM role to the Deployment Manager service account

## Deployment

### Properties

See the `properties` section in the schema file(s):

- [Cloud Run](cloud_run.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit):

    ```shell
    git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit
    ```

2. Go to the [dm](../../) directory:

    ```shell
    cd dm
    ```

3. Copy the example DM config to be used as a model for the deployment; in this
   case, [examples/cloud\_run.yaml](examples/cloud_run.yaml):

    ```shell
    cp templates/cloud_run/examples/cloud_run.yaml my_cloud_run.yaml
    ```

4. Change the values in the config file to match your specific GCP setup. 
   Generate access token with scope https://www.googleapis.com/auth/cloud-platform
   using https://developers.google.com/oauthplayground (for
   properties, refer to the schema files listed above):

    ```shell
    vim my_cloud_run.yaml  # <== change values to match your GCP setup
    ```

5. 

5. Create custom type and your deployment (replace \<YOUR\_DEPLOYMENT\_NAME\> with the relevant
   deployment name):

    ```shell
    gcloud config set project <YOUR_PROJECT_NAME>
    gcloud beta deployment-manager type-providers create cloud-run-custom-type --descriptor-url='https://run.googleapis.com/$discovery/rest?version=v1'
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_cloud_run.yaml
    ```

   To deploy with CFT:

    ```shell
    cft apply my_cloud_run.yaml
    ```

6. In case you need to delete your deployment:

    ```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
    gcloud beta deployment-manager type-providers delete cloud-run-custom-type
    ```

   To delete deployment with CFT:

    ```shell
    cft delete my_cloud_run.yaml
    ```

## Examples

- [Cloud Run](examples/cloud_run.yaml)

