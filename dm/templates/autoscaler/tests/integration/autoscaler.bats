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
    export CPU_UTILIZATION_1="0.7"
    export CUSTOM_METRIC="compute.googleapis.com/instance/disk/read_ops_count"
    export CUSTOM_METRIC_TARGET="1000"
    export CUSTOM_METRIC_TYPE="DELTA_PER_SECOND"
    export NUM_REPLICAS="2"
    export CPU_UTILIZATION_2="0.6"
    export COOL_DOWN_PERIOD="70"
    export ZONE="us-central1-c"
    export DESCRIPTION="descr"
fi

########## HELPER FUNCTIONS ##########

function create_config() {
    envsubst < "templates/autoscaler/tests/integration/${TEST_NAME}.yaml" > "${CONFIG}"
}

function delete_config() {
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
        rm -f "${RANDOM_FILE}"
        delete_config
    fi

    # Per-test teardown steps.
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    run gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config "${CONFIG}" --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

@test "Verifying that zonal autoscaler properties were set" {
    run gcloud compute instance-groups managed describe "zonal-igm-${RAND}" \
        --format "yaml(autoscaler)" --zone "${ZONE}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "zonal-autoscaler-${RAND}" ]]
    [[ "$output" =~ "coolDownPeriodSec: ${COOL_DOWN_PERIOD}" ]]
    [[ "$output" =~ "utilizationTarget: ${CPU_UTILIZATION_2}" ]]
    [[ "$output" =~ "description: ${DESCRIPTION}" ]]
}

@test "Verifying that regional autoscaler properties were set" {
    run gcloud compute instance-groups managed describe "regional-igm-${RAND}" \
        --format "yaml(autoscaler)" --region "${REGION}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "utilizationTarget: ${CPU_UTILIZATION_1}" ]]
    [[ "$output" =~ "utilizationTarget: ${CUSTOM_METRIC_TARGET}.0" ]]
    [[ "$output" =~ "maxNumReplicas: ${NUM_REPLICAS}" ]]
    [[ "$output" =~ "minNumReplicas: 1" ]] # default
    [[ "$output" =~ "metric: ${CUSTOM_METRIC}" ]]
}

@test "Deleting deployment" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}
