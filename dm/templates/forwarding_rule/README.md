# Forwarding Rule

This template creates a forwarding rule.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [compute.admin](https://cloud.google.com/compute/docs/access/iam) or
[compute.networkAdmin](https://cloud.google.com/compute/docs/access/iam) IAM
role to the Deployment Manager service account

## Deployment

### Resources

- [compute.v1.globalForwardingRule](https://cloud.google.com/compute/docs/reference/latest/globalForwardingRules)
- [compute.v1.forwardingRule](https://cloud.google.com/compute/docs/reference/latest/forwardingRules)

### Properties

See the `properties` section in the schema file(s):
- [Forwarding Rule](forwarding_rule.py.schema)

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
case, [examples/forwarding\_rule\_global.yaml](examples/forwarding_rule_global.yaml):

```shell
    cp templates/forwarding_rule/examples/forwarding_rule_global.yaml \
       my_forwarding_rule.yaml
```

4. Change the values in the config file to match your specific GCP setup (for
properties, refer to the schema files listed above):

```shell
    vim my_forwarding_rule.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant
deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
    --config my_forwarding_rule.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Global Forwarding Rule](examples/forwarding_rule_global.yaml)
- [Regional Forwarding Rule](examples/forwarding_rule_regional.yaml)
