# Simple Project

==================================================

## NAME

  simple-project

## SYNOPSIS

  Config Connector compatible YAML files to create a project in an organization.

## CONSUMPTION

  1. Clone GoogleCloudPlatform/cloud-foundation-toolkit repository:
  
      ```bash
      git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
      ```

  1. Go to the simple-project folder:

      ```bash
      cd cloud-foundation-toolkit/config-connector/solutions/projects/helm/simple-project
      ```

## REQUIREMENTS

1. GKE Cluster with Config Connector and [Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity#enable_workload_identity_on_a_new_cluster).

1. A working Config Connector cluster using the "cnrm-system" service account with _minimally_ the permissions given by the following role in the desired organization:
    - `roles/resourcemanager.projectCreator`

1. [Helm](../../../README.md#helm)

## USAGE

All steps are run from the current directory ([config-connector/solutions/projects/helm/simple-project](.)).

1. Review and update the values in `./values.yaml`.

1. Validate and install the sample with Helm. `ORG_ID` should be your organization ID, `PROJECT_ID` should be a new project ID unique within GCS, and `BILLING_ID` should be your desired billing ID for the new project.

    ```bash
    # validate your chart
    helm lint . --set billingID=BILLING_ID,orgID=ORG_ID,projectID=PROJECT_ID

    # check the output of your chart
    helm template . --set billingID=BILLING_ID,orgID=ORG_ID,projectID=PROJECT_ID

    # do a dryrun on your chart and address issues if there are any
    helm install . --dry-run --set billingID=BILLING_ID,orgID=ORG_ID,projectID=PROJECT_ID --generate-name

    # install your chart
    helm install . --set billingID=BILLING_ID,orgID=ORG_ID,projectID=PROJECT_ID --generate-name
    ```

1. Check the created helm release to verify the installation:
    ```bash
    helm list
    ```

    Check the status of the applied YAML:
    ```bash
    kubectl describe gcpprojects PROJECT_ID
    ```
    where `PROJECT_ID` is your project ID above.

1. Clean up the installation:

    ```bash
    # list Helm releases to obtain release name
    helm list

    # delete release specifying release name from the previous command output.
    helm delete [release_name]
    ```

## LICENSE

Apache 2.0 - See [LICENSE](/LICENSE) for more information.
