# Service Account

==================================================

## NAME

  service account

## SYNOPSIS

  Config Connector compatible YAML files to create a service account in your desired project, and grant a specific member a role (default to roles/iam.serviceAccountKeyAdmin) for accessing the service account that just created.

## CONSUMPTION

  1. Clone GoogleCloudPlatform/cloud-foundation-toolkit repository
  
      ```bash
      git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
      ```

  1. Go to the service account folder:

      ```bash
      cd cloud-foundation-toolkit/config-connector/solutions/iam/helm/service-account
      ```

## REQUIREMENTS

1. GKE Cluster with Config Connector and [Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity#enable_workload_identity_on_a_new_cluster).
1. [Helm](../../../README.md#helm)

## USAGE

All steps are run from this directory.

1. Review and update the values in `./values.yaml`.

1. Validate and install the sample with Helm.

    ```bash
    # validate your chart
    helm lint . --set iamPolicyMember.iamMember=user:name@example.com

    # check the output of your chart
    helm template . --set iamPolicyMember.iamMember=user:name@example.com

    # Dryrun your chart and address issues if there are any
    helm install . --dry-run --set iamPolicyMember.iamMember=user:name@example.com --generate-name

    # install your chart
    helm install . --set iamPolicyMember.iamMember=user:name@example.com --generate-name
    ```

1. _Optionaly_, you can costomize optional values by explictly set them when installing the solution:
    ```bash
    # install your chart with a new service account name
    helm install . --set serviceAccount.name=new-service-account --set iamPolicyMember.iamMember=user:name@example.com --generate-name
    ```  
    Or,
    ```bash
    # install your chart with a new service account name
    helm install . --set iamPolicyMember.role=roles/iam.serviceAccountTokenCreator --set iamPolicyMember.iamMember=user:name@example.com --generate-name
    ```
    Or set them both in one command.

1. Clean up the installation:

    ```bash
    # list Helm releases to obtain release name
    helm list

    # delete release specifying release name from the previous command output.
    helm delete [release_name]
    ```

## LICENSE

Apache 2.0 - See [LICENSE](/LICENSE) for more information.