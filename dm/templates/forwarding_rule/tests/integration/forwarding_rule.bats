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
    # test specific variables
    export REGION="us-central1"
    export EXTERNAL_RES_NAME="external-global-fr-${RAND}"
    export INTERNAL_RES_NAME="internal-regional-fr-${RAND}"
    export PROXY_NAME="http-proxy-${RAND}"
    export EXTERNAL_LB_SCHEME="EXTERNAL"
    export EXTERNAL_PORT="80"
    export INTERNAL_NAME="fr-internal-regional-${RAND}"
    export INTERNAL_DESC="Internal description"
    export INTERNAL_PORT="80"
    export INTERNAL_LB_SCHEME="INTERNAL"
    export ZONE="us-central1-f"
fi

########## HELPER FUNCTIONS ##########

function create_config() {
    echo "Creating ${CONFIG}"
    envsubst < "templates/forwarding_rule/tests/integration/${TEST_NAME}.yaml" > "${CONFIG}"
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

@test "Verifying global external forwarding rule" {
    TARGET_PROXY="global/targetHttpProxies/${PROXY_NAME}"
    run gcloud compute forwarding-rules describe \
        "${EXTERNAL_RES_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
        --global
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "portRange: ${EXTERNAL_PORT}-${EXTERNAL_PORT}" ]]
    [[ "$output" =~ "loadBalancingScheme: ${EXTERNAL_LB_SCHEME}" ]]
    [[ "$output" =~ "$TARGET_PROXY" ]]
}

@test "Verifying regional internal forwarding rule" {
    run gcloud compute forwarding-rules describe "${INTERNAL_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
        --region "${REGION}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "- '${INTERNAL_PORT}'" ]]
    [[ "$output" =~ "description: ${INTERNAL_DESC}" ]]
    [[ "$output" =~ "name: ${INTERNAL_NAME}" ]]
    [[ "$output" =~ "loadBalancingScheme: ${INTERNAL_LB_SCHEME}" ]]
    [[ "$output" =~ "regional-internal-backend-service-${RAND}" ]]
}

@test "Deleting deployment" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}
