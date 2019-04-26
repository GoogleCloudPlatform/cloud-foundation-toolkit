# Cloud Foundation Infra

IaC for launching CICD infrastructure using [Concourse CI](https://concourse-ci.org/).

## Infrastructure

Concourse runs on [GKE](https://cloud.google.com/kubernetes-engine/docs/) and [Cloud SQL for Postgres](https://cloud.google.com/sql/docs/postgres/).

Concourse is deployed with [Helm](https://helm.sh/) and the official [Concourse chart](https://github.com/helm/charts/tree/master/stable/concourse).

This is all managed with [Terraform](https://www.terraform.io/). All Terraform configurations are located in the [terraform](./terraform) directory of this repository. At this time, `terraform` is intended to be executed locally on an adminstrator workstation.

### Workspaces

Terraform configurations were written to accommodate multiple [workspaces](https://www.terraform.io/docs/state/workspaces.html). Currently, only a workspace called, "primary", is in use - as it is intended to be the primary CICD infrastructure for repositories under the [terraform-google-modules](https://github.com/terraform-google-modules) Github organization.

### Cloud SQL for PostgreSQL

The [terraform/postgres](terraform/postgres) configuration uses the [private IP](https://cloud.google.com/sql/docs/postgres/private-ip) connection strategy. This functionality is currently in beta, and requires that the [google-beta provider](https://github.com/terraform-providers/terraform-provider-google-beta) be compiled and installed. This was chosen over using the Cloud SQL Proxy Docker image which is documented alongside private ip in "[Connecting from Google Kubernetes Engine](https://cloud.google.com/sql/docs/postgres/connect-kubernetes-engine)".

Cloud SQL was chosen over running Postgres with GKE with the intent of reducing operational complexity. Though both methods are complex, I'd argue that Cloud SQL's builtin data management features like automated backups make it worthwhile. This could also be considered more secure as it allows us to keep all Postgres secrets out of Helm - as mentioned in Concourse's [chart docs](https://github.com/helm/charts/tree/master/stable/concourse#postgresql).

## Concourse Configurations

There three different locations where Concourse configurations are found:

1. [terraform/applications/concourse.tf](./terraform/applications/concourse.tf) - Terraform configuration for Helm release on GKE.
1. [terraform/test_fixtures](./terraform/test_fixtures) - test fixtures for each module such as GCP projects, service accounts, and K8s secrets for the pipelines to consume.
1. [concourse](./concourse) - configurations for pipelines and image builds. These are not managed with Terraform.

### Test Fixtures

All modules have certain GCP objects that are expected to be pre-provisioned. All projects should have the following dedicated resources:

* GCP project
* GCP service account
* A kubernetes secret containing service account credentials, project ID, and Github webhook token.

These resources and any other necessary fixtures should be managed in the [terraform/test_fixtures](./terraform/test_fixtures) configuration. Each module should have a dedicated `.tf` file that creates its fixtures and K8s secret.

### Concourse Pipelines

The [concourse/pipelines](./concourse/pipelines) directory contains [pipeline](https://concourse-ci.org/pipelines.html) configurations for CFT modules. There are Make targets to facilitate creating/upating pipelines. For example, to update the terraform-google-project-factory pipeline, run the following:

```
cd concourse
make project-factory
```

### Concourse Docker Images

The [concourse/build](./concourse/build) directory contains Dockerfiles to use in tasks. The images are currently built manually and pushed to the GCP Container Registry associated with the `cloud-foundation-cicd` project. There are Make targets to facilitate building and pushing the image - see [concourse/Makefile.BUILD](./concourse/Makefile.BUILD). For example, to build and release the `cft/lint` image, run the following:

```
cd concourse
make build-lint-image
make release-lint-image
```

### Github Pull Request Integration

Concourse can be integrated with Github pull requests using the [github-pr-resource](https://github.com/telia-oss/github-pr-resource). The integration require s a `webhook_token` and an `access_token`. For example, the `project-factory` pull request resource can be configured with:

```
resource_types:

- name: pull-request
  type: docker-image
  source:
    repository: teliaoss/github-pr-resource

resources:

- name: pull-request
  type: pull-request
  webhook_token: ((project-factory.github_webhook_token))
  source:
    repository: terraform-google-modules/terraform-google-project-factory
    access_token: ((github.pr-access-token))

...
```

#### Webhook Token

The `webhook_token` is user defined. We use the Terraform `random_id` resource to generate it. The value is written to a Kubernetes secret. The `webhook_token` value in the pipeline's resource configuration must correspond with a Github webhook payload URL. For convenience, the [terraform/test_fixtures](./terraform/test_fixtures) outputs the full payload URLs:

```
cd terraform/test_fixtures
terraform output
```

Output:

```
github_webhook_urls = {
  terraform-google-container-vm = https://...
  terraform-google-kubernetes-engine = https://...
  terraform-google-project-factory = https://...
  ...
}
```

At this time, the webhooks must be configured manually on each Github repository.

#### Access Token

The access token is a personal access token of the Github machine user, `cloud-foundation-cicd`. The token lives in the `github` K8s secret of the `concourse-cft` namespace. The login password also lives in this secret.

```
kubectl get secret github --namespace concourse-cft -o yaml
```

## Workstation Setup

### Terraform

All configurations currently use Terraform 0.11.10.

You will need to create two service accounts to interact with Terraform - one for the google.com project `cloud-foundation-cicd` and one for the `phoogle.net` organization.

The google.com service account should have the Project Owner role on the `cloud-foundation-cicd` project, and the phoogle.net service account should have Organization Admin and Folder Admin roles.

Copy `.env.sample` locally and update as necessary. You will need to specify paths to credientials for those two service accounts. As well, there is a `TF_VAR_postgres_concourse_user_password` environment variable - this is only necessary for interacting with the [terraform/postgres](./terraform/postgres) configuration. You may obtain this value from Kubernetes: `kubectl get secret concourse-concourse -o yaml` and Base64 decode the value for `postgresql-password`.

### Concourse

The Concourse UI URL isk https://concourse.infra.cft.tips/.

From the UI, download the [Fly CLI](https://concourse-ci.org/fly.html).
