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
    export RESOURCE_NAME_PREFIX="test-healthcheck-${RAND}"
    export CHECK_INTERVAL_SEC="5"
    export TIMEOUT_SEC="5"
    export UNHEALTHY_THRESHOLD="2"
    export HEALTHY_THRESHOLD="2"
    export PORT_80="80"
fi
########## HELPER FUNCTIONS ##########

function create_config() {
    echo "Creating ${CONFIG}"
    envsubst < "templates/healthcheck/tests/integration/${TEST_NAME}.yaml" > "${CONFIG}"
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
    # Global teardown; this is executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        delete_config
        rm -f ${RANDOM_FILE}
    fi

    # Per-test teardown steps.
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config "${CONFIG}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
}

@test "HTTP healthcheck was created" {
    RESOURCE_NAME=${RESOURCE_NAME_PREFIX}-legacy-http
    run gcloud compute http-health-checks describe ${RESOURCE_NAME} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "checkIntervalSec: ${CHECK_INTERVAL_SEC}" ]]
    [[ "$output" =~ "timeoutSec: ${TIMEOUT_SEC}" ]]
    [[ "$output" =~ "unhealthyThreshold: ${UNHEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "healthyThreshold: ${HEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "port: ${PORT_80}" ]]
}

@test "HTTPS healthcheck was created" {
    RESOURCE_NAME=${RESOURCE_NAME_PREFIX}-legacy-https
    run gcloud compute https-health-checks describe ${RESOURCE_NAME}\
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "checkIntervalSec: ${CHECK_INTERVAL_SEC}" ]]
    [[ "$output" =~ "timeoutSec: ${TIMEOUT_SEC}" ]]
    [[ "$output" =~ "unhealthyThreshold: ${UNHEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "healthyThreshold: ${HEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "port: 443" ]]
}

@test "TCP healthcheck was created" {
    RESOURCE_NAME=${RESOURCE_NAME_PREFIX}-tcp
    run gcloud compute health-checks describe ${RESOURCE_NAME} \
         --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "checkIntervalSec: ${CHECK_INTERVAL_SEC}" ]]
    [[ "$output" =~ "timeoutSec: ${TIMEOUT_SEC}" ]]
    [[ "$output" =~ "unhealthyThreshold: ${UNHEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "healthyThreshold: ${HEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "port: ${PORT_80}" ]]
    [[ "$output" =~ "type: TCP" ]]
}

@test "SSL healthcheck was created" {
    RESOURCE_NAME=${RESOURCE_NAME_PREFIX}-ssl
    run gcloud compute health-checks describe ${RESOURCE_NAME} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "checkIntervalSec: ${CHECK_INTERVAL_SEC}" ]]
    [[ "$output" =~ "timeoutSec: ${TIMEOUT_SEC}" ]]
    [[ "$output" =~ "unhealthyThreshold: ${UNHEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "healthyThreshold: ${HEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "port: ${PORT_80}" ]]
    [[ "$output" =~ "type: SSL" ]]
}

@test "Request path healthcheck was created" {
    RESOURCE_NAME=${RESOURCE_NAME_PREFIX}-requestpath-https
    run gcloud compute https-health-checks describe ${RESOURCE_NAME}\
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "checkIntervalSec: ${CHECK_INTERVAL_SEC}" ]]
    [[ "$output" =~ "timeoutSec: ${TIMEOUT_SEC}" ]]
    [[ "$output" =~ "unhealthyThreshold: ${UNHEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "healthyThreshold: ${HEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "requestPath: /health.html" ]]
    [[ "$output" =~ "port: 443" ]]
}

@test "TCP w/ request/response data was created" {
    RESOURCE_NAME=${RESOURCE_NAME_PREFIX}-response-tcp
    run gcloud compute health-checks describe ${RESOURCE_NAME}\
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "checkIntervalSec: ${CHECK_INTERVAL_SEC}" ]]
    [[ "$output" =~ "timeoutSec: ${TIMEOUT_SEC}" ]]
    [[ "$output" =~ "unhealthyThreshold: ${UNHEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "healthyThreshold: ${HEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "port: ${PORT_80}" ]]
    [[ "$output" =~ "type: TCP" ]]
    [[ "$output" =~ "request: request-data" ]]
    [[ "$output" =~ "response: response-data" ]]
}

@test "HTTP healthcheck was created" {
    RESOURCE_NAME=${RESOURCE_NAME_PREFIX}-beta-http
    run gcloud beta compute http-health-checks describe ${RESOURCE_NAME}\
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "checkIntervalSec: ${CHECK_INTERVAL_SEC}" ]]
    [[ "$output" =~ "timeoutSec: ${TIMEOUT_SEC}" ]]
    [[ "$output" =~ "unhealthyThreshold: ${UNHEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "healthyThreshold: ${HEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "port: ${PORT_80}" ]]
}

@test "HTTPS healthcheck was created" {
    RESOURCE_NAME=${RESOURCE_NAME_PREFIX}-beta-https
    run gcloud beta compute https-health-checks describe ${RESOURCE_NAME}\
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "checkIntervalSec: ${CHECK_INTERVAL_SEC}" ]]
    [[ "$output" =~ "timeoutSec: ${TIMEOUT_SEC}" ]]
    [[ "$output" =~ "unhealthyThreshold: ${UNHEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "healthyThreshold: ${HEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "port: 443" ]]
}

@test "HTTPS healthcheck was created" {
    RESOURCE_NAME=${RESOURCE_NAME_PREFIX}-beta-http2
    run gcloud beta compute health-checks describe ${RESOURCE_NAME}\
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "checkIntervalSec: ${CHECK_INTERVAL_SEC}" ]]
    [[ "$output" =~ "timeoutSec: ${TIMEOUT_SEC}" ]]
    [[ "$output" =~ "unhealthyThreshold: ${UNHEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "healthyThreshold: ${HEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "port: ${PORT_80}" ]]
}

@test "TCP healthcheck was created" {
    RESOURCE_NAME=${RESOURCE_NAME_PREFIX}-beta-tcp
    run gcloud beta compute health-checks describe ${RESOURCE_NAME} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "checkIntervalSec: ${CHECK_INTERVAL_SEC}" ]]
    [[ "$output" =~ "timeoutSec: ${TIMEOUT_SEC}" ]]
    [[ "$output" =~ "unhealthyThreshold: ${UNHEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "healthyThreshold: ${HEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "port: ${PORT_80}" ]]
    [[ "$output" =~ "type: TCP" ]]
}

@test "SSL healthcheck was created" {
    RESOURCE_NAME=${RESOURCE_NAME_PREFIX}-beta-ssl
    run gcloud beta compute health-checks describe ${RESOURCE_NAME} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "checkIntervalSec: ${CHECK_INTERVAL_SEC}" ]]
    [[ "$output" =~ "timeoutSec: ${TIMEOUT_SEC}" ]]
    [[ "$output" =~ "unhealthyThreshold: ${UNHEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "healthyThreshold: ${HEALTHY_THRESHOLD}" ]]
    [[ "$output" =~ "port: ${PORT_80}" ]]
    [[ "$output" =~ "type: SSL" ]]
}

@test "Deleting deployment" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q
    [[ "$status" -eq 0 ]]
}
