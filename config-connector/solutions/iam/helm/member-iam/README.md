# Member IAM

==================================================

## NAME

  member-iam

## SYNOPSIS

  Config Connector compatible YAML files to create a service account in your desired project, and grant it a specific role (defaults to `compute.networkAdmin`) in the project.

## CONSUMPTION

  1. Clone GoogleCloudPlatform/cloud-foundation-toolkit repository:

      ```bash
      git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
      ```

  1. Go to the member-iam folder:

      ```bash
      cd cloud-foundation-toolkit/config-connector/solutions/iam/helm/member-iam
      ```

## REQUIREMENTS

1. GKE Cluster with [Config Connector installed using a GKE Workload Identity](https://cloud.google.com/config-connector/docs/how-to/install-upgrade-uninstall#workload-identity).
1. [Helm](../../../README.md#helm)
1. The "cnrm-system" service account assigned with `roles/resourcemanager.projectIamAdmin` and `roles/iam.serviceAccountAdmin` or `roles/owner`
    role in your desired project (it doesn't need to be the project managed by Config Connector)
1. Cloud Resource Manager API enabled in the project where Config Connector is installed

## USAGE

All steps are run from the current directory ([config-connector/solutions/iam/helm/member-iam](.)).

1. Review and update the values in `./values.yaml`.

1. Create a namespace. If you want to use your existing namespace skip this step and use your own namespace name instead of `member-iam-solution` in all other steps.

    ```bash
    kubectl create namespace member-iam-solution
    ```

1. Validate and install the sample with Helm. `PROJECT_ID` should be your desired project ID unique with in GCP.

    ```bash
    # validate your chart
    helm lint . --set projectID=PROJECT_ID --namespace member-iam-solution

    # check the output of your chart
    helm template . --set projectID=PROJECT_ID --namespace member-iam-solution

    # do a dryrun on your chart and address issues if there are any
    helm install . --dry-run --set projectID=PROJECT_ID --namespace member-iam-solution --generate-name

    # install your chart
    helm install . --set projectID=PROJECT_ID --namespace member-iam-solution --generate-name
    ```

1. _Optionally_, you can also change the service account name `iamPolicyMember.iamMember` (defaults to `member-iam-test`) and role `iamPolicyMember.role` (defaults to `roles/compute.networkAdmin`)
  (you can find all the predefined GCP IAM roles [here](https://cloud.google.com/iam/docs/understanding-roles#predefined_roles)):

    ```bash
    # install your chart with a diffirent service account name
    helm install . --set projectID=PROJECT_ID,iamPolicyMember.iamMember=service-account-name --namespace member-iam-solution --generate-name
    ```
    Or,
    ```bash
    # install your chart with a diffirent role
    helm install . --set projectID=PROJECT_ID,iamPolicyMember.role=roles/compute.networkUser --namespace member-iam-solution --generate-name
    ```
    Or set both in one command.

1. Check the created helm release to verify the installation:

    ```bash
    helm list
    ```

    Check the status of the IAM Service Account:

    ```bash
    kubectl describe iamserviceaccount [Value of iamPolicyMember.iamMember]
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
