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
1. [Helm](https://helm.sh/docs/using_helm/)

## USAGE

All steps are run from this directory.

1. Review and update the values in `./charts/wordpress-gcp/values.yaml` .
1. Validate and install the sample with Helm

    ```bash
    # validate your chart
    helm lint ./charts/wordpress-gcp/ --set google.projectId=[PROJECT_ID]

    # check the output of your chart
    helm template ./charts/wordpress-gcp/ --set google.projectId=[PROJECT_ID]

    # install your chart
    helm install ./charts/wordpress-gcp/ --set google.projectId=[PROJECT_ID] --generate-name
    ```

1. Check the status of your database by running `kubectl describe sqlinstance wp-db`. Once the database is created, obtain the external IP address of your WordPress application by checking `kubectl get svc wordpress-external`. Navigate to this address and validate that you see WordPress installation page.

1. Clean up the installation:

    ```bash
    # list Helm releases
    helm list

    # delete release
    helm delete [release_name]
    ```

## LICENSE

Apache 2.0 - See [LICENSE](/LICENSE) for more information.