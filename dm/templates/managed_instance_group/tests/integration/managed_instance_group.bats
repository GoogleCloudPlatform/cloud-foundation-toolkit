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
    export ZONAL_MIG_NAME="zonal-mig-${RAND}"
    export ZONAL_MIG_RES_NAME="mig-${RAND}"
    export ZONE="us-central1-c"
    export REGION="us-central1"
    export AUTOSCALER_NAME="autoscaler-${RAND}"
    export COOL_DOWN_PERIOD="70"
    export MIN_SIZE="1"
    export TARGET_SIZE="2"
    export UTILIZATION_TARGET="0.7"
    export PORT_NAME="http"
    export PORT="80"
    export BASE_INSTANCE_NAME="bin-${RAND}"
    export INSTANCE_TEMPLATE_NAME="it-${RAND}"
    export IT_NETWORK="default"
    export IT_BASE_IMAGE="projects/ubuntu-os-cloud/global/images/family/ubuntu-1804-lts"
    export REGIONAL_MIG_NAME="regional-mig-${RAND}"
    export HEALTH_CHECK_NAME="test-healthcheck-http-${RAND}"
    export SECOND_HEALTH_CHECK_NAME="second-test-healthcheck-http-${RAND}"
    export INITIAL_DELAY_SEC="450"
fi

########## HELPER FUNCTIONS ##########

function create_config() {
    envsubst < ${BATS_TEST_DIRNAME}/${TEST_NAME}.yaml > "${CONFIG}"
}

function delete_config() {
    rm -f "${CONFIG}"
}

function setup() {
    # Global setup; executed once per test file.
    if [ ${BATS_TEST_NUMBER} -eq 1 ]; then
        create_config
        # Needed for testing resource creation with preexisting (not referenced)
        # health check
        gcloud compute http-health-checks create "${SECOND_HEALTH_CHECK_NAME}" \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    fi

    # Per-test setup steps.
}

function teardown() {
    # Global teardown; executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        rm -f "${RANDOM_FILE}"
        gcloud compute http-health-checks delete "${SECOND_HEALTH_CHECK_NAME}" \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q
        delete_config
    fi

    # Per-test teardown steps.
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    run gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config "${CONFIG}" --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

@test "Verifying that a zonal intance group was created" {
    run gcloud compute instance-groups managed list \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${ZONAL_MIG_NAME}" ]]
    [[ "$output" =~ "${ZONE}" ]]
}

@test "Verifying regional instance group properties" {
    run gcloud compute instance-groups managed list \
        --filter "name=(${REGIONAL_MIG_NAME})" \
        --format "yaml(region)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${REGION}" ]]
}

@test "Verifying zonal instance group properties" {
    run gcloud compute instance-groups managed describe "${ZONAL_MIG_NAME}" \
        --zone "${ZONE}" --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "baseInstanceName: ${BASE_INSTANCE_NAME}" ]]
    [[ "$output" =~ "instanceGroups/${ZONAL_MIG_NAME}" ]]
    [[ "$output" =~ "instanceTemplates/${INSTANCE_TEMPLATE_NAME}" ]]
    [[ "$output" =~ "name: ${PORT_NAME}" ]]
    [[ "$output" =~ "port: ${PORT}" ]]
}

@test "Verifying autoscaler properties" {
    run gcloud compute instance-groups managed describe "${ZONAL_MIG_NAME}" \
        --zone "${ZONE}" --format="yaml(autoscaler)"\
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "cpuUtilization:" ]]
    [[ "$output" =~ "utilizationTarget: ${UTILIZATION_TARGET}" ]]
    [[ "$output" =~ "coolDownPeriodSec: ${COOL_DOWN_PERIOD}" ]]
    [[ "$output" =~ "maxNumReplicas: ${TARGET_SIZE}" ]]
    [[ "$output" =~ "minNumReplicas: ${MIN_SIZE}" ]]
    [[ "$output" =~ "name: ${AUTOSCALER_NAME}" ]]
}

@test "Verifying instance template properties" {
    run gcloud compute instance-templates describe "${INSTANCE_TEMPLATE_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${IT_BASE_IMAGE}" ]]
    [[ "$output" =~ "networks/${IT_NETWORK}" ]]
}

@test "Verifying regional instance group health check properties" {
    run gcloud beta compute instance-groups managed describe \
        "${REGIONAL_MIG_NAME}" --region "${REGION}" \
        --format "yaml(autoHealingPolicies)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${HEALTH_CHECK_NAME}" ]]
    [[ "$output" =~ "initialDelaySec: ${INITIAL_DELAY_SEC}" ]]
}

@test "Verifying zonal instance group health check properties" {
    run gcloud beta compute instance-groups managed describe \
        "${ZONAL_MIG_NAME}" --zone "${ZONE}" \
        --format "yaml(autoHealingPolicies)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${HEALTH_CHECK_NAME}" ]]
    [[ "$output" =~ "initialDelaySec: ${INITIAL_DELAY_SEC}" ]]
}

@test "Deleting deployment" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

