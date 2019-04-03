
# Cloud Foundation Toolkit - User Guide

<!-- TOC -->

- [Overview](#overview)
- [CFT Configs](#cft-configs)
    - [Extra YAML Directives](#extra-yaml-directives)
        - [name](#name)
        - [project](#project)
        - [description](#description)
    - [Extra Features](#extra-features)
        - [Cross-deployment References with the `$(out)` Tag](#cross-deployment-references-with-the-out-tag)
        - [Jinja Templating](#jinja-templating)
    - [Samples](#samples)
        - [network.yaml](#networkyaml)
        - [firewall.yaml](#firewallyaml)
        - [instance.yaml](#instanceyaml)
- [Templates](#templates)
- [Toolkit Installation and Configuration](#toolkit-installation-and-configuration)
    - [Installing Prerequisites](#installing-prerequisites)
        - [Python 2.7 + pip](#python-27--pip)
        - [Google Cloud SDK](#google-cloud-sdk)
    - [Getting the CFT Code](#getting-the-cft-code)
    - [Installing the CFT](#installing-the-cft)
    - [Uninstalling the CFT](#uninstalling-the-cft)
    - [Updating the CFT](#updating-the-cft)
- [CLI Usage](#cli-usage)
    - [Syntax](#syntax)
    - [Actions](#actions)
        - [The "create" Action](#the-create-action)
        - [The "update" Action](#the-update-action)
        - [The "apply" Action](#the-apply-action)
        - [The "delete" Action](#the-delete-action)

<!-- /TOC -->

## Overview

The GCP Deployment Manager service does not support cross-deployment
references, and the `gcloud` utility does not support concurrent deployment of
multiple inter-dependent configs. The `Cloud Foundation toolkit` (henceforth,
`CFT`) expands the capabilities of Deployment Manager and `gcloud` to support
the following scenarios:

- Creation, update, and deletion of multiple deployments in a single operation
  which:
  - Accepts multiple config files as input
  - Automatically resolves dependencies between these configs
  - Creates/updates deployments in the dependency-stipulated order, or
    deletes deployments in a reverse dependency order
- Cross-deployment (including cross-project) referencing of deployment outputs,
  which removes the need for hard-coding many parameters in the configs

For example, if config file `A` contained all network resources, config file
`B` contained all instances, and config `C` contained firewall rules, router,
and VPN, in `gcloud` you would need to *manually* define the config deployment
order according to the resource dependencies. The VPN would depend on the cloud
router, both of them would depend on the network, etc. The `CFT` computes the
dependencies *automatically*, which eliminates the need for manual deployment
ordering.

`Note:` This User Guide assumes that you are familiar with the Google Cloud SDK
operations related to resource deployment and management. For additional
information, refer to the
[SDK documentation](https://cloud.google.com/sdk/docs/).

The CFT includes:

- A command-line interface (henceforth, CLI) that deploys resources defined in
  single or multiple CFT-compliant config files
- A comprehensive set of production-ready resource [templates](#templates) that follow
  Google's best practices, which can be used with the CFT or the `gcloud`
  utility. (`gcloud` is part of the Google Cloud SDK).

You can use the CFT "as is" or modify it to suit your specific needs. Instructions
and recommendations for the CFT code modifications are in the
[CFT Developer Guide](tool_dev_guide.md).

## CFT Configs

To use the CFT, you need to first create the config files for the desired
deployments. These configs are YAML structures very similar to, and compatible
with, the `gcloud` config files. The difference is that they contain extra YAML
directives and features to support the expanded capabilities of the CFT
(multi-config deployment and cross-deployment references).

### Extra YAML Directives

#### name

This directive is used to specify the name of the deployment; for example:

```yaml
name: my-network
```

If not specified, the name of the deployment is inferred from the config
file name. For example, if the path to the config file is
`path/to/configs/my-network.yaml`, and the config does not specify the `name`
directive, the deployment name is set to `my-network`. This is meant as a
workaround for maintaining compatibility between the `CFT` and `gcloud` configs.
However, **it is strongly recommended that the `name` directive is specified**.

#### project

This directive defines the project in which the resource is deployed; for
example:

```yaml
project: my-project
```

While this directive is optional, **its use in your configs is highly
recommended**. In addition to the project directive in the config file,
the project for a deployment to be created in can be specified by other means.
The order of precedence is as follows:

1. The `--project` command-line option. If a project is specified via this
   option, all configs in the run use that project. This is a way of
   quickly overriding the project specified in a config file, which should be
   used with caution.
2. The `project` directive in the config file.
3. The `CLOUD_FOUNDATION_PROJECT_ID` environment variable.
4. The "default project" configured with the GCP SDK.

`Note:` When deployments utilize cross-project resources, the `project`
directive becomes mandatory in at least one of the deployments.

#### description

This directive is the deployment description, which allows you
to add some documentation to your configs; for example:

```yaml
description: My firewall deployment for {{environment}} environment
```

### Extra Features

#### Cross-deployment References with the `$(out)` Tag

A config/deployment can specify a dependency on another deployment's output
without the need to create the dependent deployment in advance. This is the
mechanism the CFT uses to determine the order of execution of the deployments.

```yaml
$(out.<project>.<deployment>.<resource>.<output>)

# or

$(out.<deployment>.<resource>.<output>)
```

wherein:

- `$(out)` is the prefix that indicates that the value references an output
  from a resource defined in an external deployment (in another config file)
- `project` is the ID of the project in which the external deployment is
   created
- `deployment` is the he name of the external deployment (config) that
  defines the referenced resource
- `resource` is the DM name of the referenced resource
- `output` is the name of the output parameter to be referenced

The above construct works very similarly to Deployment Manager's
`$(ref.<resource>.<property>)`. However, it allows defining not only references
to resource properties not only *within* a deployment, but also
*inter-deployment/inter-project* references, using deployment outputs. The
value of output of a dependent deployment is only looked up during the current
deployment's execution, which allows you to create config files without knowing
in advance the actual values of the outputs in the dependent deployments, or
even having to create these deployments.

For example:

```yaml
network: $(out.my-network-prod.my-network-prod.name)
```

#### Jinja Templating

All configs submitted via the CFT CLI are rendered by the [Jinja Template
Engine](http://jinja.pocoo.org/). This supports compact code by using the DRY
pattern. For example, by using variable substitution and `for loops`:

```yaml
{% set environment = 'prod' %}
{% set applications = ['app1', 'app2', 'app3'] %}

name: my-network-{{environment}}
description: Network deployment for {{environment}} environment
project: sourced-gus-1
imports:
  - path: templates/network/network.py
resources:
{% for application in applications %}
  - type: templates/network/network.py
    name: {{application}}-{{environment}}-network
    properties:
      autoCreateSubnetworks: false
{% endfor %}
```

An alternative to using Jinja in your configs is to write wrapper DM Python
templates and reference these templates in your configs (see the
[Templates](#templates) section).

### Samples

Following are three sample config files that illustrate the above directives
and features. These will be used as examples in the action-specific sections of
this User Guide:

- [network.yaml](#network.yaml) - two networks that have no dependencies
- [firewall.yaml](#firewall.yaml) - two firewall rules, which depend on the
  corresponding networks
- [instance.yaml](#instance.yaml) - one VM instance, which depends on the
  network

#### network.yaml

```yaml
name: my-networks
description: my networks deployment

imports:
  - path: templates/network/network.py

resources:
  - type: templates/network/network.py
    name: my-network-prod
    properties:
      autoCreateSubnetworks: true

  - type: templates/network/network.py
    name: my-network-dev
    properties:
      autoCreateSubnetworks: false
```

#### firewall.yaml

```yaml
name: my-firewalls
description: My firewalls deployment

imports:
  - path: templates/firewall/firewall.py
resources:
  - type: templates/firewall/firewall.py
    name: my-firewall-prod
    properties:
      network: $(out.my-networks.my-network-prod.name)
      rules:
        - name: allow-proxy-from-inside-prod
      allowed:
            - IPProtocol: tcp
              ports:
                - "80"
                - "444"
          description: This rule allows connectivity to the HTTP proxies
          direction: INGRESS
          sourceRanges:
            - 10.0.0.0/8
        - name: allow-dns-from-inside-prod
          allowed:
            - IPProtocol: udp
              ports:
                - "53"
            - IPProtocol: tcp
              ports:
                - "53"
          description: this rule allows DNS queries to google's 8.8.8.8
          direction: EGRESS
          destinationRanges:
            - 8.8.8.8/32
  - type: templates/firewall/firewall.py
    name: my-firewall-dev
    properties:
      network: $(out.my-networks.my-network-dev.name)
      rules:
        - name: allow-proxy-from-inside-dev
          allowed:
            - IPProtocol: tcp
              ports:
                - "80"
                - "444"
          description: This rule allows connectivity to the HTTP proxies
          direction: INGRESS
          sourceRanges:
            - 10.0.0.0/8
```

#### instance.yaml

```yaml
name: my-instance-prod-1
description: My instance deployment for prod environment

imports:
  - path: templates/instance/instance.py
    name: instance.py

resources:
  - name: my-instance-prod-1
    type: instance.py
    properties:
      zone: us-central1-a
      diskImage: projects/ubuntu-os-cloud/global/images/family/ubuntu-1804-lts
      diskSizeGb: 100
      machineType: f1-micro
      diskType: pd-ssd
      network: $(out.my-networks.my-network-prod.name)
      metadata:
        items:
          - key: startup-script
            value: sudo apt-get update && sudo apt-get install -y nginx
```

## Templates

CFT-compliant configs can use templates written in Python or Jinja2. [Templates
included in the toolkit](../templates/README.md) are recommended (although not mandatory)
as they offer robust functionality, ease of use, and adherence to best
practices.

You can use the templates included in our library "as is," and/or modify them
to suit your needs, as well as develop your own templates. Instructions and
recommendations for template development are in the
[Template Developer Guide](template_dev_guide.md).  

## Toolkit Installation and Configuration

This toolkit was developed primarily on/for Linux. Therefore, the Linux platform
is expected to offer the most seamless user experience.

### Installing Prerequisites

#### Python 2.7 + pip

Follow your OS package manager instructions. For example, for Ubuntu:

```shell
sudo apt-get install python2.7
sudo apt-get install python-pip
```

#### Google Cloud SDK

1. Install the [Google Cloud SDK](https://cloud.google.com/sdk/docs/quickstarts).
2. Ensure that the `gcloud` command is in the user's PATH:

```shell
which gcloud
```

### Getting the CFT Code

Proceed as follows:

```shell
git clone https://github.com/GoogleCloudPlatform/deploymentmanager-samples
cd deploymentmanager-samples/
```

### Installing the CFT

Proceed as follows:

```shell
cd community/cloud-foundation
sudo make prerequisites        # Installs prerequisites in system python
make build                     # builds the package
sudo make install              # installs the package in /usr/local
```

### Uninstalling the CFT

If you need to uninstall the CFT, proceed as follows:

```shell
sudo make uninstall
```

### Updating the CFT

To update CFT to a newer version, proceed as follows:

```shell
cd community/cloud-foundation
make clean
sudo make prerequisites
make build
sudo make uninstall
sudo make install
```

## CLI Usage

### Syntax

The CLI commands adhere to the following syntax:

```shell
cft [action] [configs] [action-options]
```

The above syntactic structure includes the following elements:

- `[action]` - one of the supported actions/commands:
  - **create** - creates deployments defined in the specified config files,
    in the dependency order
  - **update** - updates deployments defined in the specified config files,
    in the dependency order
  - **apply** - checks if the resources defined in the specified configs
    already exist; if they do, updates them; if they don't, creates them
  - **delete** - deletes deployments defined in the specified config files,
    in the reverse dependency order
- `[config]` - The path(s) to the config files to be affected by the specified
  action (files with extensions `.yaml`, `.yml`, or `.jinja`). It can be:
  - A space-separated list of paths to the config files and/or directories,
    optionally with wildcards; for example:
    - ../deployments/config_1.yaml ../tests/test*.yaml ../dev/config*.yml
    - ../deployments/ ../tests/ - this will submit all files with extensions
      `.yaml`, `.yml`, or `.jinja` found in the ../deployments/ and ../tests/
      directories
  - A space-separated list of yaml-serialized strings, each representing a
    config; useful when another tool is generating configs on the fly\
    For example: `name: my-networks\nproject: my-project\nimports:\n  - path: templates/network/network.py\n    name: network.py\resources:\n  - type: templates/network/network.py\n    name: my-network-prod`
- `[action-options]` - one or more action-specific options; see the
  action-specific `--help` option for details:

```shell
cft --help
usage: cft [-h] [--version] [--project PROJECT] [--dry-run]
           [--verbosity VERBOSITY]
           {apply,create,update,delete} ...

positional arguments:
  {apply,create,update,delete}

optional arguments:
  -h, --help            show this help message and exit
  --version, -v         Print version information and exit
  --project PROJECT     The ID of the GCP project in which ALL config files
                        will be executed. This option will override the
                        "project" directive in the config files, so be careful
                        when using this
  --dry-run             Prints the order of execution of the configs. No
                        changes are made
  --verbosity VERBOSITY
                        The log level
```

### Actions

The `CFT` parses the submitted config files and computes the dependencies
between them. Based on the computed dependency graph, the script
determines the sequence of deployments to be executed. It then proceeds to
execute the action in the computed order.

#### The "create" Action

`Note:` Make sure that the deployments you are going to create do not exist in
your DM. An attempt to create a deployment that already exists will result in
an error. Yon can, however, do one of the following:

- Use the **update** action to update the existing deployments - see
  [The "update" Action](#the-update-action) section
- Use the **apply** action which will attempt to create the deployment if it
  doesn't already exist in DM, or update the deployment it already exist - see
  [The "apply" Action](#the-apply-action) section

To create multiple deployments, in the CLI, type:

```shell
cft create [configs] [create-options]
```

If you submit the [sample configs described above](#samples)

```shell
cft create instance.yaml firewall.yaml network.yaml
```

the following response appears in the CLI terminal:

```shell
---------- Stage 1 ----------
Waiting for insert my-network-prod (fingerprint 7OyDHEL8-ZGbay4dTcXXEg==) [operation-1538159159516-576f2964f9b61-e64bdb44-8ab51124]...done.
NAME             TYPE                STATE      ERRORS  INTENT
my-network-dev   compute.v1.network  COMPLETED  []
my-network-prod  compute.v1.network  COMPLETED  []
---------- Stage 2 ----------
Waiting for insert my-instance-prod-1 (fingerprint tdbkal-dX_ppamFJVtBGew==) [operation-1538159204094-576f298f7d030-9707b687-a3f822d9]...done.
NAME                TYPE                 STATE      ERRORS  INTENT
my-instance-prod-1  compute.v1.instance  COMPLETED  []
Waiting for insert my-firewall-prod (fingerprint Yuhd7khES_en86QtLYFV8w==) [operation-1538159238360-576f29b02abc2-b29dacc3-1b74eb12]...done.
NAME                          TYPE                   STATE      ERRORS  INTENT
allow-dns-from-inside-prod    compute.beta.firewall  COMPLETED  []
allow-proxy-from-inside-dev   compute.beta.firewall  COMPLETED  []
allow-proxy-from-inside-prod  compute.beta.firewall  COMPLETED  []
---------- Stage 3 ----------
Waiting for insert my-instance-prod-2 (fingerprint z-lJJimsanFI6cIYLU8D_w==) [operation-1538159270905-576f29cf344a8-d28b6852-52527e20]...done.
NAME                TYPE                 STATE      ERRORS  INTENT
my-instance-prod-2  compute.v1.instance  COMPLETED  []
```

In this example, the network config has no dependencies, and the firewall and
instance configs depend on the network. Therefore, the network config is
deployed first (Stage 1), and the firewall and instance are deployed next
(Stage 2).

`Note:` The order in which the configs are provided in the `cft create` command
does not affect the deployment creation order. That order is defined
exclusively by the dependency between the configs, which is, in turn, defined
by analyzing and ordering the cross-dependency tokens (`$(out.a.b.c.d)`).

The following conditions will result in the action failure,
with an error message displayed:

- One or more of the specified deployments already exist
- One or more resources in the submitted config files depend on resources that
  neither exist nor being created by the current `create` action
- One or more of the submitted config files are invalid
- One or more of the submitted config files contain circular dependencies
  (i.e., deployment A depends on deployment B, and B depends on A)

#### The "update" Action

`Note:` Make sure that the deployments you are going to update already exist in
DM. An attempt to update deployment that does not exist will result in an
error. Yon can, however, do one of the following:

- Use the **create** action to create the required deployments - see
  [The "create" Action](#the-create-action) section
- Use the **apply** action which will attempt to create the deployment if it
  doesn't already exist in DM, or update the deployment it already exist - see
  [The "apply" Action](#the-apply-action) section

To update multiple configs, in the CLI, type:

```shell
cft update [configs] [create-options]
```

If you submit the [sample configs described above](#samples)

```shell
cft update instance.yaml firewall.yaml network.yaml
```

the following response appears in the CLI terminal:

```shell
---------- Stage 1 ----------
Waiting for update my-network-prod (fingerprint 7OyDHEL8-ZGbay4dTcXXEg==) [operation-1538159159516-576f2964f9b61-e64bdb44-8ab51124]...done.
NAME             TYPE                STATE      ERRORS  INTENT
my-network-dev   compute.v1.network  COMPLETED  []
my-network-prod  compute.v1.network  COMPLETED  []
---------- Stage 2 ----------
Waiting for update my-instance-prod-1 (fingerprint tdbkal-dX_ppamFJVtBGew==) [operation-1538159204094-576f298f7d030-9707b687-a3f822d9]...done.
NAME                TYPE                 STATE      ERRORS  INTENT
my-instance-prod-1  compute.v1.instance  COMPLETED  []
Waiting for update my-firewall-prod (fingerprint Yuhd7khES_en86QtLYFV8w==) [operation-1538159238360-576f29b02abc2-b29dacc3-1b74eb12]...done.
NAME                          TYPE                   STATE      ERRORS  INTENT
allow-dns-from-inside-prod    compute.beta.firewall  COMPLETED  []
allow-proxy-from-inside-dev   compute.beta.firewall  COMPLETED  []
allow-proxy-from-inside-prod  compute.beta.firewall  COMPLETED  []
---------- Stage 3 ----------
Waiting for update my-instance-prod-2 (fingerprint z-lJJimsanFI6cIYLU8D_w==) [operation-1538159270905-576f29cf344a8-d28b6852-52527e20]...done.
NAME                TYPE                 STATE      ERRORS  INTENT
my-instance-prod-2  compute.v1.instance  COMPLETED  []
```

In this example, the network config has no dependencies, and the firewall and
instance configs depend on the network. Therefore, the network config is
updated first (Stage 1), and the firewall and instance are updated next
(Stage 2).

The following conditions will result in the actin failure, with an error
message displayed:

- One or more of the specified deployments do not exist
- One or more resources in the submitted config files depend on resources that
  do not exist
- One or more of the submitted config files are invalid
- One or more of the submitted config files contain circular dependencies
  (i.e., deployment A depends on deployment B, and B depends on A)

You can use the `--preview` option with the `update` action; for example:

```shell
cft update test/fixtures/configs/ --preview
```

The CFT puts each deployment in the `preview` mode within DM, displays a
preview of the action results, and enables you to approve/decline the action
for each of the submitted configs. The following prompt is displayed after
the Stage 1 log:

```shell
Update(u), Skip (s), or Abort(a) Deployment?
```

Having reviewed the displayed information, enter one of the following
responses:

- **u (update)** - confirms the deployment change as shown in the preview
- **s (skip)** - cancels the update (no change) and continues to the next
  config in the sequence
- **a (abort)** - cancels the update (no change) and aborts the script
  execution

#### The "apply" Action

The **apply** action makes the CFT decide which deployments must be created
(because they do not exist), and which ones must be updated (because they do
exist).

To create or update multiple configs, in the CLI, type:

```shell
cft apply [configs] [create-options]
```

If you submit the [sample configs described above](#samples)

```shell
cft apply instance.yaml firewall.yaml network.yaml
```

the following response appears in the CLI terminal:

```shell
---------- Stage 1 ----------
Waiting for update my-network-prod (fingerprint 7OyDHEL8-ZGbay4dTcXXEg==) [operation-1538159159516-576f2964f9b61-e64bdb44-8ab51124]...done.
NAME             TYPE                STATE      ERRORS  INTENT
my-network-dev   compute.v1.network  COMPLETED  []
my-network-prod  compute.v1.network  COMPLETED  []
---------- Stage 2 ----------
Waiting for update my-instance-prod-1 (fingerprint tdbkal-dX_ppamFJVtBGew==) [operation-1538159204094-576f298f7d030-9707b687-a3f822d9]...done.
NAME                TYPE                 STATE      ERRORS  INTENT
my-instance-prod-1  compute.v1.instance  COMPLETED  []
Waiting for update my-firewall-prod (fingerprint Yuhd7khES_en86QtLYFV8w==) [operation-1538159238360-576f29b02abc2-b29dacc3-1b74eb12]...done.
NAME                          TYPE                   STATE      ERRORS  INTENT
allow-dns-from-inside-prod    compute.beta.firewall  COMPLETED  []
allow-proxy-from-inside-dev   compute.beta.firewall  COMPLETED  []
allow-proxy-from-inside-prod  compute.beta.firewall  COMPLETED  []
---------- Stage 3 ----------
Waiting for update my-instance-prod-2 (fingerprint z-lJJimsanFI6cIYLU8D_w==) [operation-1538159270905-576f29cf344a8-d28b6852-52527e20]...done.
NAME                TYPE                 STATE      ERRORS  INTENT
my-instance-prod-2  compute.v1.instance  COMPLETED  []
```

The following conditions will result in the action failure, with an error
message displayed:

- One or more resources in the submitted config files depend on resources that
  neither exist nor being created by the current `apply` action
- One or more of the submitted config files are invalid
- One or more of the submitted config files contain circular dependencies
  (i.e., deployment A depends on deployment B, and B depends on A)

You can use the `--preview` option with the `apply` action; for example:

```shell
cft apply test/fixtures/configs/ --preview
```

The CFT puts each deployment in the `preview` mode within DM, displays a
preview of the action results, and enables you to approve/decline the action
for each of the submitted configs. The following prompt is displayed after
the Stage 1 log:

```shell
Update(u), Skip (s), or Abort(a) Deployment?
```

Having reviewed the displayed information, enter one of the following
responses:

- **u (update)** - confirms the deployment change as shown in the preview
- **s (skip)** - cancels the update (no change) and continues to the next
  config in the sequence
- **a (abort)** - cancels the update (no change) and aborts the script
  execution

`Note:` If the `apply` action is creating (rather than updating) a set of
resources, and if you choose to skip the creation of a deployment on which
subsequent deployments depends (e.g., **skip** network in Stage 1 and
**update** firewall in Stage 2), the operation will fail with an error message.

#### The "delete" Action

To delete the previously created/updated multiple deployments, in the CLI, type:

```shell
cft delete [configs] [create-options]
```

If you submit the [sample configs described above](#samples)

```shell
cft delete instance.yaml firewall.yaml network.yaml
```

the following response appears in the CLI terminal:

```shell
---------- Stage 1 ----------
Waiting for delete my-instance-prod-2 (fingerprint 3IWMMfbjsUWjtWgvs6Evdw==) [operation-1538159406282-576f2a504f510-2dceed8f-b222b564]...done.
---------- Stage 2 ----------
Waiting for delete my-instance-prod-1 (fingerprint ifQgUyTSOtVE1H6VgaIlYA==) [operation-1538159505990-576f2aaf66170-fcc5246d-2d44d005]...done.
Waiting for delete my-firewall-prod (fingerprint xFs1fcZiLJPVV1hUw61-og==) [operation-1538159629835-576f2b2581af9-a83468de-d3685d90]...done.
---------- Stage 3 ----------
Waiting for delete my-network-prod (fingerprint EhMN6C5IeADJYRo40CmuAg==) [operation-1538159649120-576f2b37e5f02-35da3a44-cf279bfa]...done.
```

The order of execution for `delete` is reversed (compared to `create` or
`update`). This prevents DM from attempting to delete, for example, a network
resource while an instance resource (dependent on the network) still exists.

`Note:` The CFT silently ignores deletion of deployments that do not exits.
This covers those cases where the deletion of a specific deployment had
failed and the problem was then fixed. You do not have to figure out which
deployments to delete; you simply re-run the command.
