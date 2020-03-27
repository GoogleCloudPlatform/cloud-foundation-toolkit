Pubsub Topic
==================================================

# NAME

  pubsub-topic

# SYNOPSIS

  Config Connector compatible YAML files to grant a role to a particular IAM member for a PubSub topic.

# CONSUMPTION

  Using [kpt](https://googlecontainertools.github.io/kpt/):

  ```
  SRC=https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
  kpt pkg get $SRC/config-connector/solutions/iam/kpt/pubsub-topic pubsub-topic
  ```

# REQUIREMENTS

  A cluster with Config Connector installed managing a GCP project with the
PubSub API enabled.

# USAGE

  Replace `${PROJECT_ID?}` with your desired project ID:
  ```
  kpt cfg set . project-id your-project-id
  ```

  Replace `${IAM_MEMBER?}` with the GCP identity to grant access to:
  ```
  kpt cfg set . iam-member user:name@example.com
  ```

  Optionally set the name of the PubSub topic (defaults to `allowed-topic`) and
the role to grant (defaults to `roles/pubsub.editor`, full list of roles
[here](https://cloud.google.com/iam/docs/understanding-roles#pub-sub-roles))
  ```
  kpt cfg set . topic-name your-topic
  kpt cfg set . role roles/pubsub.viewer
  ```

  Once the values are satisfactory, simply apply the YAMLs:
  ```
  kubectl apply -f .
  ```

  Note: This will create the topic if it does not exist.


# LICENSE

  Apache 2.0 - See [LICENSE](/LICENSE) for more information.
