# VPC Network Peering

This template creates peering between VPC networks.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Create a [network](../network/README.md)
- Grant the [compute.networkAdmin](https://cloud.google.com/iam/docs/understanding-roles#compute_engine_roles) IAM role to the Deployment Manager service account

## Deployment

### Resources

- [compute.v1.network.addPeering](https://cloud.google.com/compute/docs/reference/rest/v1/networks/addPeering)
- [compute.v1.network.removePeering](https://cloud.google.com/compute/docs/reference/rest/v1/networks/removePeering)

### Properties

See the `properties` section in the schema file(s):

- [VPC network peering](network_peering.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples)

    ```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
    ```

2. Go to the [community/cloud-foundation](../../) directory

    ```shell
    cd community/cloud-foundation
    ```

3. Copy the example DM config to be used as a model for the deployment, in this case [examples/network_peering.yaml](examples/network_peering.yaml)

    ```shell
    cp templates/network_peering/examples/network_peering.yaml my_network_peering.yaml
    ```

4. Change the values in the config file to match your specific GCP setup.
   Refer to the properties in the schema files described above.

    ```shell
    vim my_network_peering.yaml  # <== change values to match your GCP setup
    ```

5. Create your deployment as described below, replacing <YOUR_DEPLOYMENT_NAME>
   with your with your own deployment name

    ```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_network_peering.yaml
    ```

6. In case you need to delete your deployment:

    ```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
    ```

**IMPORTANT**: When deleting an _ACTIVE_ peering, you may recieve an error.  
Refer to [peering limitations](https://cloud.google.com/vpc/docs/using-vpc-peering#number_of_peerings_limit) and  [troubleshooting](https://cloud.google.com/vpc/docs/using-vpc-peering#troubleshooting) for details.

## Examples

- [VPC Network Peering](examples/network_peering.yaml)
