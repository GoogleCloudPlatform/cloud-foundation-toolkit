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

* **apps** - automate creation of a canonical sample application and provisioning required GCP services with Config Connector
  * wordpress [ [helm](apps/helm/wordpress) ] - provision Wordpress application powered by GCP MySQL database
* **projects** - automate creation of GCP projects, folders and project services using Config Connector
  * project-hierarcy [ [kpt](projects/kpt/project-hierarchy) ] - enable GCP API for a project
  * project-services [ [kpt](projects/kpt/project-services) ] - get started with a folder and a project
  * simple-project [ [kpt](projects/kpt/simple-project) ] - get started with a simple project

## Usage

### helm

These samples are consumable as [helm charts](https://helm.sh/docs/topics/charts/).
Common targets for modification are listed in `values.yaml`.

* [Installing helm](https://helm.sh/docs/intro/install/)
* Showing values: `helm show PATH_TO_CHART`
* Validating chart: `helm template PATH_TO_CHART`
* Setting chart: `helm install PATH_TO_CHART -generate-name`

Comprehensive documentation at
[https://helm.sh/docs/](https://helm.sh/docs/).

### kpt

These samples are consumable as [kpt
packages](https://googlecontainertools.github.io/kpt/).

* Installing kpt: follow the instructions on [the kpt
GitHub](https://github.com/GoogleContainerTools/kpt).
* Listing setters: See which values are available for kpt to change `kpt cfg list-setters`
* Setting setters: `kpt cfg set DIR NAME VALUE --set-by NAME`

Comprehensive documentation at
[https://googlecontainertools.github.io/kpt/](https://googlecontainertools.github.io/kpt/).

## License

Apache 2.0 - See [LICENSE](/LICENSE) for more information.
