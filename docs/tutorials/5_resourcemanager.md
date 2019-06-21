#  Installation

## Introduction

<walkthrough-tutorial-duration duration="10"></walkthrough-tutorial-duration>

This tutorial explains how to set up a hierarchy for managing resources in GCP.

## Choose Project

<walkthrough-project-billing-setup billing="true"></walkthrough-project-billing-setup>

First select a project to install Forseti in.
This can either be a dedicated Forseti project or an existing DevSecOps project.

## Configure Forseti
To install Forseti, you will need to update a few settings in the <walkthrough-editor-open-file filePath="terraform-google-forseti/examples/install_simple/terraform.tfvars">terraform.tfvars</walkthrough-editor-open-file>.

## Activate APIs

You will need to activate a few APIs on this project for Forseti to function:
<walkthrough-enable-apis apis=
  "cloudresourcemanager.googleapis.com,
  serviceusage.googleapis.com,
  compute.googleapis.com"></walkthrough-enable-apis>

### Set project
On line 1, update the <walkthrough-editor-select-regex
  filePath="terraform-google-forseti/examples/install_simple/terraform.tfvars"
  regex="my-project-id">project ID</walkthrough-editor-select-regex>
to match your chosen project (`{{project_id}}`).

### Set organization ID
On line 2, update the <walkthrough-editor-select-regex
  filePath="terraform-google-forseti/examples/install_simple/terraform.tfvars"
  regex="11111111">organization ID</walkthrough-editor-select-regex>
to match your organization ID, which can be found in the URL bar.

### Set domain
On line 3, update the <walkthrough-editor-select-regex
  filePath="terraform-google-forseti/examples/install_simple/terraform.tfvars"
  regex="mydomain.com">domain</walkthrough-editor-select-regex>
to match your company Cloud Identity domain, which can be found in the URL bar.

### Choose region
On line 5, update the <walkthrough-editor-select-regex
  filePath="terraform-google-forseti/examples/install_simple/terraform.tfvars"
  regex="us-east4">region</walkthrough-editor-select-regex>
you wish to deploy Forseti in.

### Choose network
On line 6, update the <walkthrough-editor-select-regex
  filePath="terraform-google-forseti/examples/install_simple/terraform.tfvars"
  regex="default">network</walkthrough-editor-select-regex>
you wish to deploy Forseti in.
You also need to update the <walkthrough-editor-select-line
  filePath="terraform-google-forseti/examples/install_simple/terraform.tfvars"
  startLine=6
  endLine=6
  startCharacterOffset=12
  endCharacterOffset=19>subnetwork</walkthrough-editor-select-line>
on line 7.

If you are deploying on a Shared VPC, you need to set the <walkthrough-editor-select-line
  filePath="terraform-google-forseti/examples/install_simple/terraform.tfvars"
  startLine=7
  endLine=7
  startCharacterOffset=17
  endCharacterOffset=17>network project</walkthrough-editor-select-line>
on line 8. Otherwise, you can leave this empty.

## Enable Optional Features
There are additional settings which you can configure in the settings file to enable advanced Forseti functionality.

If you don't need these features, you can skip these steps.

### Configure G Suite
On line 10, set the <walkthrough-editor-select-regex
  filePath="terraform-google-forseti/examples/install_simple/terraform.tfvars"
  regex="admin@mydomain.com">G Suite super admin email</walkthrough-editor-select-regex>.
Ask your G Suite Admin if you don’t know the super admin email.

This is part of the [G Suite data collection](https://forsetisecurity.org/docs/latest/configure/inventory/gsuite.html). The following functionalities will not work without G Suite integration:

- G Suite groups and users in Inventory
- Group Scanner
- Group expansion in Explain

### Configure email notifications
Forseti can be configured to [send email notifications](https://forsetisecurity.org/docs/latest/configure/notifier/index.html#email-notifications).

To enable this, you need to add a <walkthrough-editor-select-line
  filePath="terraform-google-forseti/examples/install_simple/terraform.tfvars"
  startLine=10
  endLine=10
  startCharacterOffset=18
  endCharacterOffset=18>SendGrid API key</walkthrough-editor-select-line>
on line 12 and update the <walkthrough-editor-select-line
  filePath="terraform-google-forseti/examples/install_simple/terraform.tfvars"
  startLine=11
  endLine=11
  startCharacterOffset=22
  endCharacterOffset=22>sender</walkthrough-editor-select-line>
and <walkthrough-editor-select-line
  filePath="terraform-google-forseti/examples/install_simple/terraform.tfvars"
  startLine=12
  endLine=12
  startCharacterOffset=25
  endCharacterOffset=25>recipient</walkthrough-editor-select-line>
settings.

## Install Forseti
Now that you have updated your configuration settings, you are ready to install Forseti.
This will be done using Terraform, which comes preinstalled with this Cloud Shell.

### Initialize Terraform
To download the Forseti module, you will need to initialize Terraform:
```bash
terraform init
```

### Apply Terraform
You are now ready to install Forseti with Terraform by running the apply command:

```bash
terraform apply -auto-approve
```

This can take a few minutes as all the necessary resources are provisioned.

If you encounter errors during installation, you can check your configuration and permissions, then run `terraform apply` again.

## Save state to GCS
Congratulations, you have now installed Forseti.
As a final step, you will want to save your configuration so it can be used to upgrade Forseti in the future.

### Create Terraform state bucket
Create a Google Cloud Storage bucket to [store your Terraform state](https://www.terraform.io/docs/state/).

```bash
gsutil mb gs://{{project_id}}-tfstate
```

### Update state configuration
Open <walkthrough-editor-open-file filePath="terraform-google-forseti/examples/install_simple/backend.tf">backend.tf</walkthrough-editor-open-file> and uncomment the contents.

On line 3, change the <walkthrough-editor-select-regex
  filePath="terraform-google-forseti/examples/install_simple/backend.tf"
  regex="my-project">project ID</walkthrough-editor-select-regex>
project ID to match your project ID (`{{project_id}}`).

Finally, re-initialize Terraform to upload your state to Cloud Storage:

```bash
terraform init
```

At the prompt, type `yes`.

## Save configuration to git
As a best practice, you should save your Terraform configuration to source control. This can be done using Cloud Source Repositories.

### Create a repo
```bash
gcloud source repos create terraform-forseti
```

### Initalize git
```bash
git init
```

### Add and commit your files

```bash
git add -A
```

```bash
git commit -m "Initial commit"
```

### Push your configuration
```bash
git remote add origin https://source.developers.google.com/p/{{project_id}}/r/terraform-forseti
```

```bash
git push origin master
```
