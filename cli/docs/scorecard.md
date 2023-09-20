# CFT Scorecard

The CFT Scorecard is integrated into the [CFT CLI](../README.md) and provides
an easy integration with [Forseti Config Validator](https://github.com/forseti-security/policy-library/blob/master/docs/user_guide.md).
It can be used to print a scorecard of your GCP environment, for resources and IAM policies in Cloud Asset Inventory (CAI) exports.
The policies tested are based on constraints and constraint templates from the [Config Validator policy library](https://github.com/forseti-security/policy-library).

## Scorecard User Guide
This tutorial will walk you through setting up Scorecard for a single project.

1. Set some environment variables:
    ```
    export GOOGLE_PROJECT=my-cai-project             # For using CAI API to export CAI data to GCS
    export PUBLIC_BUCKET_NAME=my-bad-public-bucket   # Optional, for triggering a new violation
    export CAI_BUCKET_NAME=my-cai-data               # For downloading CAI data from GCS
    ```
2. Set your project:
    ```
    gcloud config set core/project $GOOGLE_PROJECT   # For exporting CAI data to GCS
    ```
3. Activate the CAI API on your project:
    ```
    gcloud services enable cloudasset.googleapis.com
    ```
4. Create a GCS bucket for storing CAI data:
    ```
    gsutil mb gs://$CAI_BUCKET_NAME
    ```
5. Optionally, create a public GCS bucket to trigger a violation:
    ```
    gsutil mb gs://$PUBLIC_BUCKET_NAME
    gsutil iam ch allUsers:objectViewer gs://$PUBLIC_BUCKET_NAME
    ```
6. Optionally, export the CAI data to GCS:
    ```
    # Export resource data
    gcloud asset export --output-path=gs://$CAI_BUCKET_NAME/resource_inventory.json \
        --content-type=resource \
        --project=$GOOGLE_PROJECT \
        # could also use --folder or --organization
    # Export IAM data
    gcloud asset export --output-path=gs://$CAI_BUCKET_NAME/iam_inventory.json \
        --content-type=iam-policy \
        --project=$GOOGLE_PROJECT \
        # could also use --folder or --organization
    ```
    The alternative is to use [integrated inventory refresh feature](#Using-integrated-inventory-refresh-feature) `cft scorecard --refresh`

7. Download the CFT CLI and make it executable:
    ```
    # OS X
    curl -o cft https://storage.googleapis.com/cft-cli/latest/cft-darwin-amd64
    # Linux
    curl -o cft https://storage.googleapis.com/cft-cli/latest/cft-linux-amd64
    # executable
    chmod +x cft

    # Windows
    curl -o cft.exe https://storage.googleapis.com/cft-cli/latest/cft-windows-amd64
    ```
The user guide in rest of this document provides examples for Linux and OS X environment. For Windows, update file names and directory paths accordingly.

8. Download the sample policy library and add a sample constraint for detecting public buckets:
    ```
    git clone https://github.com/forseti-security/policy-library.git
    cp policy-library/samples/storage_denylist_public.yaml policy-library/policies/constraints/
    ```
9. Run CFT Scorecard:
    ```
    ./cft scorecard --policy-path=./policy-library/ \
        --bucket=$CAI_BUCKET_NAME
    ```
### Using integrated inventory refresh feature
You can also use --refresh flag to create or overwrite CAI export files in GCS bucket and perform analysis, within one step.

```
# Running Cloud Asset Inventory API via Cloud SDK requires a service account and does not support end user credentials.
# Configure Application Default Credential to use a service account key if running outside GCP
# The service account needs be created in a Cloud Asset Inventory enabled project,
# with Cloud Asset Viewer role at target project/folder/org,
# and Storage Object Viewer role at $CAI_BUCKET_NAME
export GOOGLE_APPLICATION_CREDENTIALS=sa_key.json

./cft scorecard --policy-path ./policy-library \
  --bucket=$CAI_BUCKET_NAME \
  --refresh
```

### Using a local export
You can also run CFT Scorecard against locally downloaded CAI data:

```
mkdir cai-dir
gsutil cp gs://$CAI_BUCKET_NAME/resource_inventory.json ./cai-dir/
gsutil cp gs://$CAI_BUCKET_NAME/iam_inventory.json ./cai-dir/
./cft scorecard --policy-path ./policy-library \
  --dir-path ./cai-dir
```

### Full help menu
```
Print a scorecard of your GCP environment, for resources and IAM policies in Cloud Asset Inventory (CAI) exports, and constraints and constraint templates from Config Validator policy library.

	Read from a bucket:
		  cft scorecard --policy-path <path-to>/policy-library \
			  --bucket <name-of-bucket-containing-cai-export>

	Read from a local directory:
		  cft scorecard --policy-path <path-to>/policy-library \
			  --dir-path <path-to-directory-containing-cai-export>

	Read from standard input:
		  cft scorecard --policy-path <path-to>/policy-library \
			  --stdin

	As of now, CAI export file names need to be: resource_inventory.json, iam_inventory.json, org_policy_inventory.json, access_policy_inventory.json

Usage:
  cft scorecard [flags]


Flags:
      --bucket string                GCS bucket name for storing inventory (conflicts with --dir-path or --stdin)
      --workers int                  Concurrent Violations Review. If set, the CFT application will run the violations review concurrently and may improve the total execution time of the application. Default number of worker(s) is set to 1.
      --dir-path string              Local directory path for storing inventory (conflicts with --bucket or --stdin)
  -h, --help                         help for scorecard
      --output-format string         Format of scorecard outputs, can be txt, json or csv, default is txt
      --output-metadata-fields strings      List of comma delimited violation metadata fields of string type to include in output. Works when --output-format is txt or csv. By default no metadata fields in output when --output-format is txt or csv. All metadata will be in output when output-format is json.
      --output-path string           Path to directory to contain scorecard outputs. Output to console if not specified
      --policy-path string           Path to directory containing validation policies
      --refresh                      Refresh Cloud Asset Inventory export files in GCS bucket. If set, Application Default Credentials must be a service account (Works with --bucket)
      --stdin                        Passed Cloud Asset Inventory json string as standard input (conflicts with --dir-path or --bucket)
      --target-folder string         Folder ID to analyze (Works with --bucket and --refresh; conflicts with --target-project or --target--organization)
      --target-organization string   Organization ID to analyze (Works with --bucket and --refresh; conflicts with --target-project or --target--folder)
      --target-project string        Project ID to analyze (Works with --bucket and --refresh; conflicts with --target-folder or --target--organization)

Global Flags:
      --verbose   Log output to stdout

```

## Reporting
The CFT CLI can also be used for generating resource reports from CAI output files.
These resource reports are defined in Rego, including [these samples](../../reports/sample)

For example:

```bash
./cft report --query-path <path_to_cloud-foundation-toolkit>/reports/sample \
    --dir-path <path-to-directory-containing-cai-export> \
    --output-path <path-to-directory-for-report-output>
```

You could reuse the same CAI export generated for Scorecard by following these steps:
1. Download the report library from GitHub:
    ```
    git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
    ```
2. Create a directory to store report output:
    ```
    mkdir reports
    ```
3. Run the CFT report command:
    ```
    ./cft report --query-path cloud-foundation-toolkit/reports/sample \
        --dir-path ./cai-dir \
        --output-path ./reports
    ```
