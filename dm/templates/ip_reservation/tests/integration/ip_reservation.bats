#!/usr/bin/env bats

source tests/helpers.bash

TEST_NAME=$(basename "${BATS_TEST_FILENAME}" | cut -d '.' -f 1)

## Create a random 10-char string and save it in a file.
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
    envsubst < "templates/ip_reservation/tests/integration/${TEST_NAME}.yaml" > "${CONFIG}"
}

function delete_config() {
    echo "Deleting ${CONFIG}"
    rm -f "${CONFIG}"
}

function setup() {
    if [ ${BATS_TEST_NUMBER} -eq 1 ]; then
        gcloud compute networks create network-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --description "integration test ${RAND}" \
            --subnet-mode custom
        gcloud compute networks subnets create subnet-${RAND} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
        --network=network-${RAND} --region=us-central1 \
        --range=10.100.0.0/23
        create_config
    fi
}

function teardown() {
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        gcloud compute networks subnets delete subnet-${RAND} \
            --region=us-central1 --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q
        gcloud compute networks delete network-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q
        rm -f "${RANDOM_FILE}"
        delete_config
    fi
}

@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config ${CONFIG} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
}

@test "Verifying that global IPs were created as part of deployment ${DEPLOYMENT_NAME}" {
    run gcloud compute addresses describe test-myglobal-"${RAND}" --global \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "name: test-myglobal-${RAND}" ]]
    [[ "$output" =~ "status: RESERVED" ]]
    [[ "$output" =~ "description: my global ip" ]]
}

@test "Verifying that internal IPs were created as part of deployment ${DEPLOYMENT_NAME}" {
    run gcloud compute addresses describe test-myinternal-"${RAND}" \
        --region us-central1 --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "name: test-myinternal-${RAND}" ]]
    [[ "$output" =~ "status: RESERVED" ]]
    [[ "$output" =~ "addressType: INTERNAL" ]]
    [[ "$output" =~ "description: my us-central1 internal ip" ]]
}

@test "Verifying that external static IPs are created as part of deployment ${DEPLOYMENT_NAME}" {
    run gcloud compute addresses describe test-myregionalexternal-"${RAND}" \
        --region us-central1 --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "name: test-myregionalexternal-${RAND}" ]]
    [[ "$output" =~ "status: RESERVED" ]]
    [[ "$output" =~ "description: my us-central1 static external ip" ]]
}

@test "Deleting deployment" {
    gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q

    run gcloud  run gcloud compute addresses list \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ ! "$output" =~ "test-myglobal-${RAND}" ]]
    [[ ! "$output" =~ "test-myinternal-${RAND}" ]]
    [[ ! "$output" =~ "test-myregionalexternal-${RAND}" ]]
}
