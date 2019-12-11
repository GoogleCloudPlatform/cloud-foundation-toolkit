# Stackdriver Notification Channels

This template creates a Stackdriver Notification Channel.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Enable the [Stackdriver Monitoring AP](https://cloud.google.com/monitoring/api/ref_v3/rest/)
- Create a [Stackdriver Workspace](https://cloud.google.com/monitoring/workspaces/)
- Log in to the [Stackdriver Workspace](https://cloud.google.com/monitoring/workspaces/)
  where the metric has to be deployed
- Grant the [monitoring.admin](https://cloud.google.com/monitoring/access-control)
  IAM role to the Deployment Manager service account

## Deployment

### Resources

- [projects.notificationChannels](https://cloud.google.com/monitoring/api/ref_v3/rest/v3/projects.notificationChannels)

### Properties

See the `properties` section in the schema file(s):

- [Stackdriver Notification Channels](stackdriver_notification_channels.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit)

   ```shell
   git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit
   ```

2. Go to the [dm](../../) directory

   ```shell
   cd dm
   ```

3. Copy the example DM config to be used as a model for the deployment,
   in this case [examples/stackdriver\_metric\_descriptor.yaml](examples/stackdriver_notification_channels.yaml)

   ```shell
   cp templates/stackdriver_notification_channels/examples/stackdriver_notification_channels.yaml my_notification_channels.yaml
   ```

4. Change the values in the config file to match your specific GCP setup.
   Refer to the properties in the schema files described above.

   ```shell
   vim my_notification_channels.yaml  # <== Replace <FIXME:...> placeholders if any
   ```

5. Set the project context to use the Stackdriver Workspace project

   ```shell
   gcloud config set project <STACKDRIVER_WORKSPACE_PROJECT_ID>
   ```

6. Create your deployment as described below, replacing `<YOUR_DEPLOYMENT_NAME>`
   with your with your own deployment name

   ```shell
   gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
       --config my_notification_channels.yaml
   ```

7. In case you need to delete your deployment

   ```shell
   gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
   ```

## Examples

- [Stackdriver Notification Channels](examples/stackdriver_notification_channels.yaml)
