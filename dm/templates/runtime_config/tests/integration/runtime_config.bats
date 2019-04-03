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
    # Test specific variables.
    export CONFIG_NAME="test-config-${RAND}"
    export VARIABLE_1="test/dev/db/connection_string"
    export VARIABLE_1_VALUE="Server=sqlsrv;Database=mydb;Uid=uname;Pwd=pwd;"
    export VARIABLE_2="test/dev/web/appvalue"
    # 'my test text value' in base64
    export VARIABLE_2_VALUE="bXkgdGVzdCB0ZXh0IHZhbHVl"
    export WAITER_NAME="test-waiter-${RAND}"
    export WAITER_TIMEOUT="2.500s"
    export WAITER_PATH="test/dev"
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
        rm -f "${RANDOM_FILE}"
    fi

    # Per-test teardown steps.
}

########## TESTS ##########

@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    run gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config "${CONFIG}" --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

@test "Verify if CONFIG ${CONFIG_NAME} is created " {
    run gcloud beta runtime-config configs list --format="value(name)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${CONFIG_NAME}" ]]
}

@test "Verify if VARIABLE ${VARIABLE_1} is created " {
    run gcloud beta runtime-config configs variables list \
        --config-name ${CONFIG_NAME} \
        --format="value(name)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${VARIABLE_1}" ]]
}

@test "Verify if VARIABLE ${VARIABLE_1} has value ${VARIABLE_1_VALUE} " {
    run gcloud beta runtime-config configs variables get-value ${VARIABLE_1} \
        --config-name ${CONFIG_NAME}
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${VARIABLE_1_VALUE}" ]]
}

@test "Verify if VARIABLE ${VARIABLE_2} is created " {
    run gcloud beta runtime-config configs variables list \
        --config-name ${CONFIG_NAME} \
        --format="value(name)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${VARIABLE_2}" ]]
}

@test "Verify if VARIABLE ${VARIABLE_2} has value ${VARIABLE_2_VALUE} " {
    run gcloud beta runtime-config configs variables get-value ${VARIABLE_2} \
        --config-name ${CONFIG_NAME}
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "my test text value" ]]
}

@test "Verify if WAITER ${WAITER_NAME} is created " {
    run gcloud beta runtime-config configs waiters list \
        --config-name ${CONFIG_NAME} \
        --format="value(name)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${WAITER_NAME}" ]]
}

@test "Verify if WAITER ${WAITER_NAME} has timeout ${WAITER_TIMEOUT} " {
    run gcloud beta runtime-config configs waiters describe ${WAITER_NAME} \
        --config-name ${CONFIG_NAME} \
        --format="value(timeout)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${WAITER_TIMEOUT}" ]]
}

@test "Verify if WAITER ${WAITER_NAME} success path is ${WAITER_PATH} " {
    run gcloud beta runtime-config configs waiters describe ${WAITER_NAME} \
        --config-name ${CONFIG_NAME} \
        --format="value(success.cardinality.path)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${WAITER_PATH}" ]]
}

@test "Deleting deployment ${DEPLOYMENT_NAME}" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q
    [[ "$status" -eq 0 ]]
}
