# Bastion Host

This template creates a Bastion host. Once it had been deployed, one can use
`gcloud compute ssh <BASTION_HOST_NAME> --zone <ZONE>` to connect to
the Bastion host, and then use
`gcloud compute ssh <TARGET_HOST_NAME> --zone <ZONE> --internal-ip` to SSH to
another host, within the same network, that has no external IP assigned.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [compute.computeAdmin](https://cloud.google.com/compute/docs/access/iam)
  IAM role to the Deployment Manager service account

## Deployment

### Resources

- [compute.v1.instance](https://cloud.google.com/compute/docs/reference/rest/v1/instances)
- [compute.v1.firewall](https://cloud.google.com/compute/docs/reference/rest/v1/firewalls)

### Properties

See the `properties` section in the schema file(s):

- [Bastion Host](bastion.py.schema)

### Usage

1. Clone the [Deployment Manager Samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples):

```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [community/cloud-foundation](../../) directory:

```shell
    cd community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment; in this
   case, [examples/bastion.yaml](examples/bastion.yaml):

```shell
    cp templates/bastion/examples/bastion.yaml \
       my_bastion.yaml
```

4. Change the values in the config file to match your specific GCP setup (for
   properties, refer to the schema files listed above):

```shell
    vim my_bastion.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace \<YOUR\_DEPLOYMENT\_NAME\> with the relevant
   deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_bastion.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Bastion](examples/bastion.yaml)
