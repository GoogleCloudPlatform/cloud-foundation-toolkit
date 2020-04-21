Pub/Sub Subscription
==================================================
# NAME
  pubsub-subscription
# SYNOPSIS
  This package creates a pubsub subscription and configures permissions for it by creating an IAMPolicyMember resource.
# CONSUMPTION
  Download the package using [kpt](https://googlecontainertools.github.io/kpt/):
  ```
  kpt pkg get https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git/config-connector/solutions/iam/kpt/pubsub-subscription pubsub-subscription
  ```
# REQUIREMENTS
  A Config Connector installation managing a GCP project with Pub/Sub API enabled. 
# SETTERS
|       NAME        |        VALUE         |     SET BY      |         DESCRIPTION         | COUNT |
|-------------------|----------------------|-----------------|-----------------------------|-------|
| iam-member        | ${IAM_MEMBER?}       | PLACEHOLDER     | IAM member to grant role    | 1     |
| role              | roles/pubsub.viewer  | package-default | IAM role to grant           | 1     |
| subscription-name | allowed-subscription | package-default | name of PubSub subscription | 2     |
| topic-name        | allowed-topic        | package-default | name of PubSub topic        | 2     |
# USAGE
  Set the `iam-member` to grant a role to.
  ```
  kpt cfg set . iam-member user:name@example.com
  ```
  _Optionally_ set the `role` to grant. The default role is `roles/pubsub.viewer`.
  ```
  kpt cfg set . role roles/pubsub.editor
  ```
  _Optionally_ set `topic-name` and `subscription-name` in the same manner. Defaults are `allowed-topic` and `allowed-subscription`.

  Once the configuration is satisfactory, apply:
  ```
  kubectl apply -f .
  ```
# LICENSE
Apache 2.0 - See [LICENSE](/LICENSE) for more information.

