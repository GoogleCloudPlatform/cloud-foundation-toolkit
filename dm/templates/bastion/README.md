# Bastion Host

> :warning: **NOTE**

Check out SSH via IAP as an alternative to Bastion Hosts:

- [Cloud IAP enables context-aware access to VMs via SSH and RDP without bastion hosts](https://cloud.google.com/blog/products/identity-security/cloud-iap-enables-context-aware-access-to-vms-via-ssh-and-rdp-without-bastion-hosts)
- [Using IAP for TCP forwarding](https://cloud.google.com/iap/docs/using-tcp-forwarding#tunneling_with_ssh)
> :warning: **NOTE** 

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

1. Clone the [Deployment Manager Samples repository](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit):

```shell
    git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit
```

2. Go to the [dm](../../) directory:

```shell
    cd dm
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
