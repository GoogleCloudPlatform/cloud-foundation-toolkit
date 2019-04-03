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
    export ILB_RES_NAME="internal-load-balancer-${RAND}"
    export ILB_NAME="internal-load-balancer-name-${RAND}"
    export ILB_DESCRIPTION="ILB Description"
    export PROTOCOL="TCP"
    export ILB_PORT="80"
    export NETWORK_NAME="test-network-${RAND}"
    export BS_NAME="backend-service-name-${RAND}"
    export BS_DESCRIPTION="backend description"
    export BS_AFFINITY="CLIENT_IP"
    export BS_DRAINING="70"
    export TIMEOUT="40"
    export HC_NAME="tcp-healthcheck-${RAND}"
    export BACKEND_DESCRIPTION="instance group description"
    export IGM_NAME="regional-igm-${RAND}"
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

@test "Verifying forwarding rule" {
    run gcloud compute forwarding-rules describe "${ILB_NAME}" \
        --region ${REGION} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "IPProtocol: ${PROTOCOL}" ]]
    [[ "$output" =~ "${BS_NAME}" ]]
    [[ "$output" =~ "loadBalancingScheme: INTERNAL" ]]
    [[ "$output" =~ "name: ${ILB_NAME}" ]]
    [[ "$output" =~ "- '${ILB_PORT}'" ]]
    [[ "$output" =~ "${ILB_DESCRIPTION}" ]]
    [[ "$output" =~ "${NETWORK_NAME}" ]]
}

@test "Verifying backend service" {
    run gcloud compute backend-services describe "${BS_NAME}" \
        --region ${REGION} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "description: ${BS_DESCRIPTION}" ]]
    [[ "$output" =~ "protocol: ${PROTOCOL}" ]]
    [[ "$output" =~ "loadBalancingScheme: INTERNAL" ]]
    [[ "$output" =~ "sessionAffinity: ${BS_AFFINITY}" ]]
    [[ "$output" =~ "timeoutSec: ${TIMEOUT}" ]]
    [[ "$output" =~ "${HC_NAME}" ]]
    [[ "$output" =~ "drainingTimeoutSec: ${BS_DRAINING}" ]]
}

@test "Verifying backend" {
    run gcloud compute backend-services describe "${BS_NAME}" \
        --format "yaml(backends[0])" --region ${REGION} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "description: ${BACKEND_DESCRIPTION}" ]]
    [[ "$output" =~ "balancingMode: CONNECTION" ]]
    [[ "$output" =~ "${IGM_NAME}" ]]
}

@test "Deleting deployment" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}
