# Healthcheck

This template creates a load balancer healthcheck.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Create a [network](../network/README.md)
- Grant the [compute.networkAdmin](https://cloud.google.com/compute/docs/access/iam)
 IAM role to the project service account

## Deployment

### Resources

Depend on the specified healthcheck type.

#### Legacy Healthchecks

- [compute.v1.httpHealthCheck](https://cloud.google.com/sdk/gcloud/reference/compute/health-checks/create/http)
- [compute.v1.httpsHealthCheck](https://cloud.google.com/sdk/gcloud/reference/compute/health-checks/create/https)

#### TCP + SSL Healthchecks

- [compute.v1.healthChecks](https://cloud.google.com/load-balancing/docs/health-check-concepts)

#### Beta Healthchecks

- [compute.beta.healthChecks](https://cloud.google.com/sdk/gcloud/reference/beta/compute/health-checks/create/http2)

### Properties

See the `properties` section in the schema file(s):

- [Healthcheck](healthcheck.py.schema)

### Usage

1. Clone the [Deployment Manager Samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples):

```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [community/cloud-foundation](../../) directory:

```shell
    cd community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment;
 in this case, [examples/healthcheck.yaml](examples/healthcheck.yaml):

```shell
    cp templates/healthcheck/examples/healthcheck.yaml my_healthcheck.yaml
```

4. Change the values in the config file to match your specific GCP setup:

```shell
    vim my_healthcheck.yaml  # <== change values to match your GCP setup
```

5. Create your deployment:

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_healthcheck.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Healthcheck](examples/healthcheck.yaml)