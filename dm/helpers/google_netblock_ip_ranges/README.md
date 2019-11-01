# Important Google IP ranges helper

This helper creates firewall template rules for a network with Google important ranges.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Create a [network](../network/README.md)
- Grant the [compute.networkAdmin or compute.securityAdmin](https://cloud.google.com/compute/docs/access/iam) IAM role to the project service account

## Deployment

### Resources

- [compute.beta.firewall](https://cloud.google.com/compute/docs/reference/rest/beta/firewalls)
  
  `Note:` The beta API supports the firewall log feature.

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit):

```shell
    git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit
    cd cloud-foundation-toolkit
```

2. Go to the [dm](../../) directory:

```shell
    cd dm
```

3. Copy the example DM config to be used as a model for the deployment; in this case, with [firewall template](../../templates/firewall/firewall.py):

```shell
    cp helpers/google_netblock_ip_ranges/examples/google_netblock_ip_ranges_example.yaml google_netblock_ip_ranges_example.yaml
```

4. Change the values in the config file to match your specific GCP setup (like properties).
   Name of the imported YAML-file with important IP ranges must be exact "google_netblock_ip_ranges.yaml":

```shell
    vim google_netblock_ip_ranges_example.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config google_netblock_ip_ranges_example.yaml
```

## Examples

- [Firewall](examples/google_netblock_ep_ranges_example.yaml)
