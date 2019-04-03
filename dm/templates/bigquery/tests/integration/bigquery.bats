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
    export TEST_SERVICE_ACCOUNT="test-sa-${RAND}"
    export RAND=$(cat "${RANDOM_FILE}")
    DEPLOYMENT_NAME="${CLOUD_FOUNDATION_PROJECT_ID}-${TEST_NAME}-${RAND}"
    # Replace underscores with dashes in the deployment name.
    DEPLOYMENT_NAME=${DEPLOYMENT_NAME//_/-}
    CONFIG=".${DEPLOYMENT_NAME}.yaml"
fi

########## HELPER FUNCTIONS ##########

function create_config() {
    echo "Creating ${CONFIG}"
    envsubst < "${BATS_TEST_DIRNAME}/${TEST_NAME}.yaml" > "${CONFIG}"
}

function delete_config() {
    echo "Deleting ${CONFIG}"
    rm -f "${CONFIG}"
}

function setup() {
    # Global setup; executed once per test file.
    if [ ${BATS_TEST_NUMBER} -eq 1 ]; then
        gcloud iam service-accounts create "${TEST_SERVICE_ACCOUNT}" \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}"
        create_config
    fi

  # Per-test setup steps.
  }

function teardown() {
    # Global teardown; executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        gcloud iam service-accounts delete \
            "${TEST_SERVICE_ACCOUNT}@${CLOUD_FOUNDATION_PROJECT_ID}.iam.gserviceaccount.com" \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}"
        rm -f "${RANDOM_FILE}"
        delete_config
    fi

    # Per-test teardown steps
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config ${CONFIG} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
}

@test "Verifying that a dataset was created in deployment ${DEPLOYMENT_NAME}" {
    run bq show --format=prettyjson \
        "${CLOUD_FOUNDATION_PROJECT_ID}":test_bq_dataset_${RAND}
    [ "$status" -eq 0 ]
    [[ "$output" =~ "\"datasetId\": \"test_bq_dataset_${RAND}\"" ]]
}

@test "Verifying that a table was created in the dataset in deployment ${DEPLOYMENT_NAME}" {
    run bq ls --format=prettyjson \
        "${CLOUD_FOUNDATION_PROJECT_ID}":test_bq_dataset_${RAND}
    [ "$status" -eq 0 ]
    [[ "$output" =~ "\"tableId\": \"test_bq_table_${RAND}\"" ]]
}

@test "Verifying that a table schema was created in the dataset deployment ${DEPLOYMENT_NAME}" {
    run bq show --schema test_bq_dataset_${RAND}.test_bq_table_${RAND}
    [ "$status" -eq 0 ]
    [[ "$output" =~ "{\"type\":\"STRING\",\"name\":\"firstname\"}" ]]
    [[ "$output" =~ "{\"type\":\"STRING\",\"name\":\"lastname\"}" ]]
    [[ "$output" =~ "{\"type\":\"INTEGER\",\"name\":\"age\"}" ]]
}

@test "Deleting deployment" {
    gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q

    run bq show --format=prettyjson "${CLOUD_FOUNDATION_PROJECT_ID}":test_bq_dataset_${RAND}
    [[ ! "$output" =~ "\datasetId\": \"test_bq_dataset_${RAND}\"" ]]

    run bq ls --format=prettyjson "${CLOUD_FOUNDATION_PROJECT_ID}":test_bq_dataset_${RAND}
    [[ ! "$output" =~ "\"tableId\": \"test_bq_table_${RAND}\"" ]]
}
