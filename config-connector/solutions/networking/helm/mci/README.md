# Multi-cluster ingress

==================================================

## NAME

  multi-cluster ingress

## SYNOPSIS

Multi-cluster ingress solution connects provisions load balancing across multiple clusters

## CONSUMPTION

  1. Clone GoogleCloudPlatform/cloud-foundation-toolkit repository
  
      ```bash
      git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
      ```

  1. Go to the mci folder:

      ```bash
      cd cloud-foundation-toolkit/config-connector/solutions/networking/helm/mci
      ```

## REQUIREMENTS

1. GKE Cluster with Config Connector
1. [Helm](../../../README.md#helm)
1. Cloud Resource Manager API needs to be enabled on the project to use [ServiceUsage Resource](https://cloud.google.com/config-connector/docs/reference/resources#service). You can enable it by running:

    ```bash
    gcloud services enable cloudresourcemanager.googleapis.com --project [PROJECT_ID]
    ```

## USAGE

All steps are run from this directory.

1. Review and update the values in `./values.yaml`.


1. Configure clusters and load balancing resources with Helm:

    ```bash
    # validate your chart
    helm lint ./lb/ --set projectId=[PROJECT_ID]

    # check the output of your chart
    helm template ./lb/ --set projectId=[PROJECT_ID]

    # install your chart
    helm install ./lb/ --set projectId=[PROJECT_ID] --generate-name
    ```

1. Clean up the installation:

    ```bash
    # list Helm releases to obtain release name
    helm list

    # delete release specifying release name from the previous command output. Note that can take a few minutes before all K8s resources are fully deleted.
    helm delete [release_name]
    ```

## LICENSE

Apache 2.0 - See [LICENSE](/LICENSE) for more information.