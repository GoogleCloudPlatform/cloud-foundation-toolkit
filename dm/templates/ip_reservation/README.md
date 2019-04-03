# IP Reservation

This template creates an IP reservation.
Depending on the input option, the following addresses can be reserved:
- Global
- External
- Internal

## Prerequisites
- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Grant the [compute.networkAdmin](https://cloud.google.com/compute/docs/access/iam) IAM role to the project service account (unless the default Project Editor role is already granted)


## Deployment

### Resources

- [gcp-types/compute-v1:address](https://cloud.google.com/compute/docs/reference/rest/v1/addresses)
- [gcp-types/compute-v1:globalAddress](https://cloud.google.com/compute/docs/reference/rest/v1/addresses)


### Properties

See the `properties` section in the schema file(s):
-  [IP Reservation](ip_reservation.py.schema)


#### Usage

1. Clone the [Deployment Manager samples_repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples)

```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [community/cloud-foundation](../../../cloud-foundation) directory

```shell
    cd community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment, in this case [examples/ip_reservation.yaml](examples/ip_reservation.yaml)

```shell
    cp templates/ip_reservation/examples/ip_reservation.yaml my_ip_reservation.yaml
```

4. Change the values in the config file to match your specific GCP setup.
   Refer to the properties in the schema files described above.

```shell
    vim my_ip_reservation.yaml  # <== change values to match your GCP setup
```

5. Create your deployment as described below, replacing <YOUR_DEPLOYMENT_NAME>
   with your with your own deployment name

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
    --config my_ip_reservation.yaml
```

6. In case you need to delete your deployment: 

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Reserving a global, external, or internal IP address](examples/ip_reservation.yaml)
