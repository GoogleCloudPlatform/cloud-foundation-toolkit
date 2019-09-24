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
    export UMIG_NAME="umig-${RAND}"
    export UMIG_RES_NAME="umig-${RAND}"
    export ZONE="us-central1-c"
    export PORT_NAME="http"
    export PORT="80"
    export INSTANCE_NAME="test-umig-instance-${RAND}"
fi

########## HELPER FUNCTIONS ##########

function create_config() {
    echo "Creating ${CONFIG}"
    envsubst < templates/unmanaged_instance_group/tests/integration/${TEST_NAME}.yaml > "${CONFIG}"
}

function delete_config() {
    rm -f "${CONFIG}"
}

function setup() {
    # Global setup; executed once per test file.
    if [ ${BATS_TEST_NUMBER} -eq 1 ]; then
        create_config
        # Needed for testing resource creation with preexisting (not referenced)
        # instance
        gcloud compute instances create "${INSTANCE_NAME}" \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --zone "${ZONE}"
    fi

    # Per-test setup steps.
}

function teardown() {
    # Global teardown; executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        rm -f "${RANDOM_FILE}"
        gcloud compute instances delete "${INSTANCE_NAME}" \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --zone "${ZONE}" -q
        delete_config
    fi

    # Per-test teardown steps.
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    run gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config "${CONFIG}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    echo "$output"
    [[ "$status" -eq 0 ]]
}

@test "Verifying that unmanaged intance group was created" {
    run gcloud compute instance-groups unmanaged list \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    echo "$output"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${UMIG_NAME}" ]]
    [[ "$output" =~ "${ZONE}" ]]
}

@test "Verifying unmanaged instance group properties" {
    run gcloud compute instance-groups unmanaged describe "${UMIG_NAME}" \
        --zone "${ZONE}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "instanceGroups/${UMIG_NAME}" ]]
    [[ "$output" =~ "name: ${PORT_NAME}" ]]
    [[ "$output" =~ "port: ${PORT}" ]]
    [[ "$output" =~ "size: 1" ]]
}

@test "Deleting deployment" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q
    [[ "$status" -eq 0 ]]
}

