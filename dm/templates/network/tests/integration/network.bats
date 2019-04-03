#!/usr/bin/env bats

source tests/helpers.bash

# Create a random 10-char string and save it in a file.
RANDOM_FILE="/tmp/${CLOUD_FOUNDATION_ORGANIZATION_ID}-network.txt"
if [[ ! -e "${RANDOM_FILE}" ]]; then
    RAND=$(head /dev/urandom | LC_ALL=C tr -dc a-z0-9 | head -c 10)
    echo ${RAND} > "${RANDOM_FILE}"
fi

# Set variables based on the random string saved in the file.
# envsubst requires all variables used in the example/config to be exported.
if [[ -e "${RANDOM_FILE}" ]]; then
    export RAND=$(cat "${RANDOM_FILE}")
    DEPLOYMENT_NAME="${CLOUD_FOUNDATION_PROJECT_ID}-network-${RAND}"
    # Replace underscores in the deployment name with dashes.
    DEPLOYMENT_NAME=${DEPLOYMENT_NAME//_/-}
    CONFIG=".${DEPLOYMENT_NAME}.yaml"
fi

########## HELPER FUNCTIONS ##########

function create_config() {
    echo "Creating ${CONFIG}"
    envsubst < templates/network/tests/integration/network.yaml > "${CONFIG}"
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

    # Per-test setup steps.
}

function teardown() {
    # Global teardown; this is executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        delete_config
        rm -f "${RANDOM_FILE}"
    fi

    # Per-test teardown steps.
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" --config "${CONFIG}" --project "${CLOUD_FOUNDATION_PROJECT_ID}"
}

@test "Verifying that resources were created in deployment ${DEPLOYMENT_NAME}" {
    run gcloud compute networks list --filter="name:test-network-${RAND}" --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "test-network-${RAND}" ]]
}

@test "Verifying subnets were created in deployment ${DEPLOYMENT_NAME}" {
    run gcloud compute networks subnets list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "test-subnetwork-${RAND}-1" ]]
    [[ "$output" =~ "test-subnetwork-${RAND}-2" ]]
}

@test "Deleting deployment" {
    gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q --project "${CLOUD_FOUNDATION_PROJECT_ID}"

    run gcloud compute networks list --filter="name:test-network-${RAND}" --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ ! "$output" =~ "test-network-${RAND}" ]]

    run gcloud compute networks subnets list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ ! "$output" =~ "test-subnetwork-${RAND}-1" ]]
    [[ ! "$output" =~ "test-subnetwork-${RAND}-2" ]]
}
