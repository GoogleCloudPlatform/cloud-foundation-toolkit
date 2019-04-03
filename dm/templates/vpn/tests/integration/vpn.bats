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
fi

########## HELPER FUNCTIONS ##########

function create_config() {
    echo "Creating ${CONFIG}"
    envsubst < "templates/vpn/tests/integration/${TEST_NAME}.yaml" > "${CONFIG}"
}

function delete_config() {
    echo "Deleting ${CONFIG}"
    rm -f "${CONFIG}"
}

function setup() {
    # Global setup; this is executed once per test file.
    if [ ${BATS_TEST_NUMBER} -eq 1 ]; then
        gcloud compute networks create "network-${RAND}" \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --description "integration test ${RAND}" \
            --subnet-mode custom
        gcloud compute networks subnets create "subnet-${RAND}" \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --network "network-${RAND}" \
            --range 10.118.8.0/22 \
            --region us-east1
        gcloud compute routers create "router-${RAND}" \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --network "network-${RAND}" \
            --asn 65001 \
            --region us-east1
        create_config
    fi

    # Per-test setup steps.
}

function teardown() {
    # Global teardown; this is executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        gcloud compute routers delete "router-${RAND}" \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --region us-east1 -q
        gcloud compute networks subnets delete "subnet-${RAND}" \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --region us-east1 -q
        gcloud compute networks delete "network-${RAND}" \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q
        delete_config
        rm -f "${RANDOM_FILE}"
    fi

    # Per-test teardown steps.
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" --config "${CONFIG}"
}

@test "Verifying that resources were created in deployment ${DEPLOYMENT_NAME}" {
    run gcloud compute networks list --filter="name:network-${RAND}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "network-${RAND}" ]]

    run gcloud compute routers list --filter="name:router-${RAND}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "router-${RAND}  us-east1  network-${RAND}" ]]
}

@test "Verifying the the static address was created in deployment ${DEPLOYMENT_NAME}" {

    run gcloud compute addresses list --filter="name:test-vpn-${RAND}-ip" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "test-vpn-${RAND}-ip  us-east1" ]]
}

@test "Verifying that the target VPN gateway was created in deployment ${DEPLOYMENT_NAME}" {

    run gcloud compute target-vpn-gateways list \
        --filter="name:test-vpn-${RAND}-tvpng" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "test-vpn-${RAND}-tvpng  network-${RAND}  us-east1" ]]
}

@test "Verifying that the VPN tunnel was created in deployment ${DEPLOYMENT_NAME}" {

    run gcloud compute vpn-tunnels list --filter="name:test-vpn-${RAND}-vpn" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "test-vpn-${RAND}-vpn  us-east1  test-vpn-${RAND}-tvpng  1.2.3.4" ]]
}

@test "Verifying that the forwarding rules were created in deployment ${DEPLOYMENT_NAME}" {

    run gcloud compute forwarding-rules list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [ "$status" -eq 0 ]
    [[ "$output" =~ "test-vpn-${RAND}-esp-rule       us-east1" ]]
    [[ "$output" =~ "test-vpn-${RAND}-udp-4500-rule  us-east1" ]]
    [[ "$output" =~ "test-vpn-${RAND}-udp-500-rule   us-east1" ]]
}

@test "Deleting deployment" {
    gcloud deployment-manager deployments delete ${DEPLOYMENT_NAME} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q

    run gcloud compute addresses list --filter="name:test-vpn-${RAND}-ip" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ ! "$output" =~ "test-vpn-${RAND}-ip" ]]

    run gcloud compute target-vpn-gateways list \
        --filter="name:test-vpn-${RAND}-tvpng" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ ! "$output" =~ "test-vpn-${RAND}-tvpng" ]]

    run gcloud compute vpn-tunnels list --filter="name:test-vpn-${RAND}-vpn" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ ! "$output" =~ "test-vpn-${RAND}-vpn" ]]

    run gcloud compute forwarding-rules list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ ! "$output" =~ "test-vpn-${RAND}-esp-rule" ]]
    [[ ! "$output" =~ "test-vpn-${RAND}-udp-4500-rule" ]]
    [[ ! "$output" =~ "test-vpn-${RAND}-udp-500-rule" ]]
}
