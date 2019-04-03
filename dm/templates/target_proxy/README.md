# Target Proxy

This template creates one of the following target proxy resources (depending on the parameters):

- targetHttpProxy
- targetHttpsProxy
- targetTcpProxy
- targetSslProxy

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [compute.loadBalancerAdmin](https://cloud.google.com/compute/docs/access/iam)
  IAM role to the Deployment Manager service account
- To use the target TCP Proxy, request access to the Compute ALPHA features from the
  Cloud [Support](https://cloud.google.com/support/)

## Deployment

### Resources

- [compute.v1.targetHttpProxy](https://cloud.google.com/compute/docs/reference/latest/targetHttpProxies)
- [compute.v1.targetHttpsProxy](https://cloud.google.com/compute/docs/reference/latest/targetHttpsProxies)
- [compute.alpha.targetTcpProxy](https://www.googleapis.com/discovery/v1/apis/compute/alpha/rest)
- [compute.v1.targetSslProxy](https://cloud.google.com/compute/docs/reference/latest/targetSslProxies)

### Properties

See the `properties` section in the schema file(s):

- [Target Proxy](target_proxy.py.schema)

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
   case, [examples/target\_proxy.yaml](examples/target_proxy.yaml):

```shell
    cp templates/target_proxy/examples/target_proxy.yaml \
       my_target_proxy.yaml
```

4. Change the values in the config file to match your specific GCP setup (for
   properties, refer to the schema files listed above):

```shell
    vim my_target_proxy.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace \<YOUR\_DEPLOYMENT\_NAME\> with the relevant
   deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_target_proxy.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Target HTTP Proxy](examples/target_proxy_http.yaml)
- [Target HTTPS Proxy](examples/target_proxy_https.yaml)
- [Target TCP Proxy](examples/target_proxy_tcp.yaml)
- [Target SSL Proxy](examples/target_proxy_ssl.yaml)
