# Demo script for Deployment Manager

This script is part of the Take5 demo for **Deployment Manager**.
This tutorial walks you through how to start with **Deployment Manager**
and how to use the **Cloud Foundation Toolkit**.

The video will be published shortly.

## Part 1 - Firewall rules

```bash
# Clone the CFT Repo
git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
cd cloud-foundation-toolkit/dm

# Copy the Firewall config file
cp templates/firewall/examples/firewall.yaml my_firewall.yaml

# Specifying the VPC for the Firewall rules
export VpcName=default

# Edit my_firewall.yaml - Change the network name manually or via CLI: 
sed -i "s/<FIXME:network-name>/$VpcName/g" my_firewall.yaml

# Enable the compute API to create firewalls, enable the deploymentmanager API to create deployments.
gcloud services enable compute.googleapis.com deploymentmanager.googleapis.com
gcloud deployment-manager deployments create my-first-firewalls --config my_firewall.yaml

# Manually change the my_firewall.yaml and try out the changes
gcloud deployment-manager deployments update my-first-firewalls --config my_firewall.yaml

# Clean up the deployment
gcloud deployment-manager deployments delete my-first-firewalls

```
## Part 2 - Project Factory with Shared-VPC

```bash
# Prerequisites - Setting the environment specific values
export OrgID=518838582042
export ProjectNumber=700306896797
export BillingID=01BACD-32281D-31B750
export ProjectUniqueNameH=take5-host-xt-1300
export ProjectUniqueNameG=take5-host-xg-1300
export ParentFolderID=1049237988874

# Enabling the required APIs and IAM permissions
gcloud services enable deploymentmanager.googleapis.com cloudresourcemanager.googleapis.com cloudbilling.googleapis.com iam.googleapis.com servicemanagement.googleapis.com
gcloud organizations add-iam-policy-binding $OrgID --member=serviceAccount:$ProjectNumber@cloudservices.gserviceaccount.com --role=roles/resourcemanager.projectCreator

## Add <ProjectNumber>@cloudservices.gserviceaccount.com to the billing account as Billing User MANUALLY

cp templates/project/examples/project.yaml my_project.yaml

# Edit my_firewall.yaml - Change the Org/Folder ID, BillingID, UniqueProjectName manually or change it via CLI: 
sed -i "s/<FIXME:UniqueProjectName>/$ProjectUniqueNameH/g" my_project.yaml
sed -i "s/type: organization/type: folder/g" my_project.yaml
sed -i "s/<FIXME:OrgID>/$ParentFolderID/g" my_project.yaml
sed -i "s/<FIXME:BillingAccount:ID>/$BillingID/g" my_project.yaml

# Manual remove attachment to a shared VPC from my_project.yaml

# Create the project
gcloud deployment-manager deployments create my-first-project --config my_project.yaml

# add `sharedVPCHost: true` to my_project.yaml

# Enable the Deployment Manager SA to attach projects to the shared VPC
gcloud organizations add-iam-policy-binding $OrgID --member=serviceAccount:$ProjectNumber@cloudservices.gserviceaccount.com --role=roles/compute.xpnAdmin

# Update the project to a Shared-VPC host project
gcloud deployment-manager deployments update my-first-project --config my_project.yaml

cp templates/network/examples/network.yaml my_network.yaml

# Create the Shared-VPC and its subnets
gcloud deployment-manager deployments create my-first-network --config my_network.yaml --project $ProjectUniqueNameH


cp templates/project/examples/project.yaml my_guest_project.yaml


# nano my_guest_project.yaml - Change the network name and other values manually or change it via CLI: 
sed -i "s/<FIXME:UniqueProjectName>/$ProjectUniqueNameG/g" my_guest_project.yaml
sed -i "s/type: organization/type: folder/g" my_guest_project.yaml
sed -i "s/<FIXME:OrgID>/$ParentFolderID/g" my_guest_project.yaml
sed -i "s/<FIXME:BillingAccount:ID>/$BillingID/g" my_guest_project.yaml
sed -i "s/test-vpc-host-project/$ProjectUniqueNameH/g" my_guest_project.yaml
sed -i "s/subnet-1/test-subnetwork-1/g" my_guest_project.yaml

# Create the guest project and attach it to the host project
gcloud deployment-manager deployments create my-guest-project --config my_guest_project.yaml
```

## Note

- Some templates are updated since the recording of the video, there are more detailed examples
 available for the project template