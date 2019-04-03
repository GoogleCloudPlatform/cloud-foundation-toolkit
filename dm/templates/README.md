# CFT Templates

This folder contains the library of templates included in the Cloud Foundation
toolkit (henceforth, CFT).

## Overview

Each template is stored in a folder named after the templated cloud resource;
e.g., "network", "cloud_router", "healthcheck", etc. Each template folder contains:

- README.md - a textual description of the template's usage, prerequisites, etc.
- `resource`.py - the Python 2.7 template file
- `resource`.py.schema - the schema file associated with the template
- examples:
  - `resource`.yaml - a sample config file that utilizes the template
- tests:
  - integration:
    - `resource`.yaml - a test config file
    - `resource`.bats - a bats test harness for the test config

## Usage

You can use the templates included in the template library:

- Via Google's Deployment Manager / gcloud as described in the
  [Google SDK documentation](https://cloud.google.com/sdk/)
- Via the `CFT`, as described in the [CFT User Guide](../docs/userguide.md)

You can use the templates "as is," and/or modify them to suit your needs, as
well as create new ones. Instructions and recommendations for template
development are in the
[Template Developer Guide](../docs/template_dev_guide.md).