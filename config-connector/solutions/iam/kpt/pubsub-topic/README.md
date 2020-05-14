Pub/Sub Topic
==================================================
# NAME
  pubsub-topic
# SYNOPSIS
  Config Connector compatible YAML files to grant a role to a particular IAM member for a PubSub topic.
# CONSUMPTION
  Download the package using [kpt](https://googlecontainertools.github.io/kpt/).
  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/config-connector/solutions/iam/kpt/pubsub-topic pubsub-topic
  ```
# REQUIREMENTS
  -   A working Config Connector instance using the "cnrm-system" service
      account with either `roles/pubsub.admin` or `roles/owner` in the project
      managed by Config Connector
  -   Cloud Pub/Sub API enabled in the project where Config Connector is
      installed
  -   Cloud Pub/Sub API enabled in the project managed by Config Connector if it
      is a different project
# SETTERS
|    NAME    |        VALUE        |     SET BY      |         DESCRIPTION          | COUNT |
|------------|---------------------|-----------------|------------------------------|-------|
| iam-member | ${IAM_MEMBER?}      | PLACEHOLDER     | identity to grant privileges | 1     |
| role       | roles/pubsub.editor | package-default | IAM role to grant            | 1     |
| topic-name | allowed-topic       | package-default | name of PubSub topic         | 2     |
# USAGE

  Replace `${IAM_MEMBER?}` with the GCP identity to grant access to.
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
  Once the values are satisfactory, apply the YAMLs.
  ```
  kubectl apply -f .
  ```
  Note: This will create the topic if it does not exist.
# LICENSE
  Apache 2.0 - See [LICENSE](/LICENSE) for more information.
