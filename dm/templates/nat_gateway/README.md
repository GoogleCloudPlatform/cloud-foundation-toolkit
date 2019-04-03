# HA NAT Gateway

This template creates a High Availability NAT Gateway based on the number of
regions specified. Each gateway is a managed instance group of one with
auto-healing through healthchecks. The only firewall rule created is for the
instance healthcheck. Any additional traffic you wish to go through the gateway
will require additional firewall rules (for example, TCP/UDP/ICMP, etc.).

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [compute.admin](https://cloud.google.com/compute/docs/access/iam)
  IAM role to the [Deployment Manager service account](https://cloud.google.com/deployment-manager/docs/access-control#access_control_for_deployment_manager)
- Grant the [compute.networkAdmin](https://cloud.google.com/compute/docs/access/iam)
  IAM role to the [Deployment Manager service account](https://cloud.google.com/deployment-manager/docs/access-control#access_control_for_deployment_manager)
- NOTE: The NAT Gateway integration tests will need additional IAM permissions. The tests will SSH into test instances to verify the NAT functionality. Please refer to [Managing Instance Access Using OS Login](https://cloud.google.com/compute/docs/instances/managing-instance-access#enable_oslogin) and [Connecting through a bastion host](https://cloud.google.com/compute/docs/instances/connecting-advanced#bastion_host) page for additional information.

## Deployment

### Resources

- [compute.v1.addresses](https://cloud.google.com/compute/docs/reference/rest/v1/addresses)
- [compute.v1.instanceTemplate](https://cloud.google.com/compute/docs/reference/latest/instanceTemplates)
- [compute.v1.instanceGroupManagers](https://cloud.google.com/compute/docs/reference/rest/v1/instanceGroupManagers)
- [compute.v1.firewalls](https://cloud.google.com/compute/docs/reference/rest/v1/firewalls)
- [compute.v1.routes](https://cloud.google.com/compute/docs/reference/rest/v1/routes)
- [compute.v1.healthChecks](https://cloud.google.com/compute/docs/reference/rest/v1/healthChecks)

### Properties

See the `properties` section in the schema file(s):

- [NAT Gateway](nat_gateway.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples)

```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [community/cloud-foundation](../../) directory

```shell
    cd community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment, in this
   case [examples/nat\_gateway.yaml](examples/nat_gateway.yaml)

```shell
    cp templates/nat_gateway/examples/nat_gateway.yaml \
        my_nat_gateway.yaml
```

4. Change the values in the config file to match your specific GCP setup.
   Refer to the properties in the schema files described above.

```shell
    vim my_nat_gateway.yaml  # <== change values to match your GCP setup
```

5. Create your deployment as described below, replacing
   \<YOUR\_DEPLOYMENT\_NAME\> with your with your own deployment name

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_nat_gateway.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [NAT Gateway](examples/nat_gateway.yaml)
