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
    # Replace underscores with dashes in the deployment name.
    DEPLOYMENT_NAME=${DEPLOYMENT_NAME//_/-}
    CONFIG=".${DEPLOYMENT_NAME}.yaml"
    # Test specific variables.
    export HTTPS_RES_NAME="https-proxy-${RAND}"
    export URL_MAP_RES_NAME="url-map-${RAND}"
    export HTTPS_QUIC_OVERRIDE="ENABLE"
    export SSL_RES_NAME="ssl-proxy-${RAND}"
    export SSL_NAME="ssl-proxy-name-${RAND}"
    export SSL_DESCRIPTION="ssl-proxy-description-${RAND}"
    export SSL_BS_RES_NAME="ssl-backend-service-${RAND}"
    export PROXY_HEADER="PROXY_V1"
    export SSL_CERT_NAME="ssl-certificate-${RAND}"
    export SSL_POLICY_NAME="ssl-policy-${RAND}"
    export HTTP_RES_NAME="http-proxy-${RAND}"
    export HTTP_NAME="https-proxy-name-${RAND}"
    export HTTP_DESCRIPTION="http-proxy-description-${RAND}"
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
        gcloud compute ssl-policies create "${SSL_POLICY_NAME}" \
            --profile MODERN --min-tls-version 1.2 \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    fi

    # Per-test setup steps.
}

function teardown() {
    # Global teardown; executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        delete_config
        gcloud compute ssl-policies delete "${SSL_POLICY_NAME}" -q \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    fi

    # Per-test teardown steps.
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    run gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config ${CONFIG} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

@test "Verifying HTTP proxy" {
    run gcloud compute target-http-proxies describe "${HTTP_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "description: ${HTTP_DESCRIPTION}" ]]
    [[ "$output" =~ "${URL_MAP_RES_NAME}" ]]
}

@test "Verifying HTTPS proxy" {
    run gcloud compute target-https-proxies describe "${HTTPS_RES_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "quicOverride: ${HTTPS_QUIC_OVERRIDE}" ]]
    [[ "$output" =~ "${URL_MAP_RES_NAME}" ]]
    [[ "$output" =~ "${SSL_CERT_NAME}" ]]
}

@test "Verifying SSL proxy" {
    run gcloud compute target-ssl-proxies describe "${SSL_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "description: ${SSL_DESCRIPTION}" ]]
    [[ "$output" =~ "proxyHeader: ${PROXY_HEADER}" ]]
    [[ "$output" =~ "${SSL_CERT_NAME}" ]]
    [[ "$output" =~ "${SSL_POLICY_NAME}" ]]
    [[ "$output" =~ "${SSL_BS_RES_NAME}" ]]
}

@test "Deleting deployment" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

