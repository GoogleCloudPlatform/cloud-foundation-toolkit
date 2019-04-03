#!/usr/bin/env bats

source tests/helpers.bash

# Create a random 10-char string and save it in a file.
RANDOM_FILE="/tmp/${CLOUD_FOUNDATION_ORGANIZATION_ID}-natgatewayha.txt"
TEST_NAME=$(basename "${BATS_TEST_FILENAME}" | cut -d '.' -f 1)
if [[ ! -e "${RANDOM_FILE}" ]]; then
    RAND=$(head /dev/urandom | LC_ALL=C tr -dc a-z0-9 | head -c 10)
    echo ${RAND} > "${RANDOM_FILE}"
fi

# Set variables based on the random string saved in the file.
# envsubst requires all variables used in the example/config to be exported.
if [[ -e "${RANDOM_FILE}" ]]; then
    export RAND=$(cat "${RANDOM_FILE}")
    DEPLOYMENT_NAME="${CLOUD_FOUNDATION_PROJECT_ID}-natgatewayha-${RAND}"
    # Replace underscores in the deployment name with dashes.
    DEPLOYMENT_NAME=${DEPLOYMENT_NAME//_/-}
    CONFIG=".${DEPLOYMENT_NAME}.yaml"
fi

export PROJECT_NUMBER=$(gcloud projects list | grep "${CLOUD_FOUNDATION_PROJECT_ID}" | awk {'print $NF'})

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
        gcloud compute networks create "network-${RAND}" \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --description "integration test ${RAND}" \
            --subnet-mode custom
        gcloud compute networks subnets create "subnet-${RAND}" \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --network "network-${RAND}" \
            --range 10.0.1.0/24 \
            --region us-east1
        create_config
    fi
    # Per-test setup steps.
}

function teardown() {
    # Global teardown; executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
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
    run gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config "${CONFIG}" --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

@test "Verifying that resources were created in deployment ${DEPLOYMENT_NAME}" {
    run gcloud compute instances list --filter="name:test-nat-gateway-${RAND}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "test-nat-gateway-${RAND}-gateway-us-east1-b" ]]
    [[ "$output" =~ "test-nat-gateway-${RAND}-gateway-us-east1-c" ]]
    [[ "$output" =~ "test-nat-gateway-${RAND}-gateway-us-east1-d" ]]
}

@test "Verifying that external IP was created in deployment ${DEPLOYMENT_NAME}" {
    run gcloud compute addresses list \
        --filter="name:test-nat-gateway-${RAND}-ip-external" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "test-nat-gateway-${RAND}-ip-external-us-east1-b" ]]
    [[ "$output" =~ "test-nat-gateway-${RAND}-ip-external-us-east1-c" ]]
    [[ "$output" =~ "test-nat-gateway-${RAND}-ip-external-us-east1-d" ]]
}

@test "Verifying that internal IP was created in deployment ${DEPLOYMENT_NAME}" {
    run gcloud compute addresses list \
        --filter="name:test-nat-gateway-${RAND}-ip-internal" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "test-nat-gateway-${RAND}-ip-internal-us-east1-b" ]]
    [[ "$output" =~ "test-nat-gateway-${RAND}-ip-internal-us-east1-c" ]]
    [[ "$output" =~ "test-nat-gateway-${RAND}-ip-internal-us-east1-d" ]]
}

@test "Verifying that routes were created in deployment ${DEPLOYMENT_NAME}" {
    run gcloud compute routes list \
        --filter="name:test-nat-gateway-${RAND}-route" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "test-nat-gateway-${RAND}-route-us-east1-b" ]]
    [[ "$output" =~ "test-nat-gateway-${RAND}-route-us-east1-c" ]]
    [[ "$output" =~ "test-nat-gateway-${RAND}-route-us-east1-d" ]]
}

@test "Verifying that firewall rule was created in deployment ${DEPLOYMENT_NAME}" {
    run gcloud compute firewall-rules list \
        --filter="name:test-nat-gateway-${RAND}-healthcheck-firewall" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "test-nat-gateway-${RAND}-healthcheck-firewall" ]]
}

@test "Verifying NAT functionality created in deployment ${DEPLOYMENT_NAME}" {
    # SSH into the instance with external IP and SSH into the instance without
    # an external IP that is using the NAT gateway and successfully execute
    # wget on a site.
    run gcloud compute ssh "test-inst-has-ext-ip-${RAND}" --zone "us-east1-b" \
        --ssh-flag="-q" \
        --command "gcloud compute ssh test-inst-nat-no-ext-ip-${RAND} \
            --internal-ip --command 'wget google.com' --zone 'us-east1-b' \
            --quiet" \
        --quiet
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "HTTP request sent, awaiting response... 200 OK" ]]

    # SSH into the instance with external IP and SSH into the instance without
    # an external IP that is not using the NAT gateway. The wget command will
    # fail.
    run gcloud compute ssh "test-inst-has-ext-ip-${RAND}" --zone "us-east1-b" \
        --ssh-flag="-q" \
        --command "gcloud compute ssh test-inst-no-ext-ip-${RAND} --internal-ip \
            --command 'wget google.com --timeout=5' --zone 'us-east1-b' \
            --quiet" \
        --quiet
    [[ "$output" =~ "failed: Network is unreachable" ]]
}

@test "Deleting deployment" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]

    run gcloud compute instances list \
        --filter="name:test-nat-gateway-${RAND}-gw-1-us-east1-b" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ ! "$output" =~ "test-nat-gateway-${RAND}-gateway-us-east1-b" ]]
    [[ ! "$output" =~ "test-nat-gateway-${RAND}-gateway-us-east1-c" ]]
    [[ ! "$output" =~ "test-nat-gateway-${RAND}-gateway-us-east1-d" ]]

    run gcloud compute addresses list \
        --filter="name:test-nat-gateway-${RAND}-ip" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ ! "$output" =~ "test-nat-gateway-${RAND}-ip-external-us-east1-b" ]]
    [[ ! "$output" =~ "test-nat-gateway-${RAND}-ip-external-us-east1-c" ]]
    [[ ! "$output" =~ "test-nat-gateway-${RAND}-ip-external-us-east1-d" ]]
    [[ ! "$output" =~ "test-nat-gateway-${RAND}-ip-internal-us-east1-b" ]]
    [[ ! "$output" =~ "test-nat-gateway-${RAND}-ip-internal-us-east1-c" ]]
    [[ ! "$output" =~ "test-nat-gateway-${RAND}-ip-internal-us-east1-d" ]]

    run gcloud compute routes list \
        --filter="name:test-nat-gateway-${RAND}-route" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ ! "$output" =~ "test-nat-gateway-${RAND}-route-us-east1-b" ]]
    [[ ! "$output" =~ "test-nat-gateway-${RAND}-route-us-east1-c" ]]
    [[ ! "$output" =~ "test-nat-gateway-${RAND}-route-us-east1-d" ]]

    run gcloud compute firewall-rules list \
        --filter="name:test-nat-gateway-${RAND}-healthcheck-firewall" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ ! "$output" =~ "test-nat-gateway-${RAND}-healthcheck-firewall" ]]
}
