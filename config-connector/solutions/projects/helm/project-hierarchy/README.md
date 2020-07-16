# Project Hierarchy
==================================================

## NAME

  project-hierarchy

## SYNOPSIS

  Config Connector compatible YAML files to create
  a folder in an organization, and a project
  beneath it.

## CONSUMPTION

1. Clone GoogleCloudPlatform/cloud-foundation-toolkit repository:

    ```bash
    git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
    ```

1. Go to the project-hierarchy folder:

    ```bash
    cd cloud-foundation-toolkit/config-connector/solutions/projects/helm/project-hierarchy
    ```

## REQUIREMENTS

1. GKE Cluster with Config Connector and [Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity#enable_workload_identity_on_a_new_cluster).
1. [Helm](../../../README.md#helm)
1. The "cnrm-system" service account assigned with
      -   `roles/resourcemanager.folderCreator`
      -   `roles/resourcemanager.projectCreator`
  in the target organization

## USAGE

All steps are run from the current directory ([config-connector/solutions/projects/helm/project-hierarchy](.)).

1. Review and update the values in `./folder/values.yaml` and `./project/values.yaml`, except folderID, which you will find in a later step.

1. Install the folder Helm chart:

    ```bash
    # validate your chart
    helm lint ./folder
    
    # do a dryrun on your chart and address issues if there are any
    helm install ./folder --dry-run --generate-name

    # install your chart
    helm install ./folder --generate-name
    ```

1. Check the created helm release to verify the installation:
    
    ```bash
    helm list
    ```
    
    Check the status of the folder you just created:
    ```bash
    kubectl describe gcpfolder project-hierarchy-folder
    ```
    
1. Find the ID of the created folder. If you replaced the `folderName` value in step 1, replace `project-hierarchy-folder` below with the value you chose:

    ```bash
    kubectl describe gcpfolder project-hierarchy-folder | grep Name:\ *folders\/ | sed "s/.*folders\///"
    ```
    
1. Update `./project/values.yaml` with this folderID value.
    
1. Install the project helm chart:

    ```bash
    # validate your chart
    helm lint ./project

    # do a dryrun on your chart and address issues if there are any
    helm install ./project --dry-run --generate-name

    # install your chart
    helm install ./project --generate-name
    ```

1. Clean up the installations:

    ```bash
    # list Helm releases to obtain release name
    helm list

    # delete releases specifying release names from the previous command output.
    helm delete [release_name]
    helm delete [release_name]
    ```

## LICENSE

Apache 2.0 - See [LICENSE](/LICENSE) for more information.
