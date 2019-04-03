# VPN

This template creates a VPN.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Create a [network](../network/README.md)
- Grant the [compute.networkAdmin or compute.admin](https://cloud.google.com/compute/docs/access/iam) IAM role to the project service account

## Deployment

### Resources

- [compute.v1.targetVpnGateway](https://cloud.google.com/compute/docs/reference/latest/targetVpnGateways)
- [compute.v1.address](https://cloud.google.com/compute/docs/reference/rest/v1/addresses)
- [compute.v1.forwardingRule](https://cloud.google.com/compute/docs/reference/latest/forwardingRules)
- [compute.v1.vpnTunnel](https://cloud.google.com/compute/docs/reference/latest/vpnTunnels)
- [gcp-types/compute-v1:compute.routers.patch](https://www.googleapis.com/discovery/v1/apis/compute/v1/rest)

### Properties

See `properties` section in the schema file(s):

-  [VPN](../vpn/vpn.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples):

```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [community/cloud-foundation](../../) directory:

```shell
    cd community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment; in this case, [examples/vpn.yaml](examples/vpn.yaml):

```shell
    cp templates/vpn/examples/vpn.yaml my_vpn.yaml
```

4. Change the values in the config file to match your specific GCP setup (for properties, refer to the schema files listed above):

```shell
    vim my_vpn.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
    --config my_vpn.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [VPN](examples/vpn.yaml)
