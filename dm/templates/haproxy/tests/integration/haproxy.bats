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
fi

########## HELPER FUNCTIONS ##########

function create_config() {
    echo "Creating ${CONFIG}"
    envsubst < "templates/haproxy/tests/integration/${TEST_NAME}.yaml" > "${CONFIG}"
}

function delete_config() {
    echo "Deleting ${CONFIG}"
    rm -f "${CONFIG}"
}

function setup() {
    # Global setup; executed once per test file.
    if [ ${BATS_TEST_NUMBER} -eq 1 ]; then
        # Set up instance groups to be load-balanced via HAProxy.
        gcloud compute instance-templates create template-${RAND}-1 \
            --no-service-account --no-scopes --machine-type=f1-micro \
            --image-project=debian-cloud --image-family=debian-9 \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}"

        gcloud compute instance-groups managed create group-${RAND}-1 \
            --zone us-central1-a --template template-${RAND}-1 --size 1 \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}"

        gcloud compute instance-groups managed create group-${RAND}-2 \
            --zone us-central1-c --template template-${RAND}-1 --size 1 \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}"

        create_config
    fi

  # Per-test setup steps.
}

function teardown() {
    # Global teardown; executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        gcloud compute instance-groups managed delete group-${RAND}-1 \
            --zone us-central1-a --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q

        gcloud compute instance-groups managed delete group-${RAND}-2 \
            --zone us-central1-c --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q

        gcloud compute instance-templates delete template-${RAND}-1 \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q

        rm -f "${RANDOM_FILE}"
        delete_config
    fi

    # Per-test teardown steps.
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    run gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config "${CONFIG}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

@test "Verifying that the HAProxy instance was created in deployment ${DEPLOYMENT_NAME}" {
    run gcloud compute instances list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "ilb-proxy-${RAND}" ]]
}

@test "Verifying that haproxy.cfg was populated with instances and had all properties set" {
     # Wait for the HAProxy instance to be configured.
     until gcloud compute instances get-serial-port-output "ilb-proxy-${RAND}" \
            --zone us-central1-a \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" | grep /etc/haproxy/haproxy.cfg; do

            sleep 10;
     done

    # Verify VM serial output
    run gcloud compute ssh "ilb-proxy-${RAND}" --zone us-central1-a \
        --command "sudo tail -n 15 /etc/haproxy/haproxy.cfg" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "group-${RAND}-1" ]]   # has instances from group 1
    [[ "$output" =~ "group-${RAND}-2" ]]   # has instances from group 2
    [[ "$output" =~ "mode tcp" ]]          # the mode was set
    [[ "$output" =~ "balance leastconn" ]] # load labalcing algorithm is set
    [[ "$output" =~ ":9999" ]]             # Load balancer's port
    [[ "$output" =~ ":8888" ]]             # Instance group's port
}

@test "Verifying that update interval was set" {
    run gcloud compute ssh "ilb-proxy-${RAND}" --zone us-central1-a \
        --command "sudo crontab -l" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"

    [[ "$status" -eq 0 ]]
    [[ "$output" = "*/15 * * * * /sbin/haproxy-conf-updater" ]]
}

@test "Deleting deployment" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}
