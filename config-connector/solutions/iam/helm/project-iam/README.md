# Project IAM

==================================================

## NAME

  project-iam

## SYNOPSIS

  Config Connector compatible YAML files to grant a role for a member in a desired project.

## CONSUMPTION

  1. Clone GoogleCloudPlatform/cloud-foundation-toolkit repository:

      ```bash
      git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
      ```

  1. Go to the project-iam folder:

      ```bash
      cd cloud-foundation-toolkit/config-connector/solutions/iam/helm/project-iam
      ```

## REQUIREMENTS

1. GKE Cluster with [Config Connector installed using a GKE Workload Identity](https://cloud.google.com/config-connector/docs/how-to/install-upgrade-uninstall#workload-identity).
1. [Helm](../../../README.md#helm)
1. The "cnrm-system" service account that has the `roles/resourcemanager.projectIamAdmin`
   role in your desired project (it doesn't need to be the project managed by Config Connector).
1. The project managed by Config Connector has Cloud Resource Manager API enabled.

## USAGE

All steps are run from the current directory ([config-connector/solutions/iam/helm/project-iam](.)).

1. Review and update the values in `./values.yaml`.

1. Validate and install the sample with Helm. `PROJECT_ID` should be your desired project ID unique with in GCP.

    ```bash
    # validate your chart
    helm lint . --set iamPolicyMember.iamMember=user:name@example.com,projectID=PROJECT_ID

    # check the output of your chart
    helm template . --set iamPolicyMember.iamMember=user:name@example.com,projectID=PROJECT_ID

    # do a dryrun on your chart and address issues if there are any
    helm install . --dry-run --set iamPolicyMember.iamMember=user:name@example.com,projectID=PROJECT_ID --generate-name

    # install your chart
    helm install . --set iamPolicyMember.iamMember=user:name@example.com,projectID=PROJECT_ID --generate-name
    ```

1. _Optionaly_, you can also change the role (defaults to `roles/logging.viewer`):

    ```bash
    # install your chart with a diffirent role
    helm install . --set iamPolicyMember.iamMember=user:name@example.com,iamPolicyMember.role=roles/logging.admin,projectID=PROJECT_ID --generate-name
    ```

1. Check the created helm release to verify the installation:

    ```bash
    helm list
    ```

    Check the status of the IAM Policy Member:

    ```bash
    kubectl describe iampolicymember project-iam-member
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
