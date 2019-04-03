# Google Kubernetes Engine (GKE)

This template creates a Google Kubernetes Engine cluster.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Create a [network and subnetwork](../network/README.md)
- Grant the [container.admin](https://cloud.google.com/kubernetes-engine/docs/how-to/iam) IAM role to the Deployment Manager service account

## Deployment

### Resources

- [container-v1beta1:projects.locations.clusters](https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1beta1/projects.locations.clusters)
- [container.v1.cluster](https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1/projects.zones.clusters)

### Properties

See the `properties` section in the schema file(s):

- [GKE cluster schema](gke.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples)

    ```shell
        git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
    ```

2. Go to the [community/cloud-foundation](../../) directory

    ```shell
        cd community/cloud-foundation
    ```

3. Copy the example DM config to be used as a model for the deployment, in this case [examples/gke.yaml](examples/gke.yaml)

    ```shell
        cp templates/gke/examples/gke_zonal.yaml my_gke_zonal.yaml
    ```

4. Change the values in the config file to match your specific GCP setup.
   Refer to the properties in the schema files described above.

    ```shell
        vim my_gke_zonal.yaml  # <== change values to match your GCP setup
    ```

5. Create your deployment as described below, replacing <YOUR_DEPLOYMENT_NAME>
   with your with your own deployment name

    ```shell
        gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
            --config my_gke_zonal.yaml
    ```

6. In case you need to delete your deployment:

    ```shell
        gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
    ```

## Examples

- [GKE Zonal Cluster](examples/gke_zonal.yaml)
- [GKE Regional Cluster](examples/gke_regional.yaml)
- [GKE Private Regional Cluster](examples/gke_regional_private.yaml)
