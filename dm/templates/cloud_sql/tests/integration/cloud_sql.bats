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
    export MASTER_INSTANCE_NAME="cloud-sql-master-instance-${RAND}"
    export VERSION="MYSQL_5_6"
    export MASTER_INSTANCE_TIER="db-n1-standard-1"
    export MASTER_ZONE="us-central1-c"
    export REPLICA_ZONE="us-central1-a"
    export REGION="us-central1"
    export REPLICA_INSTANCE_NAME="cloud-sql-replica-instance-${RAND}"
    export REPLICA_INSTANCE_TIER="db-n1-standard-2"
    export REPLICA_INSTANCE_TYPE="READ_REPLICA_INSTANCE"
    export BACKUP_START_TIME="02:00"
    export BACKUP_ENABLED="true"
    export BACKUP_BL_ENABLED="true"
    export USER1_NAME="user-1"
    export USER1_HOST="10.1.1.1"
    export USER2_NAME="user-2"
    export USER2_HOST="10.1.1.2"
    export DB1="db-1"
    export DB2="db-2"
fi

########## HELPER FUNCTIONS ##########

function create_config() {
    envsubst < ${BATS_TEST_DIRNAME}/${TEST_NAME}.yaml > "${CONFIG}"
}

function delete_config() {
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
        rm -f "${RANDOM_FILE}"
        delete_config
    fi

    # Per-test teardown steps.
}


@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    run gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config "${CONFIG}" --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

@test "Verifying that both instances were created" {
    run gcloud sql instances list \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${MASTER_INSTANCE_NAME}" ]]
    [[ "$output" =~ "${REPLICA_INSTANCE_NAME}" ]]
}

@test "Verifying master instance" {
    run gcloud sql instances describe ${MASTER_INSTANCE_NAME} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${VERSION}" ]]
    [[ "$output" =~ "${MASTER_INSTANCE_TIER}" ]]
    [[ "$output" =~ "instanceType: CLOUD_SQL_INSTANCE" ]]
    [[ "$output" =~ "region: ${REGION}" ]]
    [[ "$output" =~ "${MASTER_ZONE}" ]]
}

@test "Verifying master replica list" {
    run gcloud sql instances describe ${MASTER_INSTANCE_NAME} \
        --format="yaml(replicaNames)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${REPLICA_INSTANCE_NAME}" ]]
}

@test "Verifying master database list" {
    run gcloud sql databases list --instance ${MASTER_INSTANCE_NAME} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${DB1}" ]]
    [[ "$output" =~ "${DB2}" ]]
}

@test "Verifying master user list" {
    run gcloud sql users list --instance ${MASTER_INSTANCE_NAME} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${USER1_NAME}" ]]
    [[ "$output" =~ "${USER2_NAME}" ]]
    [[ "$output" =~ "${USER1_HOST}" ]]
    [[ "$output" =~ "${USER2_HOST}" ]]
}

@test "Verifying master backup settings" {
    run gcloud sql instances describe ${MASTER_INSTANCE_NAME} \
        --format="yaml(settings.backupConfiguration)" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "binaryLogEnabled: ${BACKUP_BL_ENABLED}" ]]
    [[ "$output" =~ "enabled: ${BACKUP_ENABLED}" ]]
    [[ "$output" =~ "startTime: ${BACKUP_START_TIME}" ]]
}

@test "Verifying replica instance" {
    run gcloud sql instances describe ${REPLICA_INSTANCE_NAME} \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${VERSION}" ]]
    [[ "$output" =~ "${REPLICA_INSTANCE_TIER}" ]]
    [[ "$output" =~ "${REPLICA_INSTANCE_TYPE}" ]]
    [[ "$output" =~ "region: ${REGION}" ]]
    [[ "$output" =~ "${REPLICA_ZONE}" ]]
    [[ "$output" =~ "${MASTER_INSTANCE_NAME}" ]]
}

@test "Deleting deployment" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}
