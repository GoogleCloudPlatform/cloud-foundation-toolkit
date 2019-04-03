# Cloud DNS Managed Zone

This template creates a managed zone in the Cloud DNS (Domain Name System).

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [dns.admin](https://cloud.google.com/dns/access-control) IAM role to the Deployment Manager service account

## Deployment

### Resources

- [dns.v1.managedZone](https://cloud.google.com/dns/docs/)

### Properties

See the `properties` section in the schema file(s):
- [Cloud DNS Managed Zone](dns_managed_zone.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples):

```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [community/cloud-foundation](../../) directory:

```shell
    cd community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment; in this case, [examples/dns_managed_zone.yaml](examples/dns_managed_zone.yaml):

```shell
    cp templates/dns_managed_zone/examples/dns_managed_zone.yaml my_dns_managed_zone.yaml
```

4. Change the values in the config file to match your specific GCP setup (for properties, refer to the schema files listed above):

```shell
    vim my_dns_managed_zone.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
    --config my_dns_managed_zone.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Cloud DNS Managed Zone](examples/dns_managed_zone.yaml)
