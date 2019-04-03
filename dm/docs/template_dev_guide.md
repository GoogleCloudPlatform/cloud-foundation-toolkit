# Template Developer Guide

<!-- TOC -->

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Testing](#testing)
    - [Bats Installation](#bats-installation)
    - [Testing Environment Setup](#testing-environment-setup)
        - [Using the Cloud Foundation Config File](#using-the-cloud-foundation-config-file)
        - [Using environment variables](#using-environment-variables)
    - [Running Tests](#running-tests)
    - [Temporary Files and Fixtures](#temporary-files-and-fixtures)

<!-- /TOC -->

## Overview

The Cloud Foundation toolkit (henceforth, CFT) includes the following parts:

- A comprehensive set of [production-ready resource templates](../templates/README.md)
  that follow Google's best practices, which can be used with the CFT or the
  gcloud utility (part of the Google Cloud SDK)
- A command-line interface (henceforth, CLI) that deploys resources defined in
  single or multiple CFT-compliant config files - see the
  [CFT User Guide](userguide.md)

This Guide is intended for the developers who are planning to modify the
existing templates or create new ones.

## Prerequisites

1. Install and set up the [Google Cloud SDK](https://cloud.google.com/sdk/).
2. Install the template development prerequisites:

```shell
make template-prerequisites
```

## Testing

The template consistency and quality control in this project are backed by
simple integration tests using the
[Bats testing framework](https://github.com/sstephenson/bats).

### Bats Installation

To install Bats:

1. Follow the instructions on the Bats
   [website](https://github.com/sstephenson/bats).
2. Make sure the `bats` executable is in your PATH.
3. Alternatively, set up a *development environment* as described in the
   [CFT Developer Guide](tool_dev_guide.md).

### Testing Environment Setup

#### Using the Cloud Foundation Config File

To run tests, you need to modify the organization, project, and
account-specific values in the configuration file. Proceed as follows:

1. Copy `tests/cloud-foundation-tests.conf.example` to
   `~/.cloud-foundation-tests.conf`.
2. Change the values as required.

`Note:` You can modify the configuration file path by changing the
CLOUD_FOUNDATION_CONF environment variable. For example:

```shell
export CLOUD_FOUNDATION_CONF=/etc/cloud-foundation-tests.conf
```

You need to enter the site-specific information (for yourself or for your
organization) in the test config file. See, for example,
`tests/cloud-foundation-tests.conf.example`.

#### Using environment variables

An alternative to using the Cloud Foundation config file is to use environment
variables. Make sure to export all variables described in the
`tests/cloud-foundation-tests.conf.example` file, with your organization-specific
changes.

### Running Tests

`Note:` Currently, only one test file can be executed at a time.

Always run the test from the root of the `cloud-foundation` project:

```shell
./tests/integration/network.bats
 ✓ Creating deployment my-gcp-project-network from my-gcp-project-network.yaml
 ✓ Verifying resources were created in deployment my-gcp-project-network
 ✓ Verifying subnets were created in deployment my-gcp-project-network
 ✓ Deployment Delete
 ✓ Verifying resources were deleted in deployment my-gcp-project-network
 ✓ Verifying subnets were deleted in deployment my-gcp-project-network
```

For the sake of consistency, keep the test files similar, as much as possible,
to the *example configs* available in each template's `examples/` directory.

### Temporary Files and Fixtures

When running tests, temporary Deployment Manager configs and fixtures
are often created and deleted by the *teardown()* function.

Due to the fact that a DM config file must be located relative to the
template(s) it uses, the configs are usually created in the root of the
project. For example, in the `network` template, the config
`.${CLOUD_FOUNDATION_PROJECT_ID}-network.yaml` will be temporarily created
(and deleted at the end of the execution).

Other temporary files are created under `/tmp`; for example:

```shell
/tmp/${CLOUD_FOUNDATION_ORGANIZATION_ID}-network.txt
/tmp/${CLOUD_FOUNDATION_ORGANIZATION_ID}-project.txt
```

The names of the "artifacts" could change. However, if a problem is observed
during the test execution, the root and the /tmp directory are good places to
look for hints about what had caused the problem.