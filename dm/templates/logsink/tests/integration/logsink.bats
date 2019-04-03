#!/usr/bin/env bats

source tests/helpers.bash

TEST_NAME=$(basename "${BATS_TEST_FILENAME}" | cut -d '.' -f 1)

# Create a random 10-char string and save it in a file.
RANDOM_FILE="/tmp/${CLOUD_FOUNDATION_ORGANIZATION_ID}-${TEST_NAME}.txt"
if [[ ! -e "${RANDOM_FILE}" ]]; then
    RAND=$(head /dev/urandom | LC_ALL=C tr -dc a-z0-9 | head -c 10)
    echo ${RAND} > "${RANDOM_FILE}"
fi

# Set variables based on the random string saved in the file.
# envsubst requires all variables used in the example/config to be exported.
if [[ -e "${RANDOM_FILE}" ]]; then
    export RAND=$(cat "${RANDOM_FILE}")
    DEPLOYMENT_NAME="${CLOUD_FOUNDATION_PROJECT_ID}-${TEST_NAME}-${RAND}"
    # Replace underscores in the deployment name with dashes.
    DEPLOYMENT_NAME=${DEPLOYMENT_NAME//_/-}
    CONFIG=".${DEPLOYMENT_NAME}.yaml"
fi

########## HELPER FUNCTIONS ##########

function create_config() {
    echo "Creating ${CONFIG}"
    envsubst < ${BATS_TEST_DIRNAME}/${TEST_NAME}.yaml > "${CONFIG}"
}

function delete_config() {
    echo "Deleting ${CONFIG}"
    rm -f "${CONFIG}"
}

function get_test_folder_id() {
    # Get the test folder ID and make it available
    TEST_ORG_FOLDER_NAME=$(gcloud alpha resource-manager folders list \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
        --organization "${CLOUD_FOUNDATION_ORGANIZATION_ID}" | \
        grep "test-org-folder-${RAND}")

    export TEST_ORG_FOLDER_NAME=`echo ${TEST_ORG_FOLDER_NAME} | cut -d ' ' -f 3`
}

function setup() {
    # Global setup; this is executed once per test file.
    if [ ${BATS_TEST_NUMBER} -eq 1 ]; then
        gcloud alpha resource-manager folders create \
            --display-name="test-org-folder-${RAND}" \
            --organization="${CLOUD_FOUNDATION_ORGANIZATION_ID}"
        get_test_folder_id
        create_config
        gcloud pubsub topics create test-topic-${RAND}
        gsutil mb -l us-east1 gs://test-bucket-${RAND}/
        bq mk test_dataset_${RAND}
    fi

    # Per-test setup as per documentation.
    get_test_folder_id
}

function teardown() {
    # Global teardown; this is executed once per test file
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        get_test_folder_id
        gcloud alpha resource-manager folders delete "${TEST_ORG_FOLDER_NAME}"
        gsutil rm -r gs://test-bucket-${RAND}/
        gcloud pubsub topics delete test-topic-${RAND}
        bq rm -rf test_dataset_${RAND}
        delete_config
        rm -f "${RANDOM_FILE}"
    fi

    # Per-test teardown as per documentation.
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    run gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config "${CONFIG}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

@test "Verifying project sinks were created each with the requested destination in deployment ${DEPLOYMENT_NAME}" {
    run gcloud logging sinks list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "test-logsink-project-bq-${RAND}" ]]
    [[ "$output" =~ "test-logsink-project-pubsub-${RAND}" ]]
    [[ "$output" =~ "test-logsink-project-storage-${RAND}" ]]
}

@test "Verifying organization sinks were created each with a different as the destination in deployment ${DEPLOYMENT_NAME}" {
    run gcloud logging sinks list \
        --organization "${CLOUD_FOUNDATION_ORGANIZATION_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "test-logsink-org-bq-${RAND}" ]]
    [[ "$output" =~ "test-logsink-org-pubsub-${RAND}" ]]
    [[ "$output" =~ "test-logsink-org-storage-${RAND}" ]]
}

@test "Verifying billing account sinks were created each with a different as the destination in deployment ${DEPLOYMENT_NAME}" {
    run gcloud logging sinks list --billing-account \
        "${CLOUD_FOUNDATION_BILLING_ACCOUNT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "test-logsink-billing-bq-${RAND}" ]]
    [[ "$output" =~ "test-logsink-billing-pubsub-${RAND}" ]]
    [[ "$output" =~ "test-logsink-billing-storage-${RAND}" ]]
}

@test "Verifying folder sinks were created each with a different as the destination in deployment ${DEPLOYMENT_NAME}" {
    run gcloud logging sinks list --folder "${TEST_ORG_FOLDER_NAME}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "test-logsink-folder-bq-${RAND}" ]]
    [[ "$output" =~ "test-logsink-folder-pubsub-${RAND}" ]]
    [[ "$output" =~ "test-logsink-folder-storage-${RAND}" ]]
}

@test "Verifying project sinks and the destination resource were created in deployment ${DEPLOYMENT_NAME}" {
    run gcloud logging sinks list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    #[[ "$output" =~ "test-logsink-project-bq-${RAND}" ]]
    [[ "$output" =~ "test-logsink-project-pubsub-create-${RAND}" ]]
    [[ "$output" =~ "test-logsink-project-storage-create-${RAND}" ]]

    run gcloud beta pubsub topics get-iam-policy \
        "test-logsink-project-pubsub-topic-dest-${RAND}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "@gcp-sa-logging.iam.gserviceaccount.com" ]]
    [[ "$output" =~ "user:${CLOUD_FOUNDATION_USER_ACCOUNT}" ]]
    [[ "$output" =~ "role: roles/pubsub.admin" ]]

    run gsutil iam get "gs://test-logsink-project-storage-dest-${RAND}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "@gcp-sa-logging.iam.gserviceaccount.com" ]]
    [[ "$output" =~ "roles/storage.admin" ]]
    [[ "$output" =~ "user:${CLOUD_FOUNDATION_USER_ACCOUNT}" ]]
    [[ "$output" =~ "roles/storage.objectViewer" ]]

    #TODO: Add test for BQ
}

@test "Verifying org sinks and the destination resource were created in deployment ${DEPLOYMENT_NAME}" {
    run gcloud logging sinks list \
        --organization "${CLOUD_FOUNDATION_ORGANIZATION_ID}"
    [[ "$status" -eq 0 ]]
    #[[ "$output" =~ "test-logsink-org-bq-${RAND}" ]]
    [[ "$output" =~ "test-logsink-org-pubsub-create-${RAND}" ]]
    [[ "$output" =~ "test-logsink-org-storage-create-${RAND}" ]]

    run gcloud beta pubsub topics get-iam-policy "test-logsink-org-pubsub-topic-dest-${RAND}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "@gcp-sa-logging.iam.gserviceaccount.com" ]]
    [[ "$output" =~ "user:${CLOUD_FOUNDATION_USER_ACCOUNT}" ]]
    [[ "$output" =~ "role: roles/pubsub.admin" ]]

    run gsutil iam get "gs://test-logsink-org-storage-dest-${RAND}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "@gcp-sa-logging.iam.gserviceaccount.com" ]]
    [[ "$output" =~ "roles/storage.admin" ]]
    [[ "$output" =~ "user:${CLOUD_FOUNDATION_USER_ACCOUNT}" ]]
    [[ "$output" =~ "roles/storage.objectViewer" ]]

    #TODO: Add test for BQ
}

@test "Verifying billing sinks and the destination resource were created in deployment ${DEPLOYMENT_NAME}" {
    run gcloud logging sinks list \
        --billing-account "${CLOUD_FOUNDATION_BILLING_ACCOUNT_ID}"
    [[ "$status" -eq 0 ]]
    #[[ "$output" =~ "test-logsink-billing-bq-${RAND}" ]]
    [[ "$output" =~ "test-logsink-billing-pubsub-create-${RAND}" ]]
    [[ "$output" =~ "test-logsink-billing-storage-create-${RAND}" ]]

    run gcloud beta pubsub topics get-iam-policy \
        "test-logsink-billing-pubsub-topic-dest-${RAND}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "@gcp-sa-logging.iam.gserviceaccount.com" ]]
    [[ "$output" =~ "user:${CLOUD_FOUNDATION_USER_ACCOUNT}" ]]
    [[ "$output" =~ "role: roles/pubsub.admin" ]]

    run gsutil iam get "gs://test-logsink-billing-storage-dest-${RAND}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "@gcp-sa-logging.iam.gserviceaccount.com" ]]
    [[ "$output" =~ "roles/storage.admin" ]]
    [[ "$output" =~ "user:${CLOUD_FOUNDATION_USER_ACCOUNT}" ]]
    [[ "$output" =~ "roles/storage.objectViewer" ]]

    #TODO: Add test for BQ
}

@test "Verifying folder sinks and the destination resource were created in deployment ${DEPLOYMENT_NAME}" {
    run gcloud logging sinks list --folder "${TEST_ORG_FOLDER_NAME}"
    [[ "$status" -eq 0 ]]
    #[[ "$output" =~ "test-logsink-folder-bq-${RAND}" ]]
    [[ "$output" =~ "test-logsink-folder-pubsub-create-${RAND}" ]]
    [[ "$output" =~ "test-logsink-folder-storage-create-${RAND}" ]]

    run gcloud beta pubsub topics get-iam-policy \
        "test-logsink-folder-pubsub-topic-dest-${RAND}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "@gcp-sa-logging.iam.gserviceaccount.com" ]]
    [[ "$output" =~ "user:${CLOUD_FOUNDATION_USER_ACCOUNT}" ]]
    [[ "$output" =~ "role: roles/pubsub.admin" ]]

    run gsutil iam get "gs://test-logsink-folder-storage-dest-${RAND}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "@gcp-sa-logging.iam.gserviceaccount.com" ]]
    [[ "$output" =~ "roles/storage.admin" ]]
    [[ "$output" =~ "user:${CLOUD_FOUNDATION_USER_ACCOUNT}" ]]
    [[ "$output" =~ "roles/storage.objectViewer" ]]

    #TODO: Add test for BQ
}

@test "Deleting deployment" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]

    run gcloud logging sinks list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    #[[ ! "$output" =~ "test-logsink-project-bq-${RAND}" ]]
    [[ ! "$output" =~ "test-logsink-project-pubsub-${RAND}" ]]
    [[ ! "$output" =~ "test-logsink-project-storage-${RAND}" ]]

    run gcloud logging sinks list \
        --organization "${CLOUD_FOUNDATION_ORGANIZATION_ID}"
    [[ "$status" -eq 0 ]]
    #[[ ! "$output" =~ "test-logsink-org-bq-${RAND}" ]]
    [[ ! "$output" =~ "test-logsink-org-pubsub-${RAND}" ]]
    [[ ! "$output" =~ "test-logsink-org-storage-${RAND}" ]]

    # TODO: Bug where billing accounts are not deleted during deployment delete.
    #       Re-enable this check once its fixed.
    #run gcloud logging sinks list --billing-account \
    #    "${CLOUD_FOUNDATION_BILLING_ACCOUNT_ID}"
    #[[ "$status" -eq 0 ]]
    #[[ ! "$output" =~ "test-logsink-billing-bq-${RAND}" ]]
    #[[ ! "$output" =~ "test-logsink-billing-pubsub-${RAND}" ]]
    #[[ ! "$output" =~ "test-logsink-billing-storage-${RAND}" ]]

    run gcloud logging sinks list --folder "${TEST_ORG_FOLDER_NAME}"
    [[ "$status" -eq 0 ]]
    #[[ ! "$output" =~ "test-logsink-folder-bq-${RAND}" ]]
    [[ ! "$output" =~ "test-logsink-folder-pubsub-${RAND}" ]]
    [[ ! "$output" =~ "test-logsink-folder-storage-${RAND}" ]]
}
