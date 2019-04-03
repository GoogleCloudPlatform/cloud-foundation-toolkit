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
    # Test specific variables
    export CLOUDBUILD_NAME="test-build-${RAND}"
    export BUILD_TIMEOUT="500s"
    export IMAGE_NAME="test-npm-helloworld-${RAND}"
    export IMAGE_TAG="gcr.io/${CLOUD_FOUNDATION_PROJECT_ID}/${IMAGE_NAME}"
    export LOGURL_BASE="https://console.cloud.google.com/gcr/builds/"
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
        # delete the container image from registry
        gcloud container images delete ${IMAGE_TAG}:latest -q
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

# Get the BuildID after DM deployment and store it in a variable for reuse
export ID=$(gcloud builds list --format="value(id)" --filter="(${IMAGE_NAME})")

@test "Verify if build ${CLOUDBUILD_NAME} was created " {
    run gcloud builds list --format="value(id,images)" \
        --filter="(${IMAGE_NAME})"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${IMAGE_NAME}" ]]
}

@test "Verify if build ${CLOUDBUILD_NAME} status is a SUCCESS" {
    run gcloud builds describe $ID --format="value(status)"
    [[ "$status" -eq 0 ]]
    [[ "$output" -eq "SUCCESS" ]]
}

@test "Verify if build timeout is set to ${BUILD_TIMEOUT}" {
    run gcloud builds describe $ID --format="value(timeout)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${BUILD_TIMEOUT}" ]]
}

@test "Verify if cloud-builder in STEP 1 is git" {
    run gcloud builds describe $ID --format="value(steps[0].name)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "git" ]]
}

@test "Verify if first build arg in STEP 1 is clone" {
    run gcloud builds describe $ID --format="value(steps[0].args[0])"
    [[ "$status" -eq 0 ]]
    [[ "$output" -eq "clone" ]]
}

@test "Verify if cloud-builder in STEP 2 is docker" {
    run gcloud builds describe $ID --format="value(steps[1].name)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "docker" ]]
}

@test "Verify if first build arg in STEP 2 is build" {
    run gcloud builds describe $ID --format="value(steps[1].args[0])"
    [[ "$status" -eq 0 ]]
    [[ "$output" -eq "build" ]]
}

@test "Verify if relative dir in STEP 2 is npm/examples/hello_world" {
    run gcloud builds describe $ID --format="value(steps[1].dir)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "cloud-builders/npm/examples/hello_world" ]]
}

@test "Verify if image stored in container repo is ${IMAGE_TAG} " {
    run gcloud builds describe $ID --format="value(images)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${IMAGE_TAG}" ]]
}

@test "Deleting deployment ${DEPLOYMENT_NAME}" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q
    [[ "$status" -eq 0 ]]
}
