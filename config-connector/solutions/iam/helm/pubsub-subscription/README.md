# Pub/Sub Subscription

==================================================

## NAME

  pubsub-subscription

## SYNOPSIS

  Config Connector compatible YAML files to create a Pub/Sub subscription and Pub/Sub topic if it doesn't exist in your desired project and grant a specific member a role (default to roles/pubsub.viewer) for accessing the Pub/Sub subscription that just created

## CONSUMPTION

  1. Clone GoogleCloudPlatform/cloud-foundation-toolkit repository:

      ```bash
      git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
      ```

  1. Go to the service account folder:

      ```bash
      cd cloud-foundation-toolkit/config-connector/solutions/iam/helm/pubsub-subscription
      ```

## REQUIREMENTS

1. GKE Cluster with Config Connector and [Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity#enable_workload_identity_on_a_new_cluster).
1. [Helm](../../../README.md#helm)
1. The "cnrm-system" service account assigned with either `roles/pubsub.viewer` or `roles/owner` in the project managed by Config Connector
1. Cloud Pub/Sub API enabled in the project where Config Connector is installed
1. Cloud Pub/Sub API enabled in the project managed by Config Connector if it is a different project

## USAGE

All steps are run from the current directory ([config-connector/solutions/iam/helm/pubsub-subscription](.)).

1. Review and update the values in `./values.yaml`.

1. Validate and install the sample with Helm.

    ```bash
    # validate your chart
    helm lint . --set iamPolicyMember.iamMember=user:name@example.com

    # check the output of your chart
    helm template . --set iamPolicyMember.iamMember=user:name@example.com

    # do a dryrun on your chart and address issues if there are any
    helm install . --dry-run --set iamPolicyMember.iamMember=user:name@example.com --generate-name

    # install your chart
    helm install . --set iamPolicyMember.iamMember=user:name@example.com --generate-name
    ```

1. _Optionaly_, you can set the name of the pubsub subscription (defaults to `allowed-subscription`), set the name of the PubSub topic (defaults to `allowed-topic`) and the role to grant (defaults to `roles/pubsub.viewer`, full list of roles [here](https://cloud.google.com/iam/docs/understanding-roles#pub-sub-roles)) by explictly setting them when installing the solution:

    ```bash
    # install your chart with a difirent pubsub subscription and pubsub topic
    helm install . --set PubSubTopic.name=your-topic-name,PubSubSubscription.name=your-subscription-name,iamPolicyMember.iamMember=user:name@example.com --generate-name
    ```
    Or,
    ```bash
    # install your chart with a new role
    helm install . --set iamPolicyMember.role=roles/pubsub.editor,iamPolicyMember.iamMember=user:name@example.com --generate-name
    ```
    Or set there in one command.

1. Check the created helm release to verify the installation:
    ```bash
    helm list
    ```
    Check the status of the pubsub subscription and topic resource by running:

    Note: By default value of pubsub topic name is ```allowed-topic``` and pubsub subscription name is ```allowed-subscription```

    ```bash
    kubectl describe pubsubsubscriptions [pubsub subscription name]

    kubectl describe pubsubtopic [pubsub topic name]
    ```
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
