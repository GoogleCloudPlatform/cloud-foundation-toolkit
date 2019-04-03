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
    export SSL_RES_NAME="test-ssl-elb-${RAND}"
    export SSL_TARGET_NAME="${SSL_RES_NAME}-target"
    export SSL_PORT_RANGE="443"
    export SSL_PORT_NAME="https"
    export SSL_PROXY_HEADER="PROXY_V1"
    export SSL_BACKEND_NAME="test-ssl-backend-service-${RAND}"
    export SSL_HEALTHCHECK_NAME="test-ssl-healthcheck-${RAND}"
    export HTTPS_IGM_NAME="test-zonal-igm-https-${RAND}"
    export HTTPS_RES_NAME="test-https-elb-${RAND}"
    export HTTPS_CERT_NAME="${HTTPS_RES_NAME}-target-ssl-cert"
    export HTTPS_URL_MAP_NAME="${HTTPS_RES_NAME}-url-map"
    export HTTPS_FIRST_BACKEND_NAME="first-bs-${RAND}"
    export HTTPS_HEALTHCHECK_NAME="test-healthcheck-https-${RAND}"
    export HTTPS_PORT_RANGE="443"
    export HTTPS_PORT_NAME="https"
    export HTTPS_TARGET_NAME="${HTTPS_RES_NAME}-target"
    export QUIC_OVERRIDE="ENABLE"
    export HTTP_RES_NAME="http-elb-${RAND}"
    export HTTP_NAME="http-elb-name-${RAND}"
    export HTTP_URL_MAP_NAME="${HTTP_NAME}-url-map"
    export HTTP_TARGET_NAME="${HTTP_NAME}-target"
    export HTTP_DESCRIPTION="http-elb-description"
    export HTTP_PORT_RANGE="80"
    export HTTP_FIRST_BACKEND_NAME="first-http-bs-${RAND}"
    export HTTP_FIRST_BACKEND_DESC="backend-service-description"
    export HTTP_SECOND_BACKEND_NAME="second-http-bs-${RAND}"
    export HTTP_PORT_NAME="http"
    export HTTP_ENABLE_CDN="true"
    export HTTP_HEALTHCHECK_NAME="test-healthcheck-http-${RAND}"
    export HTTP_IGM_NAME="zonal-igm-http-${RAND}"
    export TIMEOUT_SEC="70"
    export SESSION_AFFINITY="GENERATED_COOKIE"
    export SESSION_AFFINITY_TTL="1000"
    export DRAINING_TIMEOUT="100"
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
    fi

    # Per-test teardown steps.
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    run gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config ${CONFIG} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

@test "Verifying HTTP ELB forwarding rule" {
    run gcloud compute forwarding-rules describe "${HTTP_NAME}" \
        --global \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "IPProtocol: TCP" ]]
    [[ "$output" =~ "loadBalancingScheme: EXTERNAL" ]]
    [[ "$output" =~ "description: ${HTTP_DESCRIPTION}" ]]
    [[ "$output" =~ "portRange: ${HTTP_PORT_RANGE}-${HTTP_PORT_RANGE}" ]]
    [[ "$output" =~ "targetHttpProxies/${HTTP_TARGET_NAME}" ]]
}

@test "Verifying HTTP ELB URL Map references for two backend services" {
    run gcloud compute url-maps describe "${HTTP_URL_MAP_NAME}" \
        --format="value(defaultService)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${HTTP_FIRST_BACKEND_NAME}" ]]

    run gcloud compute url-maps describe "${HTTP_URL_MAP_NAME}" \
        --format="value(pathMatchers[0].defaultService)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${HTTP_SECOND_BACKEND_NAME}" ]]
}

@test "Verifying HTTP ELB first backend service" {
    run gcloud compute backend-services describe "${HTTP_FIRST_BACKEND_NAME}" \
        --global \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "sessionAffinity: ${SESSION_AFFINITY}" ]]
    [[ "$output" =~ "affinityCookieTtlSec: ${SESSION_AFFINITY_TTL}" ]]
    [[ "$output" =~ "${HTTP_IGM_NAME}" ]]
    [[ "$output" =~ "drainingTimeoutSec: ${DRAINING_TIMEOUT}" ]]
    [[ "$output" =~ "description: ${HTTP_FIRST_BACKEND_DESC}" ]]
    [[ "$output" =~ "enableCDN: ${HTTP_ENABLE_CDN}" ]]
    [[ "$output" =~ "${HTTP_HEALTHCHECK_NAME}" ]]
    [[ "$output" =~ "loadBalancingScheme: EXTERNAL" ]]
    [[ "$output" =~ "portName: ${HTTP_PORT_NAME}" ]]
    [[ "$output" =~ "protocol: HTTP" ]]
    [[ "$output" =~ "timeoutSec: ${TIMEOUT_SEC}" ]]
}

@test "Verifying HTTP ELB second backend service" {
    run gcloud compute backend-services describe \
        "${HTTP_SECOND_BACKEND_NAME}" --global \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

@test "Verifying HTTPS ELB forwarding rule" {
    run gcloud compute forwarding-rules describe "${HTTPS_RES_NAME}" \
        --global \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "loadBalancingScheme: EXTERNAL" ]]
    [[ "$output" =~ "portRange: ${HTTPS_PORT_RANGE}-${HTTPS_PORT_RANGE}" ]]
    [[ "$output" =~ "targetHttpsProxies/${HTTPS_TARGET_NAME}" ]]
}

@test "Verifying HTTPS ELB proxy settings" {
    run gcloud compute target-https-proxies describe "${HTTPS_TARGET_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "quicOverride: ${QUIC_OVERRIDE}" ]]
    [[ "$output" =~ "sslCertificates/${HTTPS_CERT_NAME}" ]]
    [[ "$output" =~ "urlMaps/${HTTPS_URL_MAP_NAME}" ]]
}

@test "Verifying SSL ELB forwarding rule" {
    run gcloud compute forwarding-rules describe "${SSL_RES_NAME}" \
        --global \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "portRange: ${SSL_PORT_RANGE}-${SSL_PORT_RANGE}" ]]
    [[ "$output" =~ "targetSslProxies/${SSL_TARGET_NAME}" ]]
}

@test "Verifying SSL ELB proxy settings" {
    run gcloud compute target-ssl-proxies describe "${SSL_TARGET_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "proxyHeader: ${SSL_PROXY_HEADER}" ]]
    [[ "$output" =~ "sslCertificates/${HTTPS_CERT_NAME}" ]]
}

@test "Deleting deployment" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

