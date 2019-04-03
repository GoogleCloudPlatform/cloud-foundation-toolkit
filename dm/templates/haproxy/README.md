# HAProxy

This template:
- Creates a Compute Instance with an [HAProxy](http://www.haproxy.org/) installed
- Configures HAProxy to load-balance traffic between one or more of the provided
[instance groups](https://cloud.google.com/compute/docs/reference/rest/v1/instanceGroups)

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant [compute.viewer](https://cloud.google.com/compute/docs/access/iam) role to 
Compute Engine [default service account](https://cloud.google.com/compute/docs/access/service-accounts#compute_engine_default_service_account).
Alternatively, create a new service account with the above role, and add it to 
the template's `resources.properties.serviceAccountEmail` property.
- Create one or more [instanceGroups](https://cloud.google.com/compute/docs/reference/rest/v1/instanceGroups)
to be load-balanced, and add them to `resources.properties.instances.groups` collection.

## Deployment

### Resources

- [compute.v1.instance](https://cloud.google.com/compute/docs/reference/rest/v1/instances)

### Properties

See the `properties` section in the schema file(s):
- [HAProxy](haproxy.py.schema)

### Usage

1. Clone the [Deployment Manager Samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples):

```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [community/cloud-foundation](../../) directory:

```shell
    cd community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment; in this case, [examples/haproxy.yaml](examples/haproxy.yaml):

```shell
    cp templates/haproxy/examples/haproxy.yaml my_haproxy.yaml
```

4. Change the values in the config file to match your specific GCP setup (for properties, refer to the schema files listed above):

```shell
    vim my_haproxy.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_haproxy.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [HAProxy](examples/haproxy.yaml)
