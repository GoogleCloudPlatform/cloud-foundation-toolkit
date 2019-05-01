# Cloud Foundation Infra

IaC for launching CICD infrastructure using [Concourse CI](https://concourse-ci.org/).

## Infrastructure

Concourse runs on [GKE](https://cloud.google.com/kubernetes-engine/docs/) and [Cloud SQL for Postgres](https://cloud.google.com/sql/docs/postgres/).

Concourse is deployed with [Helm](https://helm.sh/) and the official [Concourse chart](https://github.com/helm/charts/tree/master/stable/concourse).

This is all managed with [Terraform](https://www.terraform.io/). All Terraform configurations are located in the [terraform](./terraform) directory of this repository. At this time, `terraform` is intended to be executed locally on an administrator workstation.

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

You will need to create two service accounts to interact with Terraform - one for the google.com project `cloud-foundation-cicd` and one for the `phoogle.net` organization.

The google.com service account should have the Project Owner role on the `cloud-foundation-cicd` project, and the phoogle.net service account should have Organization Admin and Folder Admin roles bound to the [cloud-foundation-cicd][cicd-folder] folder.

#### Workspaces

When managing fixtures for new modules, switch to the primary workspace:

```
terraform workspace select primary
```

#### Google SA

Log into [Service accounts for project cloud-foundation-cicd](https://pantheon.corp.google.com/iam-admin/serviceaccounts?folder=&organizationId=433637338589&project=cloud-foundation-cicd) using your google.com identity.

 1. Create `<user>@cloud-foundation-cicd.iam.gserviceaccount.com`
 2. Select role: Project / Owner for `cloud-foundation-cicd`.
 3. Create key.  This key is used as the value of `GOOGLE_CREDENTIALS` in
    .env.sample.

#### Phoogle SA

Log into [Permissions for organization
phoogle.net](https://console.cloud.google.com/iam-admin/iam?organizationId=826592752744&orgonly=true&supportedpurview=project) using your phoogle.net identity.

If you do not have a [seed
project](https://github.com/terraform-google-modules/terraform-google-project-factory#script-helper)
yet, create one.  The cloud-foundation-infra SA may be created in any project,
the seed project is used only because it contains the project factory SA and is
a per-user project.

Navigate to Service accounts for your seed project.

 1. Create `cloud-foundation-infra@<ldap>-seed.iam.gserviceaccount.com`
 2. Add role: Service Usage / Service Usage Viewer for the
    cloud-foundation-infra SA from step 1.
 3. Create key.  This key is used as the value of
    `TF_VAR_phoogle_credentials_path` in .env.sample.
 4. Navigate to the [cloud-foundation-cicd][cicd-folder] org level.
 5. Add role: Owner for the cloud-foundation-infra SA from step 1.
 6. Add role: Project Creator for the cloud-foundation-infra SA from step 1.
 7. Navigate to the [phoogle.net Organization][phoogle-org]
 8. Add role: Resource Manager / Organization Administrator for the
    cloud-foundation-infra SA from step 1.
 9. Add role: IAM / Organization Role Administrator for the
    cloud-foundation-infra SA from step 1.
 10. Add role: Resource Manager / Folder Administrator for the
    cloud-foundation-infra SA from step 1.
 11. Add a binding resource to the billing account for role/billing.user the
     memeber SA from step 1.  See [Missing roles/billing.user][billing-user] for
     step by step instructions.

#### Environment Variables

Copy `.env.sample` locally and update as necessary. You will need to specify paths to credientials for those two service accounts. As well, there is a `TF_VAR_postgres_concourse_user_password` environment variable - this is only necessary for interacting with the [terraform/postgres](./terraform/postgres) configuration. You may obtain this value from Kubernetes: `kubectl get secret concourse-concourse -o yaml` and Base64 decode the value for `postgresql-password`.

#### Terraform plan

With the workspace set and Service Accounts configured as per the environment
variables above, a terraform plan in `terraform/test_fixtures/` should succeed.
If you get 403 errors, check the IAM bindings carefully as per above.

#### Future Improvements

The use of two service accounts is questionable because a single logical process
is managing all of these fixtures.  A future improvement could consolidate the
Terraform fixtures under a single machine/process identity.

### Concourse

The Concourse UI URL isk https://concourse.infra.cft.tips/.

From the UI, download the [Fly CLI](https://concourse-ci.org/fly.html).

[cicd-folder]: https://console.cloud.google.com/iam-admin/iam?organizationId=826592752744&orgonly=true&project=&folder=853002531658
[phoogle-org]: https://console.cloud.google.com/iam-admin/iam?organizationId=826592752744&orgonly=true&project=
[billing-user]: https://github.com/terraform-google-modules/terraform-google-project-factory/blob/master/docs/TROUBLESHOOTING.md#missing-rolesbillinguser-role

#### Managing Concourse

##### Setup Concourse `gcloud` and `kubectl` configuration

In order to manage the Concourse GKE cluster you'll need to configure `gcloud`
and `kubectl` with the CICD project and the credentials you set up in the
[Google SA][#google-sa] step.

1. Create a new gcloud configuration
    ```
    $ gcloud config configurations create cloud-foundation-cicd
    ```
2. Activate your cloud-foundation-cicd service account (if not already activated)
    ```
    $ gcloud auth activate-service-account \
        <user>@cloud-foundation-cicd.iam.gserviceaccount.com \
        --key-file=<GOOGLE_CREDENTIALS key file>
    ```
3. Configure the gcloud account and project
    ```
    $ gcloud config set account <user>@cloud-foundation-cicd.iam.gserviceaccount.com
    $ gcloud config set project cloud-foundation-cicd
    $ gcloud config set container/cluster cicd-primary
    ```
4. Load the CI Kubernetes credentials
    ```
    $ gcloud container clusters get-credentials cicd-primary --region us-west1
    ```
5. Verify that you can access Kubernetes resources
    ```
    $ kubectl get nodes
    NAME                                     STATUS   ROLES    AGE   VERSION
    gke-cicd-primary-pool-00-3754cb77-24s1   Ready    <none>   47d   v1.11.5-gke.5
    gke-cicd-primary-pool-00-3754cb77-n9qk   Ready    <none>   47d   v1.11.5-gke.5
    gke-cicd-primary-pool-00-51233741-lbdm   Ready    <none>   19d   v1.11.5-gke.5
    gke-cicd-primary-pool-00-51233741-pr78   Ready    <none>   19d   v1.11.5-gke.5
    gke-cicd-primary-pool-00-60ee1986-s48n   Ready    <none>   76d   v1.11.5-gke.5
    ```

##### Check for stuck workers

Occasionally Concourse workers become stuck. This shows up when jobs start
queueing up even when Concourse jobs can be run in parallel, indicating that
there aren't enough workers running to service all requests.

1. Check Concourse for stalled workers
    **Note** - You can check worker status with either the CFT group (`-t cft`)
    or the main group (`-t main`). In practice this tends to be done with the
    CFT group as those are the credentials you'll use on a daily basis, though
    the main group works equally well if you've already escalated privileges.
    ```
    $ fly -t cft workers
    name                containers  platform  tags  team  state    version
    concourse-worker-0  21          linux     none  none  running  2.1
    concourse-worker-1  13          linux     none  none  running  2.1
    concourse-worker-2  19          linux     none  none  running  2.1

    the following workers have not checked in recently:

    name                containers  platform  tags  team  state    version
    concourse-worker-3  0           linux     none  none  stalled  2.1
    ```
2. Check Kubernetes pod status
    ```
    $ kubectl get pods
    NAME                             READY   STATUS             RESTARTS   AGE
    concourse-web-767bbdf675-6lbns   1/1     Running            0          9d
    concourse-worker-0               1/1     Running            1          9d
    concourse-worker-1               1/1     Running            0          9d
    concourse-worker-2               1/1     Running            16         9d
    concourse-worker-3               0/1     CrashLoopBackOff   3637       9d
    ```

##### Logging into the **main** group

Concourse defines a **main** group that has permissions to directly administer
Concourse. Members of this group can perform maintenance tasks with `fly`, such
as prune stalled workers.

Logging into Concourse with your Google LDAP places you in the **cft** group, so
by default you will not have access to functionality like pruning workers.
You'll need to use the **concourse** user to authenticate as a user in the **main**
group.

1. Fetch the **concourse** user credentials
    ```
    $ kubectl --namespace default get secrets concourse-concourse -o yaml|grep local-users
      local-users: bm90IHRoZSBhY3R1YWwgY3JlZGVudGlhbHMK
    ```
2. Decode the **concourse** credentials
    ```
    $ echo bm90IHRoZSBhY3R1YWwgY3JlZGVudGlhbHMK | base64 --decode
    ```
3.  Log into concourse with the **concourse** user
    **Note** - Log out of the Concourse GUI before you run this step, otherwise
    you'll automatically log into Concourse with your Google LDAP.
    ```
    $ fly login --target main -n main -c https://concourse.infra.cft.tips
    ```


##### Terminate stalled workers

If a worker has become stuck you can delete the failed pod and prune the
stalled worker from Concourse. Kubernetes will automatically re-create the
deleted pod, and Concourse will start sending jobs to the new worker once the
old worker has been pruned.

1.  Terminate a stalled Concourse worker
    **Note** - deleting the pod will cause a new pod to be automatically created.
    ```
    $ kubectl delete pod concourse-worker-3
    pod "concourse-worker-3" deleted
    ```
2. Prune a stalled worker from Concourse
    **Note** - deleting the pod from Kubernetes will clear the faulty pod but will
    prevent the new node from registering - this deletion allows the re-created pod
    to start accepting jobs.
    ```
    $ fly -t main prune-worker -w concourse-worker-3
    ```
