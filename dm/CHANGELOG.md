# Cloud Foundation Toolkit Change Log

All notable changes to this project will be documented in this file.

## CFT Templates

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
