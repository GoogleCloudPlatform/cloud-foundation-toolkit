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
    export CLOUDDNS_ZONE_NAME="test-managed-zone-${RAND}"
    export CLOUDDNS_DNS_NAME="${RAND}.com."
    export CLOUDDNS_DESCRIPTION="Managed DNS Zone for Testing"
fi

########## HELPER FUNCTIONS ##########

function create_config() {
    echo "Creating ${CONFIG}"
    envsubst < templates/dns_managed_zone/tests/integration/${TEST_NAME}.yaml > "${CONFIG}"
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


########## TESTS ##########

@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
   gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
       --config "${CONFIG}" --project "${CLOUD_FOUNDATION_PROJECT_ID}"
   [[ "$status" -eq 0 ]]
}

@test "Verify if a managed zone with name $CLOUDDNS_ZONE_NAME was created" {
   run gcloud dns managed-zones list --format=flattened
   [[ "$status" -eq 0 ]]
   [[ "$output" =~ "${CLOUDDNS_ZONE_NAME}" ]]
}

@test "Verify if a DNS named ${CLOUDDNS_DNS_NAME} was created" {
    run gcloud dns managed-zones list
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${CLOUDDNS_DNS_NAME}" ]]
}

@test "Deleting deployment ${DEPLOYMENT_NAME}" {
    gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" \
        -q --project "${CLOUD_FOUNDATION_PROJECT_ID}"

    run gcloud dns managed-zones list
    [[ "$status" -eq 0 ]]
    [[ ! "$output" =~ "${CLOUDDNS_ZONE_NAME}" ]]
}

