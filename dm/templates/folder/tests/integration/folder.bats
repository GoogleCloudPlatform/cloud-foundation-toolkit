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
fi

########## HELPER FUNCTIONS ##########

function create_config() {
    echo "Creating ${CONFIG}"
    envsubst < "templates/folder/tests/integration/${TEST_NAME}.yaml" > "${CONFIG}"
}

function delete_config() {
    echo "Deleting ${CONFIG}"
    rm -f "${CONFIG}"
}

function get_test_folder_id() {
        # Get the test folder ID and make it available.
        TEST_ORG_FOLDER_NAME=$(gcloud alpha resource-manager folders list \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --organization "${CLOUD_FOUNDATION_ORGANIZATION_ID}" | \
            grep "test-org-folder-${RAND}")

        export TEST_ORG_FOLDER_NAME=`echo ${TEST_ORG_FOLDER_NAME} | cut -d ' ' -f 3`
}

function setup() {
    # Global setup; this is executed once per test file.
    if [ ${BATS_TEST_NUMBER} -eq 1 ]; then
        gcloud alpha resource-manager folders create \
        --display-name="test-org-folder-${RAND}" \
        --organization="${CLOUD_FOUNDATION_ORGANIZATION_ID}"
        get_test_folder_id
        create_config
    fi

    # Per-test setup steps.
    get_test_folder_id
}

function teardown() {
    # Global teardown; this is executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        get_test_folder_id
        gcloud alpha resource-manager folders delete "${TEST_ORG_FOLDER_NAME}"
        rm -f "${RANDOM_FILE}"
        delete_config
    fi

    # Per-test teardown steps.
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
  gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
    --config ${CONFIG} \
    --project "${CLOUD_FOUNDATION_PROJECT_ID}"
}

@test "Verifying that a folder was created under organization in deployment ${DEPLOYMENT_NAME}" {
  run gcloud alpha resource-manager folders list \
    --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
    --organization "${CLOUD_FOUNDATION_ORGANIZATION_ID}"
  [[ "$output" =~ "Folder under Org ${RAND}" ]]
}

@test "Verifying that a folder was created under the specified folder in deployment ${DEPLOYMENT_NAME}" {
  run gcloud alpha resource-manager folders list \
    --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
    --folder "${TEST_ORG_FOLDER_NAME}"
  [[ "$output" =~ "Folder under Folder ${RAND}" ]]
}

@test "Deleting deployment" {
  gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" \
    --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q

  run gcloud  run gcloud alpha resource-manager folders list \
    --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
    --organization "${CLOUD_FOUNDATION_ORGANIZATION_ID}"
  [[ ! "$output" =~ "Folder Under Org ${RAND}" ]]

  run gcloud alpha resource-manager folders list \
    --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
    --folder "${TEST_ORG_FOLDER_NAME}"
  [[ ! "$output" =~ "Folder Under Folder ${RAND}" ]]
}
