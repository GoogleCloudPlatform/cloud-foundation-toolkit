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
    export REGION="us-east1"
fi

########## HELPER FUNCTIONS ##########

function create_config() {
    echo "Creating ${CONFIG}"
    envsubst < "templates/resource_policy/tests/integration/${TEST_NAME}.yaml" > "${CONFIG}"
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
    fi

  # Per-test teardown steps.
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config "${CONFIG}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
}

@test "Verifying that test-res-policy-inst-${RAND} was created in deployment ${DEPLOYMENT_NAME}" {
    run gcloud compute resource-policies list \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
        --filter="region:( ${REGION} )"
    [[ "$output" =~ "test-res-policy-inst-${RAND}" ]]
}

@test "Verifying resource policy test-res-policy-inst-${RAND}" {
    run gcloud compute resource-policies describe test-res-policy-inst-${RAND} \
        --region="${REGION}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "day: MONDAY" ]]
    [[ "$output" =~ "startTime: 00:00" ]]
}

@test "Deleting deployment" {
    gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"

    # Due to deployment does not delete resource policy it needs to be removed via CLI tool
    gcloud compute resource-policies delete test-res-policy-inst-${RAND} \
        --region="${REGION}"

    run gcloud compute resource-policies list \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
        --filter="region:( ${REGION} )"
    [[ ! "$output" =~ "test-res-policy-inst-${RAND}" ]]
}
