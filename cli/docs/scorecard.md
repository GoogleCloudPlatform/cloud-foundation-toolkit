# CFT Scorecard

The CFT Scorecard is integrated into the [CFT CLI](../README.md) and provides
an easy integration with [Forseti Config Validator](https://github.com/forseti-security/policy-library/blob/master/docs/user_guide.md).
It can be used to print a scorecard of your GCP environment, for resources and IAM policies in Cloud Asset Inventory (CAI) exports.
The policies tested are based on constraints and constraint templates from the [Config Validator policy library](https://github.com/forseti-security/policy-library).

## Scorecard User Guide
This tutorial will walk you through setting up Scorecard for a single project.

1. Set some environment variables:
    ```
    export GOOGLE_PROJECT=my-cai-project             # For using CAI API exporting CAI data to GCS
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
5. Export the CAI data to GCS:
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
6. Download the CFT CLI and make it executable:
    ```
    # OS X
    curl -o cft https://storage.googleapis.com/cft-cli/latest/cft-darwin-amd64
    # Linux
    curl -o cft https://storage.googleapis.com/cft-cli/latest/cft-linux-amd64
    # executable
    chmod +x cft
    ```
7. Download the sample policy library and add a sample constraint for detecting public buckets:
    ```
    git clone https://github.com/forseti-security/policy-library.git
    cp policy-library/samples/storage_blacklist_public.yaml policy-library/policies/constraints/
    ```
8. Optionally, create a public GCS bucket to trigger a violation:
    ```
    gsutil mb gs://$PUBLIC_BUCKET_NAME
    gsutil iam ch allUsers:objectViewer gs://$PUBLIC_BUCKET_NAME
    ```
9. Run CFT Scorecard:
    ```
    ./cft scorecard --policy-path=./policy-library/ \
        --bucket=$CAI_BUCKET_NAME
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
