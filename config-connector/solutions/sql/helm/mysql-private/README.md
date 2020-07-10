# MySQL Private

==================================================

## NAME

  mysql-private

## SYNOPSIS

  Config Connector compatible YAML files for creating a MySQL instance on a private network.

## CONSUMPTION

  1. Clone GoogleCloudPlatform/cloud-foundation-toolkit repository:

      ```bash
      git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
      ```

  1. Go to the service account folder:

      ```bash
      cd cloud-foundation-toolkit/config-connector/solutions/sql/helm/mysql-private
      ```

## REQUIREMENTS

1. GKE Cluster with Config Connector and [Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity#enable_workload_identity_on_a_new_cluster).
1. [Helm](../../../README.md#helm)
1. A working Config Connector instance.
1. The following APIs enabled in the project where Config Connector is installed:
      -   Cloud SQL Admin API
      -   Compute Engine API
1. The following APIs enabled in the project managed by Config Connector:
      -   Cloud SQL Admin API
      -   Compute Engine API
      -   Service Networking API
      -   Cloud Resource Manager API
1. The "cnrm-system" service account with either both `roles/cloudsql.admin` and
 `roles/compute.networkAdmin` roles or `roles/owner` in the project managed by Config Connector.

## USAGE

All steps are run from the current directory ([config-connector/solutions/sql/helm/mysql-private](.)).

1. Review and update the values in `./sql/values.yaml`.

1. install and check the private network with Helm.

    Due to the bug in Config Connector ([more details](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/148)), the following resources must be in a ready state before the SQLInstance YAML is applied:
    - `ComputeNetwork`
    - `ComputeAddress`
    - `ServiceNetworkingConnection`

    ```bash
    # Do a dryrun on your chart and address issues if there are any
    helm install ./network --dry-run --generate-name

    # install network chart
    helm install ./network

    # Meke sure wait Status of ComputeNetwork,ComputeAddress,ServiceNetworkingConnection is Ready
    kubectl describe ComputeNetwork,ComputeAddress,ServiceNetworkingConnection
    ```

1. install and check the MySQL instance on private network with Helm.

    ```bash
    # validate your chart
    helm lint ./sql --set SQLUser.Name=username,SQLUser.Password=$(echo -n 'your-password' | base64)

    # check the output of your chart
    helm template ./sql --set SQLUser.Name=username,SQLUser.Password=$(echo -n 'your-password' | base64)

    # Do a dryrun on your chart and address issues if there are any
    helm install ./sql --dry-run --set SQLUser.Name=username,SQLUser.Password=$(echo -n 'your-password' | base64) --generate-name

    # install your chart
    helm install ./sql --set SQLUser.Name=username,SQLUser.Password=$(echo -n 'your-password' | base64) --generate-name
    ```

1. _Optionally_ set `Database.Name`, `MySQLInstance.Name`, and `MySQLInstance.Region` in the same
manner. Note that if your instance is deleted the name you used will be
reserved for 7 days. You will need to use a new name in order to re-create the
instance:
    ```bash
    # install your chart with custom changes
    helm install ./sql --set SQLUser.Name=username,SQLUser.Password=$(echo -n 'your-password' | base64),Database.Name=mysql-private-databasename,MySQLInstance.Name=mysql-private-instancename,MySQLInstance.Region=us-west1 --generate-name
    ```

1. Check the created helm release to verify the installation:
    ```bash
    helm list
    ```
    Check the status of the sqlinstances,sqldatabase,sqlusers:
    ```bash
    kubectl describe sqlinstances,sqldatabase,sqlusers
    ```

1. Clean up both installation:

    ```bash
    # list Helm releases to obtain release name
    helm list

    # delete release specifying release name from the previous command output.
    helm delete [release_name]
    ```

## LICENSE

Apache 2.0 - See [LICENSE](/LICENSE) for more information.
