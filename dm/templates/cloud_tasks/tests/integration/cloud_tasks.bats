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
    export QUEUE_NAME="test-q-${RAND}"
    export TASK_NAME="test-task-${RAND}"
    export DISPATCHES_PER_SECOND="10.0"
    export CONCURRENT_DISPATCHES="5"
    export MAX_ATTEMPTS="100"
    export MAX_RETRY_DURATION="60s"
    export MAX_BACKOFF="3600s"
    export MIN_BACKOFF="0.100s"
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
    run gcloud beta deployment-manager deployments create "${DEPLOYMENT_NAME}"\
        --config "${CONFIG}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

sleep 2

@test "Verify if queue ${QUEUE_NAME} was created " {
    run gcloud beta tasks queues list --format="value(name)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${QUEUE_NAME}" ]]
}

@test "Verify if queue ${QUEUE_NAME} is RUNNING" {
    run gcloud beta tasks queues describe ${QUEUE_NAME} \
        --format="value(state)"
    [[ "$status" -eq 0 ]]
    [[ "$output" -eq "RUNNING" ]]
}

@test "Verify if maxDispatchesPerSecond is ${DISPATCHES_PER_SECOND}" {
    run gcloud beta tasks queues describe ${QUEUE_NAME} \
        --format="value(rateLimits.maxDispatchesPerSecond)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${DISPATCHES_PER_SECOND}" ]]
}

@test "Verify if maxConcurrentDispatches is ${CONCURRENT_DISPATCHES}" {
    run gcloud beta tasks queues describe ${QUEUE_NAME} \
        --format="value(rateLimits.maxConcurrentDispatches)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${CONCURRENT_DISPATCHES}" ]]
}

@test "Verify if retryConfig maxAttempts is set to ${MAX_ATTEMPTS}" {
    run gcloud beta tasks queues describe ${QUEUE_NAME} \
        --format="value(retryConfig.maxAttempts)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${MAX_ATTEMPTS}" ]]
}

@test "Verify if retryConfig maxBackoff is set to ${MAX_BACKOFF}" {
    run gcloud beta tasks queues describe ${QUEUE_NAME} \
        --format="value(retryConfig.maxBackoff)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${MAX_BACKOFF}" ]]
}

@test "Verify if retryConfig minBackoff is set to ${MIN_BACKOFF}" {
    run gcloud beta tasks queues describe ${QUEUE_NAME} \
        --format="value(retryConfig.minBackoff)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${MIN_BACKOFF}" ]]
}

@test "Verify if the task ${TASK_NAME} was created " {
    run gcloud beta tasks describe ${TASK_NAME} \
        --queue ${QUEUE_NAME}
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "createTime:" ]]
    [[ "$output" =~ "${TASK_NAME}" ]]
}

@test "Deleting deployment ${DEPLOYMENT_NAME}" {
    run gcloud beta deployment-manager deployments delete "${DEPLOYMENT_NAME}"\
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q
    [[ "$status" -eq 0 ]]
}
