# Folder IAM

==================================================

## NAME

  folder-iam

## SYNOPSIS

  Config Connector compatible YAML files to grant a specific member a role (default to `roles/resourcemanager.folderEditor`) to an existing folder.

## CONSUMPTION

  1. Clone GoogleCloudPlatform/cloud-foundation-toolkit repository:
  
      ```bash
      git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
      ```

  1. Go to the folder-iam folder:

      ```bash
      cd cloud-foundation-toolkit/config-connector/solutions/iam/helm/folder-iam
      ```

## REQUIREMENTS

1. GKE Cluster with Config Connector and [Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity#enable_workload_identity_on_a_new_cluster).

1. A working Config Connector cluster using the "cnrm-system" service account with _minimally_ the permissions given by the following role on the desired folder:
    - `roles/resourcemanager.folderIamAdmin`

1. Install [Helm](../../../README.md#helm)

## USAGE

All steps are running from the current directory ([config-connector/solutions/iam/helm/folder-iam](.)).

1. Review and update the values in `./values.yaml`.

1. Validate and install the sample with Helm.

    ```bash
    # validate your chart
    helm lint . --set iamPolicyMember.iamMember=user:name@example.com,folderID=VALUE

    # check the output of your chart
    helm template . --set iamPolicyMember.iamMember=user:name@example.com,folderID=VALUE

    # do a dryrun on your chart and address issues if there are any
    helm install . --dry-run --set iamPolicyMember.iamMember=user:name@example.com,folderID=VALUE --generate-name

    # install your chart
    helm install . --set iamPolicyMember.iamMember=user:name@example.com,folderID=VALUE --generate-name
    ```

1. _Optionaly_, you can also change the role granted to the member. (you can find all of the folder related IAM roles
  [here](https://cloud.google.com/iam/docs/understanding-roles#resource-manager-roles)):
    ```bash
    # install your chart with a new IAM role.
    helm install . --set iamPolicyMember.role=roles/iam.serviceAccountTokenCreator,iamPolicyMember.iamMember=user:name@example.com,folderID=VALUE --generate-name
    ```

1. Check the created helm release to verify the installation:
    ```bash
    helm list
    ```

    Check the status of the IAM Policy Member:
    ```bash
    kubectl describe iampolicymember iampolicymember-folder-iam
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
