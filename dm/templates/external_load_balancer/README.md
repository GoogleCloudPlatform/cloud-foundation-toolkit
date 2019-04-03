# External Load Balancer

This template creates an HTTP(S), SSL Proxy, or TCP Proxy external load balancer.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [compute.loadBalancerAdmin](https://cloud.google.com/compute/docs/access/iam)
  IAM role to the Deployment Manager service account
- For using the TCP Proxy load balancing, request access to the Compute ALPHA features
  from the Cloud [Support](https://cloud.google.com/support/)

## Deployment

### Resources

- [compute.v1.forwardingRule](https://cloud.google.com/compute/docs/reference/latest/forwardingRules)
- [compute.v1.targetHttpProxy](https://cloud.google.com/compute/docs/reference/latest/targetHttpProxies)
- [compute.v1.targetHttpsProxy](https://cloud.google.com/compute/docs/reference/latest/targetHttpsProxies)
- [compute.alpha.targetTcpProxy](https://www.googleapis.com/discovery/v1/apis/compute/alpha/rest)
- [compute.v1.targetSslProxy](https://cloud.google.com/compute/docs/reference/latest/targetSslProxies)
- [compute.v1.backendService](https://cloud.google.com/compute/docs/reference/rest/v1/backendServices)
- [compute.v1.sslCertificate](https://cloud.google.com/compute/docs/reference/rest/v1/sslCertificates)
- [compute.v1.urlMap](https://cloud.google.com/compute/docs/reference/rest/v1/urlMaps)

### Properties

See the `properties` section in the schema file(s):

- [External Load Balancer](external_load_balancer.py.schema)

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
   case, [examples/external\_load\_balancer.yaml](examples/external_load_balancer.yaml):

```shell
    cp templates/external_load_balancer/examples/external_load_balancer.yaml \
       my_external_load_balancer.yaml
```

4. Change the values in the config file to match your specific GCP setup (for
   properties, refer to the schema files listed above):

```shell
    vim my_external_load_balancer.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace \<YOUR\_DEPLOYMENT\_NAME\> with the relevant
   deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_external_load_balancer.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [External HTTP Load Balancer](examples/external_load_balancer_http.yaml)
- [External HTTPS Load Balancer](examples/external_load_balancer_https.yaml)
- [External SSL Load Balancer](examples/external_load_balancer_ssl.yaml)
- [External TCP Load Balancer](examples/external_load_balancer_tcp.yaml)
