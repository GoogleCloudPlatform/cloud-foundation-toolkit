# Internal Load Balancer

This template creates an internal load balancer that consists of a forwarding
rule and a regional backend service.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [compute.admin](https://cloud.google.com/compute/docs/access/iam)
  IAM role to the Deployment Manager service account

## Deployment

### Resources

- [compute.v1.forwardingRule](https://cloud.google.com/compute/docs/reference/latest/forwardingRules)
- [compute.v1.regionalBackendService](https://cloud.google.com/compute/docs/reference/latest/regionBackendServices)

### Properties

See the `properties` section in the schema file(s):

- [Internal Load Balancer](internal_load_balancer.py.schema)

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
   case, [examples/internl\_load\_balancer.yaml](examples/internal_load_balancer.yaml):

```shell
    cp templates/internal_load_balancer/examples/internal_load_balancer.yaml \
       my_internal_load_balancer.yaml
```

4. Change the values in the config file to match your specific GCP setup (for
   properties, refer to the schema files listed above):

```shell
    vim my_internal_load_balancer.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant
   deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_internal_load_balancer.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Internal Load Balancer](examples/internal_load_balancer.yaml)
