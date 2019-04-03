# Cloud Tasks

This set of two templates creates a Cloud Task and a Cloud Task Queue.

## Prerequisites

- Install [gcloud](https://cloud.google.com/sdk)
- Install gcloud **beta** components:

  ```(shell)
  gcloud components update
  gcloud components install beta
  ```

- Create a [GCP project, set up billing, enable requisite APIs](../project/README.md)
- Enable the [Cloud Tasks API](https://console.cloud.google.com/apis/library/cloudtasks.googleapis.com)
  from the Google Cloud console
- Grant the [appengine.applications.get](https://cloud.google.com/appengine/docs/admin-api/access-control)
  IAM permission to the Deployment Manager service account
- NOTE: Cloud Tasks requires an App Engine application. To run the integration tests
  please ensure that an App Engine application exists. An App Engine app can be created using the [App Engine Template](../app_engine) or by running the [App Engine Template integration tests](../app_engine/tests/integration)

## Deployment

### Resources

- [projects.locations.queues](https://cloud.google.com/tasks/docs/reference/rest/v2beta3/projects.locations.queues)
- [projects.locations.queues.tasks](https://cloud.google.com/tasks/docs/reference/rest/v2beta3/projects.locations.queues.tasks)
- [Task Queues](https://cloud.google.com/appengine/docs/standard/python/taskqueue/)
- [CloudTasks v2beta3 Descriptor URL](https://cloudtasks.googleapis.com/$discovery/rest?version=v2beta3)

### Properties

See the `properties` section in the schema file(s):

- [CloudTasks Queue schema](queue.py.schema)
- [CloudTasks Task schema](task.py.schema)

### Usage

1. Clone the [Deployment Manager samples repository](https://github.com/GoogleCloudPlatform/deploymentmanager-samples)

   ```(shell)
   git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
   ```

2. Go to the [community/cloud-foundation](../../) directory

   ```(shell)
   cd community/cloud-foundation
   ```

3. Create a custom type-provider named `cloudtasks`

   ```(shell)
   cp templates/cloud_tasks/examples/create_typeprovider.sh .
   chmod u+x create_typeprovider.sh
   ./create_typeprovider.sh
   ```

4. Copy the example DM config to be used as a model for the deployment. In this case, [examples/cloud\_tasks\_queue.yaml](examples/cloud_tasks_queue.yaml)

   ```(shell)
   cp templates/cloud_tasks/examples/cloud_tasks_queue.yaml my_cloud_tasks_queue.yaml
   ```

5. Change the values in the config file to match your specific GCP setup.
   Refer to the properties in the schema files described above.

   ```(shell)
   vim my_cloud_tasks_queue.yaml
   ```

6. Create your deployment as described below, replacing `<YOUR_DEPLOYMENT_NAME>`
   with your with your own deployment name

   ```(shell)
   gcloud beta deployment-manager deployments create <YOUR_DEPLOYMENT_NAME> \
       --config my_cloud_tasks_queue.yaml
   ```

7. In case you need to delete your deployment

   ```(shell)
   gcloud deployment-manager deployments delete <YOUR_DEPLOYMENT_NAME>
   ```

8. To delete the custom `cloudtasks` type-provider

   ```(shell)
   cp templates/cloud_tasks/examples/delete_typeprovider.sh .
   chmod u+x delete_typeprovider.sh
   ./delete_typeprovider.sh
   ```

## Examples

- [CloudTasks Queue](examples/cloud_tasks_queue.yaml)
- [CloudTasks Task](examples/cloud_tasks_task.yaml)
