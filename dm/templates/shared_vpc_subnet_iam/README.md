# Shared VPC Subnet IAM

This template grants IAM roles to a user on a shared VPC subnetwork.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Create a [network and subnetworks](../network/README.md)
- Grant the [compute.networkAdmin or compute.admin](https://cloud.google.com/compute/docs/access/iam) IAM role to the project service account

## Deployment

### Resources

- [gcp-types/compute-beta:compute.subnetworks.setIamPolicy](https://cloud.google.com/compute/docs/reference/rest/beta/subnetworks/setIamPolicy)
- [gcp-types/compute-beta:compute.subnetworks.getIamPolicy](https://cloud.google.com/compute/docs/reference/rest/beta/subnetworks/getIamPolicy)

### Properties

See `properties` section in the schema file(s):

-  [Shared VPC Subnet IAM](shared_vpc_subnet_iam.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit):

```shell
    git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit
```

2. Go to the [dm](../../) directory:

```shell
    cd dm
```

3. Copy the example DM config to be used as a model for the deployment; in this case, [examples/shared\_vpc\_subnet_iam.yaml](examples/shared_vpc_subnet_iam.yaml):

```shell
    cp templates/shared_vpc_subnet_iam/examples/shared_vpc_subnet_iam.yaml my_shared_vpc_subnet-iam.yaml
```

4. Change the values in the config file to match your specific GCP setup (for properties, refer to the schema files listed above):

```shell
    vim my_shared_vpc_subnet-iam.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
      --config my_shared_vpc_subnet-iam.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples
- [Shared VPC Subnet IAM Bindings syntax](examples/shared_vpc_subnet_iam_bindings.yaml)
- [Shared VPC Subnet IAM Policy syntax](examples/shared_vpc_subnet_iam_policy.yaml)
- [Shared VPC Subnet IAM Legacy](examples/shared_vpc_subnet_iam_legacy.yaml)

## Tests Cases
- [Shared VPC Subnet IAM Bindings syntax](tests/integration/bindings.bats)
- [Shared VPC Subnet IAM Policy syntax](tests/integration/policy.bats)
- [Shared VPC Subnet IAM Legacy syntax](tests/integration/legacy.bats)
