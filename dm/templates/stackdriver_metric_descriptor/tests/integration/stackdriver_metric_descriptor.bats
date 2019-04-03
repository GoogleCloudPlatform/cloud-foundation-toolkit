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
    export METRIC_NAME="test-metric-${RAND}"
    export METRIC_TYPE="custom.googleapis.com/agent/log_entry_retry_count"
    export METRIC_KIND="CUMULATIVE"
    export VALUE_TYPE="INT64"
    export UNIT="1"
    export LAUNCH_STAGE="ALPHA"
    export SAMPLE_PERIOD="10s"
    export INGEST_DELAY="1s"
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

@test "Deleting deployment ${DEPLOYMENT_NAME}" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q
    [[ "$status" -eq 0 ]]
}

########## NOTE ###########
#
# From Google Cloud SDK version 221.0.0, beta 2018.07.16, the only way to
# list a metric descriptor is to make an API call to the
# project.metricDescriptors resource type. Hence, no test assertions were
# written.
#
# The following logging commands do not list custom metricDescriptors:
#   `gcloud logging metrics list`
#   `gcloud beta logging metrics list`
#
#
# References:
# https://cloud.google.com/monitoring/api/ref_v3/rest/v3/projects.metricDescriptors
# https://cloud.google.com/monitoring/custom-metrics/creating-metrics
#
###########################
