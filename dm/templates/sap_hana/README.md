# SAP HANA

This template creates an SAP HANA deployment. It supports two types of deployments:
- A standalone HANA instance 
![GitHub Logo](/dm/sap_hana/images/sap_hana_standalone.png)


- A highly-available HANA deployment
![GitHub Logo](/images/sap_hana_ha.png)
Format: ![Alt Text](url)

Below are the main steps performed by the template: 
- Create a custom VPC  with two subnets: 
    - subnetwork-1 will be used as DMZ and the bastion host will be deployed into it.
    - subnetwork-2: this is where the HANA DB will be deployed.
- Set up a NAT gateway: so that your VMs can access the internet without having to have a public IP address.
- Create necessary firewall rules, to allow connectivity between bastion hosts and the instance where HANA is deployed.
- Install SAP HANA using by leveraging existing [templates](https://cloud.google.com/solutions/sap/docs/sap-hana-deployment-guide).
- Deploy two bastion-hosts. Once which is windows-based 


## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)


## Deployment

### Resources

- [compute.v1.instance](https://cloud.google.com/compute/docs/reference/rest/v1/instances)
- [compute.v1.network](https://cloud.google.com/compute/docs/reference/latest/networks)
- [compute.v1.subnetwork](https://cloud.google.com/compute/docs/reference/latest/subnetworks)
- [compute.beta.firewall](https://cloud.google.com/compute/docs/reference/rest/beta/firewalls)
- [compute.v1.router](https://cloud.google.com/compute/docs/reference/rest/v1/routers)


### Properties

See `properties` section in the schema file(s):

-  [SAP HANA template](sap_hana_template.py.schema)

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
    cp templates/sap_hana/examples/sap_hana_ha.yaml my_sap_hana_ha.yaml
```

4. Change the values in the config file to match your specific GCP setup (for properties, refer to the schema files listed above):

```shell
    vim my_sap_hana_ha.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
    --config my_sap_hana_ha.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## How to test it?
You need to perform the the testing steps are the following: 
1) Deploy the template as described in step 5 above
2) Wait until HANA is deployed successfully. This can be done by checking the Stackdriver Logs for ""INSTANCE DEPLOYMENT COMPLETE".
3) Connect vis SSH to the Linux bastion host 
4) From the bastion host, connect to the HANA instance using the command ```gcloud compute INSTANCE-NAME --internal-ip```, then:
    * For a standalone HANA deployment, follow the steps described [here](https://cloud.google.com/solutions/sap/docs/sap-hana-deployment-guide#verifying_deployment).
    * For a HANA HA deployment, follow the steps described [here](https://cloud.google.com/solutions/sap/docs/sap-hana-ha-deployment-guide#checking_the_configuration_of_the_vm_and_the_sap_hana_installation).

## Examples

- [A standlone SAP Hana deployment](examples/sap_hana_ha.yaml)
- [A highly-available SAP HANA deployment](examples/sap_hana_standalone.yaml)
