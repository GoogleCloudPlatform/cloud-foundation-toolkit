# Cloud Foundation Toolkit Change Log

All notable changes to this project will be documented in this file.

## CFT Templates

### 11.12.2019

- Updated logging sink configuration to export entries to a desired destination in external project [#77](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/496)
- Added Stackdriver Notification Channels template [#432](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/485)

### 09.12.2019 Ho-ho-ho

- SSL-Certificate template supports beta features (managed certificate). This update is backwards compatible. [#505](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/505)

### 09.12.2019

- Added 'resource_policy' DM template [#497](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/497)

### 05.12.2019

- Updated internal LB and external LB templates according to backend_service.py.schema change [#476](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/476)

### 02.12.2019

- IAM Member template support bindings on types which implement `gcp-types/storage-v1:virtual.buckets.iamMemberBinding` like syntax. ( currently storage-v1.)

### 25.11.2019

- In `cloud_sql.py`, added support for PostgreSQL 11 & fixed `ipAddress` output [#477](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/477)

### 22.11.2019

- Fixed sharedVPC for GKE use case behaviour in 'project' DM template [#469](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/469)

### 22.11.2019

- Cloud Build Trigger support for Github [#470](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/463)
- CI example for Github triggering Cloudbuild for PRs

### 21.11.2019

- Added support for unified Stackdriver Kubernetes Engine Monitoring [#463](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/463)
- Add explicit dependencies to the 'iam_member' DM template to avoid fail, in case of a large amount of bindings (30+) [#443](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/443)

### 18.11.2019

- The [GCS Bucket template](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/tree/ocsig-patch-storage1/dm/templates/gcs_bucket) supports gcp-types/storage-v1:virtual.buckets.iamMemberBinding. [#453](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/453)

### 31.10.2019

- New helper template to use firewall with Google important IP ranges, stored in YAML [#370](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/370)

### 29.10.2019

- Bigquery template schema supports single region as location

### 25.10.2019

- Fixed examples for ELB to pass the new (strict) schema validation [#392](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/392)
- Added support for GKE on SharedVPC for the Project Factory and better visibility for the GKE Service Accounts [#385](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/385)
- Fixed IAM member binding schema to truly support the lack of project property. [#384](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/384)
- The CFT CLI (go version) support complex cross deployment references [#359](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/359)
- Added a Docker container for running bats tests on your local source code for template developers [#355](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/355)
- New exapmle: [Instance with a private IP](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/dm/templates/instance/examples/instance_private.yaml) [#346](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/346)
- New template: [Cloud Filestore](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/tree/master/dm/templates/cloud_filestore) [#348](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/348)

### 17.09.2019

- New template: [Unamanged Instance Group](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/tree/master/dm/templates/unmanaged_instance_group)
- CFT Instance (DM) template support sourceInstanceTemplate property instead of properties of the instance [#330](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/330)
- New examples for CloudSQL with [private IP](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/dm/templates/cloud_sql/examples/cloud_sql_private_network.yaml)
- The [pubsub template](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/tree/master/dm/templates/pubsub) supports subscription expiration
- Github PRs are now automatically [triggering](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/pull/312) CloudBuild based tests

### 03.09.2019

 **cft-dm-dev branch merged to master**
 
#### Major updates which may break your current deployments:
- Switching legacy types to 'gcp-types'
- Switching to the latest API where it is possible
- Adding unique hash to resource names in case of iteration
  - This is fixing many issues with iterations but it is a breaking change

#### Non-breaking changes:
- Adding versioning for the templates
- Adding new properties for templates, schemas where it's applicable
- Adding support for Labels for every resource where it's possible
- Adding cross project resource creation where it is possible
- Locking down schemas:
  - Tight check on invalid properties to catch typos instead of ignoring them
  - Tight check on combination of properties. ( For example a project can't be
  a host and a guest (VPC) project at the same time.)
  
#### CI improvements:
- Our CI environment is running tests on the current master and dev branch
  - Running schema validation checks on the example yamls where it's applicable
  - Running integration tests on all the templates
- CloudBuild containers and jobs running the tests in a test organization
  - Currently working on local container based testing with local source code

### 23.08.2019

- Adding container images for test automation
- Finalizing 'cft-dm-dev' branch for merge to master

### 21.03.2019

- *Templates/iam_member*: The template is now using virtual.projects.iamMemberBinding which is and advanced
endpoint to manage IAM membership. This is a fix for concurrent IAM changes error.
   - This change should be 100% backwards compatible
   - This template should solve concurrency error with built in retries
 - *Templates/project*: This template had  concurrent IAM changes error. This update utilizes the iam_member 
 CFT template, which is referenced in the project.py.schema file. No more concurancy error!

### 20.03.2019

 - *Example Solutions*: The first exmaple demonstrate how to use Wrapper templates.
   - *Specific wrapper* template to modify the behaviour of an external template such as a CFT template
   - *Generic Wrapper* template to inject configuration for every template regardless of it's behaviour.

### 19.03.2019

 - *CloudDNS*: Changed CloudDNS Record set from actions to use gcp-types which gives native support for the API.

## CFT CLI

### 0.0.4

- Feature: Cross deployment refference support output lookup of complex DM resources 

### 0.0.3

- Feature: Cross deployment refference support complex outputs such as hashmap, list and their combination. 

### 0.0.2 

- Initial version
