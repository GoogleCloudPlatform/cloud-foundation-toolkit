# Project Services

==================================================

## NAME

  project-services

## SYNOPSIS

  Config Connector compatible YAML files to enable services on a desired project.

## CONSUMPTION

  1. Clone GoogleCloudPlatform/cloud-foundation-toolkit repository:

      ```bash
      git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
      ```

  1. Go to the project-services folder:

      ```bash
      cd cloud-foundation-toolkit/config-connector/solutions/projects/helm/project-services
      ```

## REQUIREMENTS

1. GKE Cluster with [Config Connector installed using a GKE Workload Identity](https://cloud.google.com/config-connector/docs/how-to/install-upgrade-uninstall#workload-identity).

1. Cloud Resource Manager API enabled in the project where Config Connector is installed.

1. The "cnrm-system" service account that has the `roles/serviceusage.serviceUsageAdmin` or `roles/owner` role in your desired project.

1. [Helm](../../../README.md#helm)

## USAGE

All steps are run from the current directory ([config-connector/solutions/projects/helm/project-services](.)).

1. Review and update the values in `./values.yaml`.

1. Create a namespace. If you want to use your existing namespace skip this step and use your own namespace name instead of `project-annotated` in all other steps.

    ```bash
    kubectl create namespace project-annotated
    ```

1. Validate and install the sample with Helm. `PROJECT_ID` should be the project ID of the desired project.

    ```bash
    # validate your chart
    helm lint . --set ProjectID=PROJECT_ID --namespace project-annotated

    # check the output of your chart
    helm template . --set ProjectID=PROJECT_ID --namespace project-annotated

    # do a dryrun on your chart and address issues if there are any
    helm install . --dry-run --set ProjectID=PROJECT_ID --namespace project-annotated --generate-name

    # install your chart
    helm install . --set ProjectID=PROJECT_ID --namespace project-annotated --generate-name
    ```

1. _Optionally_ set `Service.Name` in the same manner.

  ```bash
    helm install . --set ProjectID=PROJECT_ID,Service.Name=compute.googleapis.com
  ```

  The package-default value will enable [Firebase](https://firebase.google.com/docs).

1. Check the created helm release to verify the installation:

    ```bash
    helm list
    ```

    Check the status of the applied YAML by specifying namespace:

    ```bash
    kubectl describe gcpservice --namespace project-annotated
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
