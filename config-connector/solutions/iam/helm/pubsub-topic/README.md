# Pub/Sub Topic

==================================================

## NAME

  pubsub-topic

## SYNOPSIS

  
  Config Connector compatible YAML files to create a Pub/Sub topic in your desired project, and grant a specific member a role (default to roles/pubsub.editor) for accessing the Pub/Sub topic that just created

## CONSUMPTION

  1. Clone GoogleCloudPlatform/cloud-foundation-toolkit repository:

      ```bash
      git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
      ```

  1. Go to the service account folder:

      ```bash
      cd cloud-foundation-toolkit/config-connector/solutions/iam/helm/pubsub-topic
      ```

## REQUIREMENTS

1. GKE Cluster with Config Connector and [Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity#enable_workload_identity_on_a_new_cluster).
1. [Helm](../../../README.md#helm)

## USAGE

All steps are run from the current directory ([config-connector/solutions/iam/helm/pubsub-topic](.)).

1. Review and update the values in `./values.yaml`.

1. Validate and install the sample with Helm.

    ```bash
    # validate your chart
    helm lint . --set iamPolicyMember.iamMember=user:name@example.com

    # check the output of your chart
    helm template . --set iamPolicyMember.iamMember=user:name@example.com

    # Do a dryrun on your chart and address issues if there are any
    helm install . --dry-run --set iamPolicyMember.iamMember=user:name@example.com --generate-name

    # install your chart
    helm install . --set iamPolicyMember.iamMember=user:name@example.com --generate-name
    ```

1. _Optionaly_, you can customize optional values by explictly setting them when installing the solution:
    ```bash
    # install your chart with a difirent name of the PubSub topic
    helm install . --set PubSubTopic.name=your-topic,iamPolicyMember.iamMember=user:name@example.com --generate-name
    ```
    Or,
    ```bash
    # install your chart with a new role
    helm install . --set iamPolicyMember.role=roles/pubsub.viewer,iamPolicyMember.iamMember=user:name@example.com --generate-name
    ```
    Or set them both in one command.
    
    Note: These will create the topic if it does not exist.

1. Check the created helm release to verify the installation:
    ```bash
    helm list
    ```
    Check the status of the pub/sub topic resource by running:
    ```bash
    kubectl describe pubsubtopic [topic name]
    ```
    Note: By default value of Pub/Sub topic name is ```allowed-topic```
    Check the status of the IAM Policy Member:
    ```bash
    kubectl describe iampolicymember topic-iam-member
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
