# CFT Sample Pipeline

<!-- TOC -->

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Pipelines](#pipelines)

<!-- /TOC -->

## Overview

You can use the Cloud Foundation toolkit (henceforth, CFT) as a standalone
solution, via its command line interface (CLI) – see
[CFT User Guide](../docs/userguide.md) for details. Alternatively, you can
initiate CFT actions via its API, from a variety of existing orchestration
tools, or from your own application.

This document describes one of the CFT integration scenarios, wherein
you initiate the CFT actions from Jenkins. It uses as an example a
Jenkins-based "sample pipeline", which is included in this CFT directory.

`Note:` This document assumes that you are familiar with the basics of
[Jenkins](https://jenkins.io/) and of its
[Pipeline Plugin](https://jenkins.io/doc/book/pipeline/).

`Note:` The Jenkins-based process this document describes is for demonstration
purposes only. It is not intended as a product. Your Jenkins setup is likely
to be different from the one used for demonstrate. Therefore, to achieve
similar results, you need to modify certain parameters in all the demo files.

## Prerequisites

1. A working Jenkins server:
    - Different organizations have vastly different Jenkins setups. Therefore,
      this document provides no specific recommendations for fulfilling this
      prerequisite. You might use a Compute Image from
      [Marketplace](https://console.cloud.google.com/marketplace/browse?q=jenkins).
    - Install the Pipeline Utility Steps plugin.
2. GCP Service Accounts (SA):
    - `Service Account for Jenkins`: Jenkins must be configured with
      permissions sufficient for managing DM deployments. This can be achieved
      by:
      - Associating a SA with the GCP Compute Instance running Jenkins (if
        Jenkins is in GCP), or
      - Configuring the SA credentials with the Jenkins user (if running
        Jenkins outside GCP)
    - `Service Account for the GCP project` (a.k.a. the DM Service Account):
      this SA needs permissions to all APIs DM uses to create resources.
3. The Cloud Foundation toolkit:
    - CFT must be installed in the Jenkins master and slaves. For installation
      instructions, see the [CFT User
      Guide](../docs/userguide.md#toolkit-installation-and-configuration).
    - Note that the [Google Cloud SDK](https://cloud.google.com/sdk) is a
      prerequisite for the CFT.
4. The Environment Variables file:
    - An example file is [here](pipeline-vars).
    - Replace <FIXME:XXX> with values specific to you organization, and move
      the file to the Jenkins user's home directory.

## Pipelines

This directory implements deployment pipelines to show how the CFT can be used
in a *fictitious company*. In this fictitious company, three separate teams are
responsible for the corresponding separate pieces of the cloud infrastructure:

- Central Cloud Platform Team:
  - Responsible for creating GCP projects, IAM entities, Permissions,
    Billing, etc.
  - Owns the pipeline and configs in [project](project)
- Central Networking Team:
  - Responsible for networking between for all other teams, interconnects,
    on-premise integration, etc.
  - Owns the pipeline and configs in [network](network)
- Application Teams (typically, more than one):
  - Responsible for deploying the team-specific application stack (in this
    example, there is a single application team, which is responsible for
    deploying its GKE clusters in the different environments)
  - Owns the pipeline and configs in [app](app)

Each folder in this directory of the CFT repository represents and implements
a pipeline that corresponds to one of the above teams.

`Note:` This is not a typical way of organizing Jenkins pipelines. Normally,
each pipeline would be in its own Git repository, with its own access controls
for the different teams.