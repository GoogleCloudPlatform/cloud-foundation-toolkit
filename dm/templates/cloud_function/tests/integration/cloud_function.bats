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
    # Replace underscores with dashes in the deployment name.
    DEPLOYMENT_NAME=${DEPLOYMENT_NAME//_/-}
    CONFIG=".${DEPLOYMENT_NAME}.yaml"
fi

########## HELPER FUNCTIONS ##########

function create_config() {
    echo "Creating ${CONFIG}"
    envsubst < "templates/cloud_function/tests/integration/${TEST_NAME}.yaml" > "${CONFIG}"
}

function delete_config() {
    echo "Deleting ${CONFIG}"
    rm -f "${CONFIG}"
}

function setup() {
    # Global setup; executed once per test file.
    if [ ${BATS_TEST_NUMBER} -eq 1 ]; then
        gcloud pubsub topics create topic-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}"
        create_config
    fi

  # Per-test setup steps.
}

function teardown() {
    # Global teardown; executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        gcloud pubsub topics delete topic-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}"
        gsutil rm -r gs://test-function-http-${RAND}
        delete_config
    fi

    # Per-test teardown steps.
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config ${CONFIG} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
}

@test "Verifying that cloud functions were created in deployment ${DEPLOYMENT_NAME}" {
    run gcloud functions list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "test-function-https-name-${RAND}" ]]
    [[ "$output" =~ "test-function-storage-${RAND}" ]]
    [[ "$output" =~ "test-function-topic-${RAND}" ]]
}

@test "Verifying that test-function-https-name-${RAND} properties are set" {
    run gcloud functions describe test-function-https-name-${RAND} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "availableMemoryMb: 512" ]]
    [[ "$output" =~ "timeout: 120s" ]]
    [[ "$output" =~ "sourceArchiveUrl: gs://test-function-http-${RAND}/helloGET.zip" ]]
}

@test "Verifying that test-function-https-name-${RAND} trigger is set" {
    run gcloud functions describe test-function-https-name-${RAND} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "url: https://us-central1-${CLOUD_FOUNDATION_PROJECT_ID}.cloudfunctions.net/test-function-https-name-${RAND}" ]]
}

@test "Verifying that test-function-topic-${RAND} trigger is set" {
    run gcloud functions describe test-function-topic-${RAND} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "providers/cloud.pubsub/eventTypes/topic.publish" ]]
    [[ "$output" =~ "projects/${CLOUD_FOUNDATION_PROJECT_ID}/topics/topic-${RAND}" ]]
}

@test "Verifying that test-function-storage-${RAND} trigger is set" {
    run gcloud functions describe test-function-storage-${RAND} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "eventType: google.storage.object.finalize" ]]
    [[ "$output" =~ "resource: projects/${CLOUD_FOUNDATION_PROJECT_ID}/buckets/test-function-http-${RAND}" ]]
}

@test "Deleting deployment" {
    gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"

    run gcloud functions list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ ! "$output" =~ "test-function-https-name-${RAND}" ]]
    [[ ! "$output" =~ "test-function-storage-${RAND}" ]]
    [[ ! "$output" =~ "test-function-topic-${RAND}" ]]
}
