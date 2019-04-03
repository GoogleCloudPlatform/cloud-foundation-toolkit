# Stackdriver Metric Descriptor

This template creates a Stackdriver Metric Descriptor.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Create a [Stackdriver Workspace](https://cloud.google.com/monitoring/workspaces/)
- Log in to the [Stackdriver Workspace](https://cloud.google.com/monitoring/workspaces/)
  where the metric has to be deployed
- Grant the [monitoring.admin](https://cloud.google.com/monitoring/access-control)
  IAM role to the Deployment Manager service account

## Deployment

### Resources

- [projects.metricDescriptors](https://cloud.google.com/monitoring/api/ref_v3/rest/v3/projects.metricDescriptors)
- [GCP metric list](https://cloud.google.com/monitoring/api/metrics_gcp)
- [AWS metric list](https://cloud.google.com/monitoring/api/metrics_aws)
- [Stackdriver Agent metric list](https://cloud.google.com/monitoring/api/metrics_agent)
- [External metric list](https://cloud.google.com/monitoring/api/metrics_other)

### Properties

See the `properties` section in the schema file(s):

- [Stackdriver Metric Descriptor](stackdriver_metric_descriptor.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples)

   ```shell
   git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
   ```

2. Go to the [community/cloud-foundation](../../) directory

   ```shell
   cd community/cloud-foundation
   ```

3. Copy the example DM config to be used as a model for the deployment,
   in this case [examples/stackdriver\_metric\_descriptor.yaml](examples/stackdriver_metric_descriptor.yaml)

   ```shell
   cp templates/stackdriver_metric_descriptor/examples/stackdriver_metric_descriptor.yaml my_metric_descriptor.yaml
   ```

4. Change the values in the config file to match your specific GCP setup.
   Refer to the properties in the schema files described above.

   ```shell
   vim my_metric_descriptor.yaml  # <== Replace <FIXME:...> placeholders if any
   ```

5. Set the project context to use the Stackdriver Workspace project

   ```shell
   gcloud config set project <STACKDRIVER_WORKSPACE_PROJECT_ID>
   ```

6. Create your deployment as described below, replacing `<YOUR_DEPLOYMENT_NAME>`
   with your with your own deployment name

   ```shell
   gcloud deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
       --config my_metric_descriptor.yaml
   ```

7. In case you need to delete your deployment

   ```shell
   gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
   ```

## Examples

- [Stackdriver Metric Descriptor](examples/stackdriver_metric_descriptor.yaml)
