# Cloud Foundation Toolkit Project

## Overview

The Cloud Foundation toolkit (henceforth, CFT) includes the following parts:

- A comprehensive set of production-ready resource templates that follow
  Google's best practices, which can be used with the CFT or the gcloud
  utility (part of the Google Cloud SDK) - see
  [the template directory](dm/templates/)
- A command-line interface (henceforth, CLI) that deploys resources defined in
  single or multiple CFT-compliant config files - see:
  - The CFT source Python files (the `src/` directory)
  - The [CFT DM User Guide](dm/docs/userguide.md)
- A growing set of sample [Config
  Connector](https://cloud.google.com/config-connector/docs/overview)
  YAML configurations - see [the solutions directory](config-connector/solutions/)

In addition, the CFT repository includes a sample pipeline that enables running
CFT deployment operations from Jenkins - see the
[pipeline directory](dm/pipeline/README.md).

## License

Apache 2.0 - See [LICENSE](LICENSE) for more information.
