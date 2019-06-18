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
    envsubst < "templates/route/tests/integration/${TEST_NAME}.yaml" > "${CONFIG}"
}

function delete_config() {
    echo "Deleting ${CONFIG}"
    rm -f "${CONFIG}"
}

function setup() {
    # Global setup; this is executed once per test file.
    if [ ${BATS_TEST_NUMBER} -eq 1 ]; then
        gcloud compute networks create network-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --description "integration test ${RAND}" \
            --subnet-mode custom
        gcloud compute networks subnets create subnet-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --network network-${RAND} \
            --range 10.118.8.0/22 \
            --region us-east1
        gcloud compute routers create router-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --network network-${RAND} \
            --asn 65001 \
            --region us-east1
        gcloud compute target-vpn-gateways create gateway-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --network network-${RAND} \
            --region us-east1
        gcloud compute addresses create staticip-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --region us-east1
        gcloud compute forwarding-rules create esprule-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --target-vpn-gateway gateway-${RAND} \
            --region us-east1 \
            --ip-protocol "ESP" \
            --address staticip-${RAND}
        gcloud compute forwarding-rules create udp4500rule-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --target-vpn-gateway gateway-${RAND} \
            --region us-east1 \
            --ip-protocol "UDP" \
            --address staticip-${RAND} \
            --ports 4500
        gcloud compute forwarding-rules create udp500rule-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --target-vpn-gateway gateway-${RAND} \
            --region us-east1 \
            --ip-protocol "UDP" \
            --address staticip-${RAND} \
            --ports 500
        gcloud compute vpn-tunnels create vpntunnel-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --peer-address 1.2.3.4 \
            --shared-secret 'superSecretPassw0rd' \
            --target-vpn-gateway gateway-${RAND} \
            --router router-${RAND} \
            --region us-east1
        create_config
    fi

    # Per-test setup steps.
}

function teardown() {
    # Global teardown; this is executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        gcloud compute vpn-tunnels delete vpntunnel-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --region us-east1 -q
        gcloud compute forwarding-rules delete udp500rule-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --region us-east1 -q
        gcloud compute forwarding-rules delete udp4500rule-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --region us-east1 -q
        gcloud compute forwarding-rules delete esprule-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --region us-east1 -q
        gcloud compute addresses delete staticip-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --region us-east1 -q
        gcloud compute target-vpn-gateways delete gateway-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --region us-east1 -q
        gcloud compute routers delete router-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --region us-east1 -q
        gcloud compute networks subnets delete subnet-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --region us-east1 -q
        gcloud compute networks delete network-${RAND} \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q
        delete_config
        rm -f ${RANDOM_FILE}
    fi

    # Per-test teardown steps.
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
        --config "${CONFIG}"
}

@test "Verifying that resources were created in deployment ${DEPLOYMENT_NAME}" {
    run gcloud compute routes list --filter="name:gateway-route-${RAND} AND priority:1002" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [ "$status" -eq 0 ]
    [[ "${lines[1]}" =~ "gateway-route-${RAND}" ]]

    run gcloud compute routes list --filter="name:instance-route-${RAND} AND priority:1001" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [ "$status" -eq 0 ]
    [[ "${lines[1]}" =~ "instance-route-${RAND}" ]]

    run gcloud compute routes list --filter="(name:ip-route-${RAND} AND priority:20000)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [ "$status" -eq 0 ]
    [[ "${lines[1]}" =~ "ip-route-${RAND}" ]]

    run gcloud compute routes list --filter="(name:vpn-tunnel-route-${RAND} AND priority:500)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [ "$status" -eq 0 ]
    [[ "${lines[1]}" =~ "vpn-tunnel-route-${RAND}" ]]
}

@test "Deleting deployment" {
    gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q

    run gcloud compute routes list --filter="name:gateway-route-${RAND}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ ! "$output" =~ "gateway-route-${RAND}" ]]

    run gcloud compute routes list --filter="name:instance-route-${RAND}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ ! "$output" =~ "instance-route-${RAND}" ]]

    run gcloud compute routes list --filter="name:ip-route-${RAND}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ ! "$output" =~ "ip-route-${RAND}" ]]

    run gcloud compute routes list --filter="name:vpn-runnel-route-${RAND}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ ! "$output" =~ "vpn-tunnel-route-${RAND}" ]]
}
