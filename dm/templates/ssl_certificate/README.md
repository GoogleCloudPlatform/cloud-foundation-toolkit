# SSL Certificate

This template creates an SSL certificate.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [roles/compute.securityAdmin](https://cloud.google.com/compute/docs/access/iam),
  or [compute.loadBalancerAdmin](https://cloud.google.com/compute/docs/access/iam)
  IAM role to the Deployment Manager service account

## Deployment

### Resources

- [gcp-types/compute-v1:sslCertificates](https://cloud.google.com/compute/docs/reference/rest/v1/sslCertificates)
- [gcp-types/compute-beta:sslCertificates](https://cloud.google.com/compute/docs/reference/rest/beta/sslCertificates)
(If `betaFeaturesEnabled` flag turned on.) 

### Properties

See the `properties` section in the schema file(s):

- [SSL Certificate](ssl_certificate.py.schema)

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
   case, [examples/ssl\_certificate.yaml](examples/ssl_certificate.yaml):

```shell
    cp templates/ssl_certificate/examples/ssl_certificate.yaml \
       my_ssl_certificate.yaml
```

4. Change the values in the config file to match your specific GCP setup (for
   properties, refer to the schema files listed above):

```shell
    vim my_ssl_certificate.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace \<YOUR\_DEPLOYMENT\_NAME\> with the relevant
   deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_ssl_certificate.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [SSL Certificate](examples/ssl_certificate.yaml)
