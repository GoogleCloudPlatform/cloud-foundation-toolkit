#!/usr/bin/env bats

source tests/helpers.bash

TEST_NAME=$(basename "${BATS_TEST_FILENAME}" | cut -d '.' -f 1)

# Create a random 10-char string and save it in a file.
RANDOM_FILE="/tmp/${CLOUD_FOUNDATION_PROJECT_ID}-${TEST_NAME}.txt"
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

function setup() {
    # Global setup; this is executed once per test file.
    if [ ${BATS_TEST_NUMBER} -eq 1 ]; then
        create_config
    fi
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    run gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config "${CONFIG}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"

    echo "Status: $status"
    echo "Output: $output"

    [[ "$status" -eq 0 ]]
}

@test "Verifying project sinks were created each with the requested destination in deployment ${DEPLOYMENT_NAME}" {
    run gcloud logging sinks list --project "${CLOUD_FOUNDATION_PROJECT_ID}"

    echo "Status: $status"
    echo "Output: $output"

    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "test-logsink-project-bq-${RAND}" ]]
    [[ "$output" =~ "test-logsink-project-pubsub-${RAND}" ]]
    [[ "$output" =~ "test-logsink-project-storage-${RAND}" ]]

    run gcloud beta pubsub topics get-iam-policy \
        "projects/${CLOUD_FOUNDATION_PROJECT_ID}/topics/test-topic-${RAND}"

    echo "Status: $status"
    echo "Output: $output"

    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "@gcp-sa-logging.iam.gserviceaccount.com" ]]
    [[ "$output" =~ "role: roles/pubsub.admin" ]]

    run gsutil iam get "gs://test-bucket-${RAND}" --project "${TARGET_PROJECT_ID}"

    echo "Status: $status"
    echo "Output: $output"

    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "@gcp-sa-logging.iam.gserviceaccount.com" ]]
    [[ "$output" =~ "roles/storage.objectAdmin" ]]
}

@test "Deleting deployment" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"

    echo "Status: $status"
    echo "Output: $output"

    [[ "$status" -eq 0 ]]

    run gcloud logging sinks list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]

    echo "Status: $status"
    echo "Output: $output"

    [[ ! "$output" =~ "test-logsink-project-bq-${RAND}" ]]
    [[ ! "$output" =~ "test-logsink-project-pubsub-${RAND}" ]]
    [[ ! "$output" =~ "test-logsink-project-storage-${RAND}" ]]
}
