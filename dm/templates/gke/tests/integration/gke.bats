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
    export CLUSTER_NAME="testcluster-${RAND}"
    export REGION="us-east1"
    export NETWORK_NAME="test-k8nw-${RAND}"
    export SUBNET_NAME="test-k8subnet-${RAND}"
    export MACHINE_TYPE="n1-standard-1"
    export NODE_COUNT="1"
    export LOCALSSD_COUNT="1"
    export CLUSTER_VERSION="latest"
    export LOGGING_SERVICE="logging.googleapis.com"
    export MONITORING_SERVICE="monitoring.googleapis.com"
    export MASTERIPV4_CIDRBLOCK="172.16.0.0/28"
    export CLUSTERIPV4_CIDR="10.0.0.0/11"
    export SERVICESIPV4_CIDRBLOCK="10.96.0.0/18"
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
        # Create a VPC network and a subnet for deploying the cluster.
        gcloud compute networks create "${NETWORK_NAME}" \
            --subnet-mode custom
        
        gcloud compute networks subnets create "${SUBNET_NAME}" \
            --region ${REGION} --network "${NETWORK_NAME}" \
            --range 10.200.0.0/24
    fi

    # Per-test setup steps.
}

function teardown() {
    # Global teardown; executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        delete_config
        rm -f "${RANDOM_FILE}"
        # Delete the VPC subnets and network after the tests are completed.
        gcloud compute networks subnets delete "${SUBNET_NAME}" \
            --region ${REGION} -q
        
        gcloud compute networks delete "${NETWORK_NAME}" -q
    fi

    # Per-test teardown steps.
}

########## TESTS ##########

@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config "${CONFIG}" --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

@test "Verify if cluster: ${CLUSTER_NAME} was created " {
    run gcloud container clusters describe "${CLUSTER_NAME}" \
        --region ${REGION} --format="value(name)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${CLUSTER_NAME}" ]]
}

@test "Cluster ${CLUSTER_NAME} is deployed to network ${NETWORK_NAME}" {
    run gcloud container clusters describe "${CLUSTER_NAME}" \
        --region ${REGION} --format="value(network)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${NETWORK_NAME}" ]]
}

@test "Network ${NETWORK_NAME} has subnet ${SUBNET_NAME}" {
    run gcloud container clusters describe "${CLUSTER_NAME}" \
        --region ${REGION} \
        --format="value(networkConfig.subnetwork.scope(subnetworks))"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${SUBNET_NAME}" ]]
}

@test "NodeCount on ${CLUSTER_NAME} is ${NODE_COUNT}" {
    run gcloud container clusters describe "${CLUSTER_NAME}" \
        --region ${REGION} --format="value(currentNodeCount)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "2" ]]
}

@test "Cluster ${CLUSTER_NAME} machineType is ${MACHINE_TYPE}" {
    run gcloud container clusters describe "${CLUSTER_NAME}" \
        --region ${REGION} --format="value(nodeConfig.machineType)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${MACHINE_TYPE}" ]]
}

@test "Logging service on ${CLUSTER_NAME} is ${LOGGING_SERVICE}" {
    run gcloud container clusters describe "${CLUSTER_NAME}" \
        --region ${REGION} --format="value(loggingService)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${LOGGING_SERVICE}" ]]
}

@test "Monitoring service on ${CLUSTER_NAME} is ${MONITORING_SERVICE}" {
    run gcloud container clusters describe "${CLUSTER_NAME}" \
        --region ${REGION} --format="value(monitoringService)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${MONITORING_SERVICE}" ]]
}

@test "Deleting deployment ${DEPLOYMENT_NAME}" {
    gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q
    [[ "$status" -eq 0 ]]

    run gcloud container clusters describe "${CLUSTER_NAME}" \
        --region ${REGION} --format="value(name)"
    [[ "$status" -ne 0 ]]
}
