# Config Connector Solutions

## Overview

Config Connector Solutions provides best practice solutions
to common cloud applications, formatted as YAML definitions
for Config Connector CRDs. These YAMLs can be applied to
clusters running [Config
Connector](https://cloud.google.com/config-connector/docs/how-to/getting-started).

## Structure

Folders under this directory denote general solution areas.
In each solution area folder, there are folders for each package
& customization tool (currently helm and kpt), under which are nested all available solutions in
that solution area and package format.

## Solutions

The full list of solutions grouped by area:

* **apps** - automate creation of a canonical sample application and provision required GCP services with Config Connector
  * wordpress [ [helm](apps/helm/wordpress) ] - provision Wordpress application powered by GCP MySQL database
* **projects** - automate creation of GCP projects, folders and project services
  using Config Connector
  * owned-project [ [kpt](projects/kpt/owned-project) ] - grant the project
    owner role
  * project-hierarcy [ [kpt](projects/kpt/project-hierarchy) ] - get started
    with a folder and a project
  * project-services [ [kpt](projects/kpt/project-services) ] - enable GCP API
    for a project
  * shared-vpc [ [kpt](projects/kpt/shared-vpc) ] - create a shared VPC network
  * simple-project [ [kpt](projects/kpt/simple-project) ] - get started with a
    simple project
* **iam** - automate the management of IAM roles for resources using Config
  Connector
  * folder-iam [ [kpt](iam/kpt/folder-iam) ] - grant an IAM role to a GCP folder
  * kms-crypto-key [ [kpt](iam/kpt/kms-crypto-key) ] - grant an IAM role to a
    KMS crypto key
  * kms-key-ring [ [kpt](iam/kpt/kms-key-ring) ] - grant an IAM role to a KMS
    key ring
  * member-iam [ [kpt](iam/kpt/member-iam) ] - grant a service account an IAM
    role to a project
  * project-iam [ [kpt](iam/kpt/project-iam) ] - grant an IAM role to a project
  * pubsub-subscription [ [kpt](iam/kpt/pubsub-subscription) ] - grant an IAM
    role to a Pub/Sub subscription
  * pubsub-topic [ [kpt](iam/kpt/pubsub-topic) ] - grant an IAM role to a
    Pub/Sub topic
  * service-account [ [helm](iam/helm/service-account) ] \[ [kpt](
    iam/kpt/service-account) ] - grant an IAM role to a service account
  * storage-bucket-iam [ [kpt](iam/kpt/storage-bucket-iam) ] - grant an IAM role
    to a storage bucket
  * subnet [ [kpt](iam/kpt/subnetp) ] - grant an IAM role to a subnetwork
* **sql** - automate the creation of Cloud SQL instances, databases, and users
  using Config Connector
  * mysql-ha [ [kpt](sql/kpt/mysql-ha) ] - create a high availability MySQL
    cluster
  * mysql-private [ [kpt](sql/kpt/mysql-private) ] - create a private MySQL
    database
  * mysql-public [ [kpt](sql/kpt/mysql-public) ] - create a public MySQL
    database
  * postgres-ha [ [kpt](sql/kpt/postgres-ha) ] - create a high availability
    PostgreSQL cluster
  * postgres-public [ [kpt](sql/kpt/postgres-public) ] - create a public
    PostgreSQL database


## Usage

### helm

These solutions are consumable as [helm charts](https://helm.sh/docs/topics/charts/).
Common targets for modification are listed in `values.yaml`.

[Install helm](https://helm.sh/docs/intro/install/). These solutions support Helm v.3+.

Common operations, where `PATH` is the path to the relevant solution folder:
* Showing values: `helm show values PATH`
* Validating chart: `helm template PATH`
* Setting chart: `helm install PATH -generate-name`

Comprehensive documentation at
[https://helm.sh/docs/](https://helm.sh/docs/).

### kpt

These samples are consumable as [kpt
packages](https://googlecontainertools.github.io/kpt/).
Common targets for modification are provided kpt setters,
and can be listed with `kpt cfg list-setters`.

* Installing kpt: follow the instructions on [the kpt
GitHub](https://github.com/GoogleContainerTools/kpt).
* Listing setters: See which values are available for kpt to change `kpt cfg list-setters`
* Setting setters: `kpt cfg set DIR NAME VALUE --set-by NAME`

Comprehensive documentation at
[https://googlecontainertools.github.io/kpt/](https://googlecontainertools.github.io/kpt/).

## License

Apache 2.0 - See [LICENSE](/LICENSE) for more information.
