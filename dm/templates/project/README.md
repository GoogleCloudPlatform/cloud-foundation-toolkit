# Project

This template:

1. Creates a new project.
2. Sets a billing account for the new project
3. Sets IAM permissions in the new project
4. Turns on a set of APIs in the new project
5. Creates service accounts for the new project
6. Creates an usage export Cloud Storage bucket for the new project
7. Removed default networks, firewalls
8. Removes default Service Account
9. Creates VPC host or attached VPC service project

## Prerequisites

Following are the prerequisites for creating a project via Deployment Manager. You can perform some of the steps via the Cloud Console at https://console.cloud.google.com/. The `gcloud` command line tool is used to deploy the configs.

`Note:` Permission changes can take up to 20 minutes to propagate. If you run commands before the propagation is completed, you may receive errors regarding the user not having permissions.

1. Install [gcloud](https://cloud.google.com/sdk).

2.  Create a project that will create and own the deployments (henceforth referred to as *DM Creation Project*). See:  https://cloud.google.com/resource-manager/docs/creating-managing-organization.
    
    `Important:` Because of the special permissions granted to the *DM Creation Project*, it should not be used for any purpose other than creating other projects.

3.  Activate the following APIs for the *DM Creation Project*:
    * Google Cloud Deployment Manager V2 API
    * Google Cloud Resource Manager API
    * Google Cloud Billing API
    * Google Identity and Access Management (IAM) API
    * Google Service Management API

    You may use the `gcloud services enable` command to do this:

    ```shell
    gcloud services enable deploymentmanager.googleapis.com
    gcloud services enable cloudresourcemanager.googleapis.com
    gcloud services enable cloudbilling.googleapis.com
    gcloud services enable iam.googleapis.com
    gcloud services enable servicemanagement.googleapis.com
    ```

4.  Find the *Cloud Services* service account associated with the *DM Creation Project*.

    It is formatted as `<project_number>@cloudservices.gserviceaccount.com`,
    and is listed under [IAM & Admin](https://console.cloud.google.com/iam-admin/iam)
    in Google Cloud Console. This account is henceforth referred to as the *DM Service Account*. See https://cloud.google.com/resource-manager/docs/access-control-proj.

5.  Create an Organization node.

    If you do not already have an Organization node under which you can create
    projects, create that node following [these instructions](https://cloud.google.com/resource-manager/docs/creating-managing-organization).

6.  Grant the *DM Service Account* the following permissions on the Organization node:
`roles/resourcemanager.projectCreator`. This is visible in the Cloud Console's IAM permissions in *Resource Manager -> Project Creator*. See https://cloud.google.com/resource-manager/docs/access-control-proj.

7.  Create/find the *Billing Account* associated with the Organization. See: https://cloud.google.com/support/billing/. Take note of the *Billing Account*'s ID, which is formatted as follows:`00E12A-0AB8B2-078CE8`.

8.  Give the *DM Service Account* the following permissions on the *Billing Account*: `roles/billing.user`. This is visible in Cloud Console's IAM permissions in *Billing -> Billing Account User*.

9.  If the project is a VPC host project, give the *DM Service Account* the following permissions: `roles/compute.xpnAdmin`.

## Deployment

### Resources

- [cloudresourcemanager.v1.project](https://cloud.google.com/compute/docs/reference/latest/projects)
- [deploymentmanager.v2.virtual.projectBillingInfo](https://cloud.google.com/billing/reference/rest/v1/projects/updateBillingInfo)
- [iam.v1.serviceAccount](https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts)
- [deploymentmanager.v2.virtual.enableService](https://cloud.google.com/service-management/reference/rest/v1/services/enable)
- [../iam_member CFT temaplet](../iam_member/README.md)
- [gcp-types/cloudresourcemanager-v1:cloudresourcemanager.projects.setIamPolicy](https://cloud.google.com/deployment-manager/docs/configuration/supported-gcp-types)
- [gcp-types/storage-v1:buckets](https://cloud.google.com/deployment-manager/docs/configuration/supported-gcp-types)
- [gcp-types/compute-v1:compute.projects.setUsageExportBucket](https://cloud.google.com/deployment-manager/docs/configuration/supported-gcp-types)
- [compute.beta.xpnResource](https://cloud.google.com/compute/docs/reference/rest/beta/projects/enableXpnResource)
- [compute.beta.xpnHost](https://cloud.google.com/compute/docs/reference/rest/beta/projects/enableXpnHost)
- [gcp-types/compute-beta:compute.firewalls.delete](https://cloud.google.com/compute/docs/reference/rest/beta/firewalls)
- [gcp-types/compute-beta:compute.networks.delete](https://cloud.google.com/compute/docs/reference/rest/beta/networks)
- [gcp-types/iam-v1:iam.projects.serviceAccounts.delete](https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts)

### Properties

See the `properties` section in the schema file(s):

-  [project](project.py.schema)

### Usage

1. Clone the [Deployment Manager Samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples):

```shell
    git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
```

2. Go to the [community/cloud-foundation](../../) directory:

```shell
    cd community/cloud-foundation
```

3. Copy the example DM config to be used as a model for the deployment; in this case, [examples/project.yaml](examples/project.yaml):

```shell
    cp templates/project/examples/project.yaml my_project.yaml
```

4. Change the values in the config file to match your specific GCP setup (for properties, refer to the schema files listed above):

```shell
    vim my_project.yaml  # <== change values to match your GCP setup
```

5. Create your deployment (replace <YOUR_DEPLOYMENT_NAME> with the relevant deployment name):

```shell
    gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
        --config my_project.yaml
```

6. In case you need to delete your deployment:

```shell
    gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
```

## Examples

- [Project](examples/project.yaml)
