# PostgreSQL High Availability

==================================================

## NAME

  postgres-ha

## SYNOPSIS

  Config Connector compatible yaml files to configure a high availability PostgreSQL cluster.

## CONSUMPTION

  1. Clone GoogleCloudPlatform/cloud-foundation-toolkit repository:

      ```bash
      git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
      ```

  1. Go to the service account folder:

      ```bash
      cd cloud-foundation-toolkit/config-connector/solutions/sql/helm/postgres-ha
      ```

## REQUIREMENTS

1. GKE Cluster with [Config Connector installed using a GKE Workload Identity](https://cloud.google.com/config-connector/docs/how-to/install-upgrade-uninstall#workload-identity).
1. [Helm](../../../README.md#helm).
1. Cloud SQL Admin API enabled in the project where Config Connector is installed.
1. Cloud SQL Admin API enabled in the project managed by Config Connector if it is installed in a different project.
1. The "cnrm-system" service account with either `roles/cloudsql.admin` or `roles/owner` in the project managed by Config Connector.

## USAGE

All steps are run from the current directory ([config-connector/solutions/sql/helm/postgres-ha](.)).

1. Review and update the values in `./values.yaml`.

1. Configure a high availability PostgreSQL cluster with Helm.

    ```bash
    # validate your chart
    helm lint . --set User1.Name=first-username,User2.Name=second-username,User3.Name=third-username,\
    User1.Password=$(echo -n 'first-password' | base64),User2.Password=$(echo -n 'second-password' | base64),\
    User3.Password=$(echo -n 'third-password' | base64)

    # check the output of your chart
    helm template . --set User1.Name=first-username,User2.Name=second-username,User3.Name=third-username,\
    User1.Password=$(echo -n 'first-password' | base64),User2.Password=$(echo -n 'second-password' | base64),\
    User3.Password=$(echo -n 'third-password' | base64)

    # do a dryrun on your chart
    helm install . --dry-run --set User1.Name=first-username,User2.Name=second-username,User3.Name=third-username,\
    User1.Password=$(echo -n 'first-password' | base64),User2.Password=$(echo -n 'second-password' | base64),\
    User3.Password=$(echo -n 'third-password' | base64) --generate-name

    # install your chart
    helm install . --set User1.Name=first-username,User2.Name=second-username,User3.Name=third-username,\
    User1.Password=$(echo -n 'first-password' | base64),User2.Password=$(echo -n 'second-password' | base64),\
    User3.Password=$(echo -n 'third-password' | base64) --generate-name
    ```

1. _Optionally_ here the list of things you can customize.

    |       NAME        |      DEFAULT VALUE     |          DESCRIPTION           |
    |-------------------|------------------------|--------------------------------|
    | Database1.Name   | postgres-ha-database-1 | name of first SQL database     |
    | Database2.Name   | postgres-ha-database-2 | name of second SQL database    |
    | external-ip-range | 192.10.10.10/32        | ip range to allow to connect   |
    | PostgreSQLInstance.Name     | postgres-ha-solution   | name of main SQL instance      |
    | PostgreSQLInstance.Region | us-central1            | region of SQL instance         |
    | PostgreSQLInstance.Zone              | us-central1-c          | zone of main instance          |
    | PostgreSQLInstance.Replica1.Zone    | us-central1-a          | zone of first replica instance |
    | PostgreSQLInstance.Replica2.Zone    | us-central1-b          | zone of second replica instance|
    | PostgreSQLInstance.Replica3.Zone    | us-central1-c          | zone of third replica instance |

    **Note:** If your SQL Instance is deleted, the name you used will be reserved
for **7 days**. In order to re-apply this solution, you need to run the following command to update the value of PostgreSQLInstance.Name to "new-instance-name".

    ```bash
    helm install . --set User1.Name=first-username,User2.Name=second-username,User3.Name=third-username,\
    User1.Password=$(echo -n 'first-password' | base64),User2.Password=$(echo -n 'second-password' | base64),\
    User3.Password=$(echo -n 'third-password' | base64), PostgreSQLInstance.Name=new-instance-name --generate-name
    ```

1. Check the created helm release to verify the installation:
    ```bash
    helm list
    ```
    Check the status of the sqlinstances,sqldatabase,sqlusers:
    ```bash
    kubectl describe sqldatabase,sqlinstance,sqluser,secret
    ```

1. Clean up installation:

    ```bash
    # list Helm releases to obtain release names
    helm list

    # delete release specifying release name from the previous command output.
    helm delete [release_name]
    ```

## LICENSE

Apache 2.0 - See [LICENSE](/LICENSE) for more information.
