# Wordpress

==================================================

## NAME

  wordpress

## SYNOPSIS

The WordPress application demonstrates how you can configure a WordPress site powered by GCP MySQL database and using Workload Identity for authentication.

## CONSUMPTION

  1. Clone GoogleCloudPlatform/cloud-foundation-toolkit repository
  
      ```bash
      git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
      ```

  1. Go to the wordpress folder:

      ```bash
      cd cloud-foundation-toolkit/config-connector/solutions/apps/helm/wordpress
      ```

## REQUIREMENTS

1. GKE Cluster with Config Connector and [Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity#enable_workload_identity_on_a_new_cluster).
1. [Helm](../../../README.md#helm)
1. Cloud Resource Manager API needs to be enabled on the project to use [ServiceUsage Resource](https://cloud.google.com/config-connector/docs/reference/resources#service). You can enable it by running:

    ```bash
    gcloud services enable cloudresourcemanager.googleapis.com --project [PROJECT_ID]
    ```

## USAGE

All steps are run from this directory.

1. Review and update the values in `./charts/wordpress-gcp/values.yaml`.

    **Note:** Please ensure the value of `database.instanceName` (defaults to `wp-db`) is unique and hasn't been used as an SQL instance name in the last 7 days.
1. Validate and install the sample with Helm.

    ```bash
    # validate your chart
    helm lint ./charts/wordpress-gcp/ --set google.projectId=[PROJECT_ID]

    # check the output of your chart
    helm template ./charts/wordpress-gcp/ --set google.projectId=[PROJECT_ID]

    # install your chart
    helm install ./charts/wordpress-gcp/ --set google.projectId=[PROJECT_ID] --generate-name
    ```

1. The wordpress creation can take up to 10-15 minutes. Throughout the process you can check the status of various components:

    ```bash
    # check the status of sqlinstance
    kubectl describe sqlinstance [VALUE of database.instanceName]
    # check the status of wordpress pod (the output should show that both containers are ready)
    kubectl get pods wordpress-0
    ```

    **Note:** If the pods can't be scheduled because of `Insufficient CPU` issue, please increase the size of nodes in your cluster.
    Once the pods are ready, obtain the external IP address of your WordPress application by checking:

    ```bash
    kubectl get svc wordpress-external
    ```

    Navigate to this address and validate that you see WordPress installation page.

1. Clean up the installation:

    ```bash
    # list Helm releases to obtain release name
    helm list

    # delete release specifying release name from the previous command output. Note that can take a few minutes before all K8s resources are fully deleted.
    helm delete [release_name]
    ```

## LICENSE

Apache 2.0 - See [LICENSE](/LICENSE) for more information.