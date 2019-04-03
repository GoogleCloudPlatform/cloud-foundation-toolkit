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
    # Test specific variables.
    export RES_NAME="url-map-${RAND}"
    export NAME="url-map-name-${RAND}"
    export DESCRIPTION="url-map-description"
    export BACKEND_SERVICE_NAME="external-backend-service-${RAND}"
    export IGM_NAME="zonal-igm-http-${RAND}"
    export IT_NAME="instance-template-${RAND}"
    export HC_NAME="test-healthcheck-http-test"
    export PORT="80"
    export HOST="example.com"
    export PATH1="/audio"
    export PATH2="/video"
    export MATCHER_NAME="default-matcher"
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
    # Global setup; executed once per test file.
    if [ ${BATS_TEST_NUMBER} -eq 1 ]; then
        create_config
    fi

    # Per-test setup steps.
}

function teardown() {
    # Global teardown; executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        delete_config
    fi

    # Per-test teardown steps.
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    run gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config ${CONFIG} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

@test "Verifying URL map properties" {
    run gcloud compute url-maps describe "${NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "description: ${DESCRIPTION}" ]]
}

@test "Verifying URL map default backend" {
    run gcloud compute url-maps describe "${NAME}" \
        --format "yaml(defaultService)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${BACKEND_SERVICE_NAME}" ]]
}

@test "Verifying path matcher" {
    run gcloud compute url-maps describe "${NAME}" \
        --format "yaml(pathMatchers)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "name: ${MATCHER_NAME}" ]]
    [[ "$output" =~ "${BACKEND_SERVICE_NAME}" ]]
}

@test "Verifying path matcher paths" {
    run gcloud compute url-maps describe "${NAME}" \
        --format "yaml(pathMatchers)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "- ${PATH1}" ]]
    [[ "$output" =~ "- ${PATH2}" ]]
    [[ "$output" =~ "${BACKEND_SERVICE_NAME}" ]]
}

@test "Deleting deployment" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

