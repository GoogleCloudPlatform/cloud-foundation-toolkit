# Owned Project
==================================================

## NAME

  owned-project

## SYNOPSIS

  Config Connector compatible YAML files to create
  a project in a folder, binding an IAM member
  as project owner.

## CONSUMPTION

1. Clone GoogleCloudPlatform/cloud-foundation-toolkit repository:

    ```bash
    git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
    ```

1. Go to the owned-project folder:

    ```bash
    cd cloud-foundation-toolkit/config-connector/solutions/projects/helm/owned-project
    ```

## REQUIREMENTS

1. GKE Cluster with Config Connector and [Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity#enable_workload_identity_on_a_new_cluster).
1. [Helm](../../../README.md#helm)
1. The "cnrm-system" service account assigned with
      -   `roles/resourcemanager.folderViewer`
      -   `roles/resourcemanager.projectCreator`
      -   `roles/iam.securityAdmin`
  in the target folder
1.   The IAM member selected below must meet the requirements specified
      [here](https://cloud.google.com/resource-manager/reference/rest/v1/projects/setIamPolicy#top_of_page).

## USAGE

All steps are run from the current directory ([config-connector/solutions/projects/helm/owned-project](.)).

1. Review and update the values in `./values.yaml`.

1. Validate and install the sample with Helm, with all values replaced with the ones you desire.

    ```bash
    # validate your chart
    helm lint . --set billingID=BILLING_ID,folderID=FOLDER_ID,iamMember=user:name@example.com,projectID=PROJECT_ID

    # check the output of your chart
    helm template . --set billingID=BILLING_ID,folderID=FOLDER_ID,iamMember=user:name@example.com,projectID=PROJECT_ID

    # do a dryrun on your chart and address issues if there are any
    helm install . --dry-run --set billingID=BILLING_ID,folderID=FOLDER_ID,iamMember=user:name@example.com,projectID=PROJECT_ID --generate-name

    # install your chart
    helm install . --set billingID=BILLING_ID,folderID=FOLDER_ID,iamMember=user:name@example.com,projectID=PROJECT_ID --generate-name
    ```

1. Check the created helm release to verify the installation:
    ```bash
    helm list
    ```
    
    Check the status of the project, where `PROJECT_ID` is the project ID value you gave above:
    ```bash
    kubectl describe gcpproject PROJECT_ID
    ```

    Check the status of the IAM Policy Member:
    ```bash
    kubectl describe iampolicymember owned-project-iampolicymember
    ```

1. Clean up the installation:

    ```bash
    # list Helm releases to obtain release name
    helm list

    # delete release specifying release name from the previous command output.
    helm delete [release_name]
    ```

## LICENSE

Apache 2.0 - See [LICENSE](/LICENSE) for more information.
