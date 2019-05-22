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
    export IMAGE="projects/ubuntu-os-cloud/global/images/family/ubuntu-1804-lts"
fi

########## HELPER FUNCTIONS ##########

function create_config() {
    envsubst < "templates/instance_template/tests/integration/${TEST_NAME}.yaml" > "${CONFIG}"
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

@test "Verifying instance template disk properties" {
    run gcloud compute instance-templates describe it-${RAND} \
        --format "yaml(properties.disks[0].initializeParams)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "diskType: pd-ssd" ]]
    [[ "$output" =~ "sourceImage: ${IMAGE}" ]]
    [[ "$output" =~ "diskSizeGb: '50'" ]]
}

@test "Verifying instance spec properties" {
    run gcloud compute instance-templates describe it-${RAND} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "machineType: f1-micro" ]]
    [[ "$output" =~ "description: Instance description" ]]
    [[ "$output" =~ "canIpForward: true" ]]
}

@test "Verifying instance template properties" {
    run gcloud compute instance-templates describe it-${RAND} \
        --format "value(name, description, properties.labels)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "Template description" ]]
    [[ "$output" =~ "it-${RAND}" ]]
    [[ "$output" =~ "name=wrench" ]]
}

@test "Verifying instance template network tags" {
    run gcloud compute instance-templates describe it-${RAND} \
        --format "yaml(properties.tags)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "ftp" ]]
    [[ "$output" =~ "https" ]]
}

@test "Verifying instance template metadata" {
    run gcloud compute instance-templates describe it-${RAND} \
        --format "yaml(properties.metadata)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "key: createdBy" ]]
    [[ "$output" =~ "value: unitTest" ]]
}

@test "Verifying instance template network properties" {
    NET="https://www.googleapis.com/compute/v1/projects/${CLOUD_FOUNDATION_PROJECT_ID}/global/networks/test-network-${RAND}"
    run gcloud compute instance-templates describe it-${RAND} \
        --format "yaml(properties.networkInterfaces[0])" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "name: External NAT" ]]
    [[ "$output" =~ "type: ONE_TO_ONE_NAT" ]]
    [[ "$output" =~ "network: ${NET}" ]]
}

@test "Deleting deployment" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}
