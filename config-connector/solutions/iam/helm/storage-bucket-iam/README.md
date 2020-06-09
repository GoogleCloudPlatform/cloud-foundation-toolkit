# Storage Bucket IAM

==================================================

## NAME

  storage bucket iam

## SYNOPSIS

  Config Connector compatible yaml to enable permissions for a storage bucket.

## CONSUMPTION

  1. Clone GoogleCloudPlatform/cloud-foundation-toolkit repository:

      ```bash
      git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
      ```

  1. Go to the service account folder:

      ```bash
      cd cloud-foundation-toolkit/config-connector/solutions/iam/helm/storage-bucket-iam
      ```

## REQUIREMENTS

1. GKE Cluster with Config Connector and [Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity#enable_workload_identity_on_a_new_cluster).
1. [Helm](../../../README.md#helm)
1. A working Config Connector instance.
1. A storage bucket managed by [IAM](https://cloud.google.com/storage/docs/access-control#using_permissions_with_acls).
1. The "cnrm-system" service account with `roles/storage.admin` in either
  the storage bucket or the project which owns the storage bucket.

## USAGE

All steps are run from the current directory ([config-connector/solutions/iam/helm/storage-bucket-iam](.)).

1. Review and update the values in `./values.yaml`.

1. Validate and install the sample with Helm.

    ```bash
    # validate your chart
    helm lint . --set iamPolicyMember.iamMember=user:name@example.com,StorageBucket.name=your-bucket

    # check the output of your chart
    helm template . --set iamPolicyMember.iamMember=user:name@example.com,StorageBucket.name=your-bucket

    # Do a dryrun on your chart and address issues if there are any
    helm install . --dry-run --set iamPolicyMember.iamMember=user:name@example.com,StorageBucket.name=your-bucket --generate-name

    # install your chart
    helm install . --set iamPolicyMember.iamMember=user:name@example.com,StorageBucket.name=your-bucket --generate-name
    ```

1. _Optionaly_, you can customize optional value role of iam policy member (defaults to `roles/storage.objectViewer`, full list of roles [here](https://cloud.google.com/iam/docs/understanding-roles#storage-roles)):
    ```bash
    # install your chart with a new role
    helm install . --set iamPolicyMember.iamMember=user:name@example.com,StorageBucket.name=your-bucket,iamPolicyMember.role=roles/storage.admin --generate-name
    ```

1. Check the created helm release to verify the installation:
    ```bash
    helm list
    ```
    Check the status of the IAM Policy Member:
    ```bash
    kubectl describe iampolicymember storage-bucket-iam-member
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
