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
    # DM config file name must be 61 chars or less.
    # Must be a match of regex '[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?'
    DEPLOYMENT_NAME=`echo $DEPLOYMENT_NAME | cut -c 1-61`
    CONFIG=".${DEPLOYMENT_NAME}.yaml"
    # Test specific variables.

    export SLACK_CHANNEL_NAME="#slack-channel"
    export SLACK_TOKEN="token-1234567890"
    export SLACK_DISPLAY_NAME="name"
    export SLACK_TYPE="slack"
    export TEST_POLICY_NAME="1 - Availability - Cloud SQL Database - Memory usage (filtered) [MAX]"
    export CONDITION_DISPLAY_NAME="CloudSQL Memory"
    export CONDITION_FILTER="metric.type=\\\"cloudsql.googleapis.com/database/memory/usage\\\" resource.type=\\\"cloudsql_database\\\" resource.label.database_id=\\\"sql_instance_id\\\""
    export CONDITION_COMPARISON="COMPARISON_GT"
    export CONDITION_DURATION="300s"
    export CONDITION_THRESHOLD_VALUE=2750000000
    export CONDITION_TRIGGER_COUNT=1
    export CONDITION_AGGREGATION_ALIGNMENT_PERIOD="60s"
    export CONDITION_AGGREGATION_ALIGNMENT_PER_SERIES="ALIGN_MAX"
    export CONDITION_AGGREGATION_CROSS_SERIES_REDUCER="REDUCE_MEAN"
    export CONDITION_AGGREGATION_GROUP_BY_FIELD="project"
    export TEST_POLICY_DOCUMENTATION_CONTEXT="The janus rule \${condition.display_name} has generated this alert for the \${metric.display_name}."


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
        --config "${CONFIG}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

@test "Check Slack notification channel configuration" {
    run gcloud alpha monitoring channels list \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "displayName: ${SLACK_DISPLAY_NAME}" ]]
    [[ "$output" =~ "channel_name: '${SLACK_CHANNEL_NAME}'" ]]
    [[ "$output" =~ "type: ${SLACK_TYPE}" ]]
    [[ "$output" =~ "name: projects/${CLOUD_FOUNDATION_PROJECT_ID}/notificationChannels/" ]]
}

@test "Check alert policies" {
    run gcloud alpha monitoring policies list \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "combiner: OR" ]]
    [[ "$output" =~ "displayName: ${TEST_POLICY_NAME}" ]]
    [[ "$output" =~ "alignmentPeriod: ${CONDITION_AGGREGATION_ALIGNMENT_PERIOD}" ]]
    [[ "$output" =~ "crossSeriesReducer: ${CONDITION_AGGREGATION_CROSS_SERIES_REDUCER}" ]]
    [[ "$output" =~ "perSeriesAligner: ${CONDITION_AGGREGATION_ALIGNMENT_PER_SERIES}" ]]
    [[ "$output" =~ "comparison: ${CONDITION_COMPARISON}" ]]
    [[ "$output" =~ "duration: ${CONDITION_DURATION}" ]]
    [[ "$output" =~ "displayName: ${CONDITION_DISPLAY_NAME}" ]]
    [[ "$output" =~ "count: ${CONDITION_TRIGGER_COUNT}" ]]
}

@test "Deleting deployment ${DEPLOYMENT_NAME}" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q
    [[ "$status" -eq 0 ]]
}
