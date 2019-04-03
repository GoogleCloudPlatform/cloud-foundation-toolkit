# Interconnect Attachment

This template creates an Interconnect Attachment.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Create a [network](../network/README.md)
- Create a [interconnect] (../interconnect/README.md)
- Create a [cloud_router] (../cloud_router/README.md)
- Grant the [compute.networkAdmin](https://cloud.google.com/compute/docs/access/iam) i
  IAM role to the project service account

## Deployment

### Resources

- [compute.v1.interconnectAttachments](https://cloud.google.com/compute/docs/reference/rest/v1/interconnectAttachments)

### Properties

See the `properties` section in the schema file(s):

- [Interconnect Attachment](interconnect_attachment.py.schema)

### Usage

1. Clone the [Deployment Manager Samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples):

```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [community/cloud-foundation](../../) directory:

```shell
    cd community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment; in this case, [examples/interconnect_attachment.yaml](examples/interconnect_attachment.yaml):

```shell
    cp templates/interconnect_attachment/examples/interconnect_attachment.yaml my_interconnect_attachment.yaml
```

4. Change the values in the config file to match your specific GCP setup (for 
   properties, refer to the schema files listed above):

```shell
    vim my_interconnect_attachment.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace \<YOUR\_DEPLOYMENT\_NAME\> with the relevant 
   deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
    --config my_interconnect_attachment.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Dedicated Interconnect Attachment](examples/interconnect_attachment_dedicated.yaml)
- [Partner Interconnect Attachment](examples/interconnect_attachment_partner.yaml)
