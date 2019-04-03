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
    export BASTION1_RES_NAME="test-bastion-w-sudo-${RAND}"
    export DEFAULT_MACHINE_TYPE="f1-micro"
    export BASTION1_MACHINE_TYPE="n1-standard-1"
    export BASTION2_RES_NAME="test-bastion-wo-sudo-${RAND}"
    export BASTION2_NAME="test-bastion-wo-sudo-name-${RAND}"
    export ZONE="us-central1-c"
    export BASTION1_DISABLE_SUDO="false"
    export BASTION2_DISABLE_SUDO="true"
    export BASTION2_DISK_SIZE="20"
    export NETWORK_NAME="test-network-${RAND}"
    export PROVISION_COMPLETED_MARKER="provision-completed-marker"
    export BASTION2_STARTUP="echo '${PROVISION_COMPLETED_MARKER}'"
    export BASTION2_EXTRA_TAG="extra"
    export BASTION2_TAG="bastion-host"
    export SSH_TO_BASTION_RULE_NAME="allow-ssh-to-bastion-${RAND}"
    export SSH_FROM_BASTION_DEFAULT_RULE_NAME="allow-ssh-from-bastion"
    export SSH_TO_BASTION_PRIORITY="1001"
    export SSH_FROM_BASTION_PRIORITY="1002"
    export SSH_TO_BASTION_SOURCE_RANGE="0.0.0.0/0"
    export SSH_TO_BASTION_SOURCE_TAG="bastion-trustee"
    export SSH_FROM_BASTION_SOURCE_TAG="bastion-target"
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

@test "Verifying the first Bastion host" {
    run gcloud compute instances describe ${BASTION1_RES_NAME} \
        --zone ${ZONE} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "machineTypes/${BASTION1_MACHINE_TYPE}" ]]
    [[ "$output" =~ "zones/$(ZONE)" ]]
    [[ "$output" =~ "${NETWORK_NAME}" ]]
}

@test "Verifying the first Bastion's sudo is ON" {
    # Wait until VM provisioning finishes
    until gcloud compute instances get-serial-port-output \
        ${BASTION1_RES_NAME} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
        --zone ${ZONE} | grep ${PROVISION_COMPLETED_MARKER}; do

        sleep 10;
    done

    run gcloud compute ssh ${BASTION1_RES_NAME} --command "sudo whoami" \
        --zone ${ZONE} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "root" ]]
}

@test "Verifying the second Bastion host" {
    run gcloud compute instances describe ${BASTION2_NAME} \
        --zone ${ZONE} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "machineTypes/${DEFAULT_MACHINE_TYPE}" ]]
    [[ "$output" =~ "zones/${ZONE}" ]]
    [[ "$output" =~ "sudo EDITOR=tee visudo" ]] # disable sudo startup script
    [[ "$output" =~ "${BASTION2_STARTUP}" ]] # user startup script
    [[ "$output" =~ "${NETWORK_NAME}" ]]
}

@test "Verifying the second Bastion host's boot disk" {
    run gcloud compute disks describe ${BASTION2_NAME} \
        --zone ${ZONE} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "sizeGb: '${BASTION2_DISK_SIZE}'" ]]
}

@test "Verifying the second Bastion's sudo is OFF" {
    # Wait until VM provisioning finishes
    until gcloud compute instances get-serial-port-output ${BASTION2_NAME} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
        --zone ${ZONE} | grep ${PROVISION_COMPLETED_MARKER}; do

        sleep 10;
    done

    run gcloud compute ssh ${BASTION2_NAME} --command "sudo -n whoami" \
        --zone ${ZONE} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ ! "$status" -eq 0 ]]
}

@test "Verifying the second Bastion's tags" {
    run gcloud compute instances describe ${BASTION2_NAME} \
        --format "yaml(tags)" \
        --zone ${ZONE} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "- ${BASTION2_EXTRA_TAG}" ]]
    [[ "$output" =~ "- ${BASTION2_TAG}" ]]
}

@test "Verifying Bastion's inbound firewall rule" {
    run gcloud compute firewall-rules describe "${SSH_TO_BASTION_RULE_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "IPProtocol: tcp" ]]
    [[ "$output" =~ "- '22'" ]]
    [[ "$output" =~ "direction: INGRESS" ]]
    [[ "$output" =~ "disabled: false" ]]
    [[ "$output" =~ "${NETWORK_NAME}" ]]
    [[ "$output" =~ "priority: ${SSH_TO_BASTION_PRIORITY}" ]]
}

@test "Verifying Bastion's inbound firewall rule's source range" {
    run gcloud compute firewall-rules describe "${SSH_TO_BASTION_RULE_NAME}" \
        --format="yaml(sourceRanges)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${SSH_TO_BASTION_SOURCE_RANGE}" ]]
}

@test "Verifying Bastion's inbound firewall rule's source tag" {
    run gcloud compute firewall-rules describe "${SSH_TO_BASTION_RULE_NAME}" \
        --format="yaml(sourceTags)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${SSH_TO_BASTION_SOURCE_TAG}" ]]
}

@test "Verifying Bastion's inbound firewall rule's target tag" {
    run gcloud compute firewall-rules describe "${SSH_TO_BASTION_RULE_NAME}" \
        --format="yaml(targetTags)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${BASTION2_TAG}" ]]
}

@test "Verifying Bastion's outbound firewall rule" {
    run gcloud compute firewall-rules describe \
        "${SSH_FROM_BASTION_DEFAULT_RULE_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "IPProtocol: tcp" ]]
    [[ "$output" =~ "- '22'" ]]
    [[ "$output" =~ "direction: INGRESS" ]]
    [[ "$output" =~ "disabled: false" ]]
    [[ "$output" =~ "${NETWORK_NAME}" ]]
    [[ "$output" =~ "priority: ${SSH_FROM_BASTION_PRIORITY}" ]]
}

@test "Verifying Bastion's outbound firewall rule's source tag" {
    run gcloud compute firewall-rules describe \
        "${SSH_FROM_BASTION_DEFAULT_RULE_NAME}" \
        --format="yaml(sourceTags)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${BASTION2_TAG}" ]]
}

@test "Verifying Bastion's outbound firewall rule's target tag" {
    run gcloud compute firewall-rules describe \
        "${SSH_FROM_BASTION_DEFAULT_RULE_NAME}" \
        --format="yaml(targetTags)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${SSH_FROM_BASTION_SOURCE_TAG}" ]]
}

@test "Deleting deployment" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

