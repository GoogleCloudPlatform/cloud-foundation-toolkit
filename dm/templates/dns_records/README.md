# DNS Resource RecordSets

This template creates Cloud DNS records using recordsets.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [dns.admin](https://cloud.google.com/dns/access-control) IAM role to the Deployment Manager `serviceAccount`

## Deployment

### Resources

- [gcp-types/dns-v1](https://cloud.google.com/dns/api/v1/changes)

### Properties

See the `properties` section in the schema file(s):

- [DNS records](dns_records.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples):

    ```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
    ```

2. Go to the [community/cloud-foundation](../../) directory:

    ```shell
    cd community/cloud-foundation
    ```

3. Copy the example DM config to be used as a model for the deployment; in this case, [examples/dns_records.yaml](examples/dns_records.yaml):

    ```shell
    cp templates/dns_records/examples/dns_records.yaml my_dns_records.yaml
    ```

4. Change the values in the config file to match your specific GCP setup (for properties, refer to the schema files listed above):

    ```shell
    vim my_dns_records.yaml  # <== change values to match your GCP setup
    ```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

    ```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_dns_records.yaml
    ```

6. In case you need to delete your deployment:

    ```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
    ```

## Examples

- [DNS records](examples/dns_records.yaml)
