#!/usr/bin/env bats

source tests/helpers.bash

TEST_NAME=$(basename "${BATS_TEST_FILENAME}" | cut -d '.' -f 1)

export TEST_SERVICE_ACCOUNT="test-sa-${RAND}"

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
fi

########## HELPER FUNCTIONS ##########

function create_config() {
    echo "Creating ${CONFIG}"
    envsubst < "templates/iam_member/tests/integration/${TEST_NAME}.yaml" > "${CONFIG}"
}

function delete_config() {
    echo "Deleting ${CONFIG}"
    rm -f "${CONFIG}"
}

function setup() {
    # Global setup; this is executed once per test file.
    if [ ${BATS_TEST_NUMBER} -eq 1 ]; then
        create_config
        gcloud iam service-accounts create "${TEST_SERVICE_ACCOUNT}" \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    fi

    # Per-test setup steps.
}

function teardown() {
    # Global teardown; this is executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        gcloud iam service-accounts delete "${TEST_SERVICE_ACCOUNT}@${CLOUD_FOUNDATION_PROJECT_ID}.iam.gserviceaccount.com" \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}"
        delete_config
        rm -f "${RANDOM_FILE}"
    fi

    # Per-test teardown steps.
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config "${CONFIG}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
}

@test "Verifying that roles were assigned in deployment ${DEPLOYMENT_NAME}" {
    run gcloud projects get-iam-policy "${CLOUD_FOUNDATION_PROJECT_ID}" \
        --flatten="bindings[].members" \
        --format='table(bindings.role)' \
        --filter="bindings.members:${TEST_SERVICE_ACCOUNT}@${CLOUD_FOUNDATION_PROJECT_ID}.iam.gserviceaccount.com"
    [[ "$output" =~ "roles/editor" ]]
    [[ "$output" =~ "roles/viewer" ]]
}

@test "Deleting deployment" {
    gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q

    run gcloud projects get-iam-policy "${CLOUD_FOUNDATION_PROJECT_ID}" \
        --flatten="bindings[].members" \
        --format='table(bindings.role)' \
        --filter="bindings.members:${TEST_SERVICE_ACCOUNT}@${CLOUD_FOUNDATION_PROJECT_ID}.iam.gserviceaccount.com"
    [[ ! "$output" =~ "roles/editor" ]]
    [[ ! "$output" =~ "roles/viewer" ]]
}
