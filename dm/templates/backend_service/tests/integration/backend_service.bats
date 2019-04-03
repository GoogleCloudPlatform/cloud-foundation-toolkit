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
    # test specific variables
    export REGION="us-east1"
    export RES_DESCRIPTION="This is the description"
    export BC_DESCRIPTION="Backend description"
    export TIMEOUT="35"
    export ENABLE_CDN="true"
    export SESSION="CLIENT_IP"
    export REGIONAL_BALANCING_MODE="CONNECTION"
    export GLOBAL_BALANCING_MODE="RATE"
    export REGIONAL_BALANCING_SCHEME="INTERNAL"
    export REGIONAL_BALANCING_PROTOCOL="TCP"
    export GLOBAL_BALANCING_SCHEME="EXTERNAL"
    export GLOBAL_BALANCING_PROTOCOL="HTTP"
    export MAX_RATE="10000"
fi

########## HELPER FUNCTIONS ##########

function create_config() {
    echo "Creating ${CONFIG}"
    envsubst < "templates/backend_service/tests/integration/${TEST_NAME}.yaml" > "${CONFIG}"
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

@test "Verifying global external backend service" {
    run gcloud compute backend-services describe \
        "global-external-backend-service-${RAND}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
        --global
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "balancingMode: ${GLOBAL_BALANCING_MODE}" ]]
    [[ "$output" =~ "maxRate: ${MAX_RATE}" ]]
    [[ "$output" =~ "loadBalancingScheme: ${GLOBAL_BALANCING_SCHEME}" ]]
    [[ "$output" =~ "protocol: ${GLOBAL_BALANCING_PROTOCOL}" ]]
    [[ "$output" =~ "enableCDN: ${ENABLE_CDN}" ]]
}

@test "Verifying regional internal backend service" {
    run gcloud compute backend-services describe \
        "regional-internal-backend-service-${RAND}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
        --region "${REGION}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "balancingMode: ${REGIONAL_BALANCING_MODE}" ]]
    [[ "$output" =~ "loadBalancingScheme: ${REGIONAL_BALANCING_SCHEME}" ]]
    [[ "$output" =~ "protocol: ${REGIONAL_BALANCING_PROTOCOL}" ]]
    [[ "$output" =~ "  description: ${BC_DESCRIPTION}" ]]
    [[ "$output" =~ "description: ${RES_DESCRIPTION}" ]]
    [[ "$output" =~ "sessionAffinity: ${SESSION}" ]]
    [[ "$output" =~ "timeoutSec: ${TIMEOUT}" ]]
}

@test "Deleting deployment" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}
