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
    export NETWORK="testnw1-${RAND}"
    export PEER_NETWORK="testnw2--${RAND}"
    export PEER_NAME="test-peer"
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
        # create sample networks
        gcloud compute networks create ${NETWORK} --subnet-mode custom
        gcloud compute networks create ${PEER_NETWORK} --subnet-mode custom
    fi

    # Per-test setup steps.
}

function teardown() {
    # Global teardown; executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        delete_config
        rm -f "${RANDOM_FILE}"
        # delete sample networks
        gcloud compute networks delete ${NETWORK} ${PEER_NETWORK} -q
    fi

    # Per-test teardown steps.
}

########## TESTS ##########

@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    run gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config "${CONFIG}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

@test "Verify if peering ${PEER_NAME} is created " {
    run gcloud compute networks peerings list --network ${NETWORK} \
        --format="value(peerings[0].name)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${PEER_NAME}" ]]
}

@test "Verify if peer network in the PEER is ${PEER_NETWORK} " {
    run gcloud compute networks peerings list --network ${NETWORK} \
        --format="value(peerings[0].network)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${PEER_NETWORK}" ]]
}

@test "Verify if peer status is INACTIVE " {
    run gcloud compute networks peerings list --network ${NETWORK} \
        --format="value(peerings[0].status)"
    [[ "$status" -eq 0 ]]
    [[ "$output" -eq "INACTIVE" ]]
}

@test "Deleting deployment ${DEPLOYMENT_NAME}" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q
    [[ "$status" -eq 0 ]]
}
