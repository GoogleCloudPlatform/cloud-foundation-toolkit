# Pubsub

This template creates a Pub/Sub (publish-subscribe) service.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [pubsub.admin](https://cloud.google.com/pubsub/docs/access-control)
IAM role to the Deployment Manager service account

## Deployment

### Resources

- [pubsub.v1.topic](https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.topics)
- [pubsub.v1.subscription](https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.subscriptions)

### Properties

See the `properties` section in the schema file(s):

- [Pub/Sub](pubsub.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples):

```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [community/cloud-foundation](../../) directory:

```shell
    cd community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment; in this case, [examples/pubsub.yaml](examples/pubsub.yaml):

```shell
    cp templates/pubsub/examples/pubsub.yaml my_pubsub.yaml
```

4. Change the values in the config file to match your specific GCP setup (for properties, refer to the schema files listed above):

```shell
    vim my_pubsub.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
    --config my_pubsub.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Pub/Sub](examples/pubsub.yaml)
- [Pub/Sub with PUSH subscription](examples/pubsub_push.yaml)
