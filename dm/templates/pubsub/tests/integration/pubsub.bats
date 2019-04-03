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
    envsubst < "templates/pubsub/tests/integration/${TEST_NAME}.yaml" > "${CONFIG}"
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
    Global teardown; this is executed once per test file.
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

@test "Verifying that test-topic-${RAND} was created in deployment ${DEPLOYMENT_NAME}" {
    run gcloud pubsub topics list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "test-topic-${RAND}" ]]
}

@test "Verifying that test-topic-${RAND}'s IAM policy was set" {
    run gcloud beta pubsub topics get-iam-policy test-topic-${RAND} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "demo@user.com" ]]
}

@test "Verifying that two subscriptions were created in deployment ${DEPLOYMENT_NAME}" {
    run gcloud pubsub subscriptions list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "first-subscription-${RAND}" ]]
    [[ "$output" =~ "second-subscription-${RAND}" ]]
}

@test "Verifying that first-subscription-${RAND}'s topic is test-topic-${RAND}" {
    run gcloud pubsub subscriptions describe first-subscription-${RAND} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "test-topic-${RAND}" ]]
}

@test "Verifying that first-subscription-${RAND}'s IAM policy was set" {
    run gcloud beta pubsub subscriptions get-iam-policy first-subscription-${RAND} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "demo@user.com" ]]
}

@test "Verifying that second-subscription-${RAND}'s topic is test-topic-${RAND}" {
    run gcloud pubsub subscriptions describe second-subscription-${RAND} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "test-topic-${RAND}" ]]
}

@test "Verifying that second-subscription-${RAND}'s ackDeadlineSeconds was set" {
    run gcloud pubsub subscriptions describe second-subscription-${RAND} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "ackDeadlineSeconds: 15" ]]
}

@test "Deleting deployment" {
    gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"

    run gcloud pubsub topics list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ ! "$output" =~ "test-topic-${RAND}" ]]

    run gcloud pubsub subscriptions list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ ! "$output" =~ "first-subscription-${RAND}" ]]
    [[ ! "$output" =~ "second-subscription-${RAND}" ]]
}
