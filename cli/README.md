# Cloud Foundation Toolkit CLI (CFT CLI)

## Usage

The CFT CLI includes a few components, including:
- [CFT Scorecard](./docs/scorecard.md)

Follow cft --help instructions

Google Cloud Foundation Toolkit CLI

```bash
Usage:
  cft [flags]
  cft [command] [flags]

Available Commands:
  help        Help about any command
  report      Generate inventory reports based on CAI outputs in a directory.
  scorecard   Print a scorecard of your GCP environment
  version     Print version information

Flags:
  -h, --help   help for cft

Use "cft [command] --help" for more information about a command.
```


#### Report
Generate inventory reports for resources in Cloud Asset Inventory (CAI) output files, with reports defined in rego (in '<path_to_cloud-foundation-toolkit>/reports/sample' folder).

Example:

```bash
	
	  cft report --query-path <path_to_cloud-foundation-toolkit>/reports/sample \
		--dir-path <path-to-directory-containing-cai-export> \
		--output-path <path-to-directory-for-report-output>
```

## Development

### Requirements

* go 1.12 and above

### Build and Run

```
make build
```

After build find binary at bin/cft location

## License

Apache 2.0 - See [LICENSE](LICENSE) for more information.