# Cloud Build

This template creates a Google Cloud Build.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Enable the following in the [APIs & Services](https://console.cloud.google.com/apis/dashboard) section of the Google Cloud console:
  - [Cloud Build API](https://console.cloud.google.com/apis/library/cloudbuild.googleapis.com)
  - [Cloud Source Repositories API](https://console.cloud.google.com/apis/library/sourcerepo.googleapis.com)
  - [Container Registry API](https://console.cloud.google.com/apis/library/containerregistry.googleapis.com)
- Grant to the Cloud Build service account the IAM roles necessary for the steps in your build

## Deployment

### Resources

- [projects.builds](https://cloud.google.com/cloud-build/docs/api/reference/rest/v1/projects.builds)
- [projects.triggers](https://cloud.google.com/cloud-build/docs/api/reference/rest/v1/projects.triggers)
- [cloud builders](https://cloud.google.com/cloud-build/docs/cloud-builders)
- [cloud builders community](https://github.com/GoogleCloudPlatform/cloud-builders-community)

### Properties

See the `properties` section in the schema file(s):

- [CloudBuild build schema](cloudbuild.py.schema)
- [CloudBuild trigger schema](trigger.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples)

    ```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
    ```

2. Go to the [community/cloud-foundation](../../) directory

    ```shell
    cd community/cloud-foundation
    ```

3. Copy the example DM config to be used as a model for the deployment, in this case [examples/cloudbuild.yaml](examples/cloudbuild.yaml)

    ```shell
    cp templates/cloudbuild/examples/cloudbuild.yaml my_cloudbuild.yaml
    ```

4. Change the values in the config file to match your specific GCP setup.
   Refer to the properties in the schema files described above.

    ```shell
    vim my_cloudbuild.yaml  # <== change values to match your GCP setup
    ```

5. Create your deployment as described below, replacing `<YOUR_DEPLOYMENT_NAME>`
   with your with your own deployment name

    ```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_cloudbuild.yaml
    ```

6. In case you need to delete your deployment:

    ```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
    ```

## Examples

- [Cloud Build](examples/cloudbuild.yaml)
- [Cloud Build with StorageSource](examples/cloudbuild_storagesource.yaml)
- [Cloud Build with RepoSource](examples/cloudbuild_reposource.yaml)
- [Cloud Build Trigger](examples/cloudbuild_trigger.yaml)
