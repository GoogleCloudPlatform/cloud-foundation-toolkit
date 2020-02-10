# Route

This template creates a SAP HANA deployment. It supportes two type of deployments:
- A standalone HANA instance 
- A highly-available HANA deployment

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)


## Deployment

### Resources

- [compute.v1.instance](https://cloud.google.com/compute/docs/reference/rest/v1/instances)
- [compute.v1.network](https://cloud.google.com/compute/docs/reference/latest/networks)
- [compute.v1.subnetwork](https://cloud.google.com/compute/docs/reference/latest/subnetworks)


### Properties

See `properties` section in the schema file(s):

-  [HANA Template](sap_hana_template.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit):

```shell
    git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit
```

2. Go to the [dm](../../) directory:

```shell
    cd dm
```

3. Copy the example DM config to be used as a model for the deployment; in this case, [examples/hana_standalone_scenario.yaml](examples/hana_standalone_scenario.yaml) for standalone HANA or [examples/hana_ha_scenario.yaml](examples/hana_ha_scenario.yaml) for HANA HA scenario :

```shell
    cp templates/sap_hana/examples/hana_ha_scenario.yaml my_hana_ha.yaml
```

4. Change the values in the config file to match your specific GCP setup (for properties, refer to the schema files listed above):

```shell
    vim my_hana_ha.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
    --config my_hana_ha.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## How to test it?
You need to perform the fthe testing steps are the following: 
1) Deploy the template as described in step 5 above
2) Wait until HANA is deployed successfully. This can be done by checking the Stackdriver Logs for ""INSTANCE DEPLOYMENT COMPLETE".
3) Connext vis SSH to the Linux bastion host 
4) From the bastion host, connect to the HANA instance using the command ```gcloud compute INSTANCE-NAME --internal-ip```, then
4.a) For a standalone HANA deployment, follow the steps described [here](https://cloud.google.com/solutions/sap/docs/sap-hana-deployment-guide#verifying_deployment).
4.b) For an HANA HA deployment, follow the steps described [here](https://cloud.google.com/solutions/sap/docs/sap-hana-ha-deployment-guide#checking_the_configuration_of_the_vm_and_the_sap_hana_installation).

## Examples

- [A standlone SAP Hana deployment](examples/hana_ha_scenario.yaml)
- [A highly-available SAP HANA deployment](examples/hana_standalone_scenario.yaml)
