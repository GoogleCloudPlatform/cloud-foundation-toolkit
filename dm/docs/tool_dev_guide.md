# CFT Developer Guide

<!-- TOC -->

- [Overview](#overview)
- [Prerequisites](#prerequisites)
    - [Google Cloud SDK](#google-cloud-sdk)
    - [Development Environment](#development-environment)
- [Unit Tests](#unit-tests)
    - [From Outside the Development Environment](#from-outside-the-development-environment)
    - [From Within the Development Environment](#from-within-the-development-environment)

<!-- /TOC -->

## Overview

The Cloud Foundation toolkit (henceforth, CFT) includes the following parts:
project is comprised of two parts:

- A comprehensive set of [production-ready resource templates](../templates/README.md)
  that follow Google's best practices, which can be used with the CFT or the
  gcloud utility (part of the Google Cloud SDK)
- A command-line interface (henceforth, CLI) that deploys resources defined in
  single or multiple CFT-compliant config files - see the
  [CFT User Guide](userguide.md)

This Guide is intended for the developers who are planning to modify and/or
programmatically interface with the CFT.

## Prerequisites

### Google Cloud SDK

1. Install the [Google Cloud SDK](https://cloud.google.com/sdk/), which
   includes the `gcloud` CLI.

Because the SDK is not in *pypi*, its installation cannot be easily
automated from within this project, due to the fact that users on the different
platforms need the different packages. Follow the SDK installation instructions
for your platform.

2. Ensure that the `gcloud` CLI is in your user PATH (because the CFT uses
   this CLI to find the location of the Python libraries included in the SDK).

The `gcloud` CLI is usually placed in the PATH automatically when you:

- Install the SDK via the official package manager for your OS (RPM, DEB,
  etc.), or
- Use the installer (`install.sh`) bundled in a Linux tarball

However, if you used neither of the above installation methods, you need to
ensure that `gcloud` can be found in one of the directories specified by the
PATH environment variable.

### Development Environment

The CFT development environment is based on:

- [Tox](https://tox.readthedocs.io/en/latest/index.html) for streamlined
  management of Python virtual environments
- [pytest](https://docs.pytest.org/en/latest/contents.html) for unit tests

Proceed as follows:

1. Install Tox with the system Python.
2. Install CFT prerequisites:

```shell
sudo make cft-prerequisites
```

The CFT development is carried out in a virtual environment.

3. Create the virtual development environment called `venv` with `tox` in
   the root of the project directory:

```shell
make cft-venv
```

4. Activate the virtual environment:

```shell
source venv/bin/activate
source src/cftenv
```

The above activates the virtual environment, then finds the Google SDK path
and adds libraries to PYTHONPATH. These cannot be simply added to the
`Makefile` because `make` creates sanitized sub-shells for each command, and
the parent shell does not get the environment variables that the virtual
environment sets up on activation.

`Note:` The `tox.ini` file in this project is configured to
"*install*" the utility using pip's "develop" mode, i.e., the pip **does not**
actually package and install the utility in the virtual environment's
`site-packages`.

5. To install or update any of the packages in your virtual environment
   (created by `tox`), delete and re-create the environment:

- *Deactivate* the virtual environment (if it has been activated):

```shell
deactivate
unset CLOUDSDK_ROOT_DIR CLOUDSDK_PYTHON_SITEPACKAGES PYTHONPATH
```

- Delete the deactivated virtual environment:

```shell
make cft-clean-venv
```

- Create the environment as described in Step 3 above.

## Unit Tests

You can run the CFT unit tests either from withing your development
environment or from outside of it.

### From Outside the Development Environment

This testing mode is typically used when running tests from a CI tool.

1. Use `tox` to create the necessary virtual environments (not `venv`, which
   is used only for active development):

```shell
make cft-test
```

2. Run all the tests within the "test" virtual environments.

### From Within the Development Environment

This testing mode is typically used while actively developing within the
development virtual environment.

1. Activate the `venv` environment as shown in Step 4 of the
   [Development Environment](#development-environment) section.
2. Source `src/cftdev` to get PYTHONPATH set as shown in Step 5 of the
   [Development Environment](#development-environment) section.
3. Run tests as follows:

```shell
# use the make target to run all tests:
make cft-test-venv

# alternatively, use pytest directly to run all tests:
python -m pytest -v

# alternatively, run a single test file:
python -m pytest -v tests/unit/test_deployment.py
```