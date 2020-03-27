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
& customization tool (currently only kpt), under which are nested all available solutions in
that solution area and package format.

## Usage

### kpt

Samples are consumable as [kpt
packages](https://googlecontainertools.github.io/kpt/).
Common targets for modification are provided kpt setters,
and can be listed with `kpt cfg list-setters`.

## License

Apache 2.0 - See [LICENSE](/LICENSE) for more information.
