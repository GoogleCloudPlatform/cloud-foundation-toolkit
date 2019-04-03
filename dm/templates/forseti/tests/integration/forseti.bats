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
    export BUCKET_NAME="forseti-security-bucket-${RAND}"
    export SERVER_ZONE="us-central1-a"
    export CLIENT_ZONE="us-central1-a"
    export PROJECT_ID="forseti-project-${RAND}"
    export PROJECT_NAME="Forseti Security-${RAND}"
    export SQL_NAME="forseti-sql-instance-${RAND}"
    export SQL_DB_NAME="${SQL_NAME}-db"
    export SQL_REGION="us-central1"
    export SERVER_SA_PREFIX="forseti-server-gcp"
    export CLIENT_SA_PREFIX="forseti-client-gcp"
    export SERVER_NAME="forseti-server"
    export CLIENT_NAME="forseti-client"
    export CLOUD_FOUNDATION_FOLDER_NAME="test-forseti-folder-${RAND}"
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
        # Set up instance groups to be load-balanced via HAProxy.
        gsutil mb gs://${BUCKET_NAME}
        gcloud alpha resource-manager folders create \
            --display-name="${CLOUD_FOUNDATION_FOLDER_NAME}" \
            --organization="${CLOUD_FOUNDATION_ORGANIZATION_ID}" > ~/output.txt
        export CLOUD_FOUNDATION_FOLDER_ID=$(gcloud alpha resource-manager folders list \
            --project "${CLOUD_FOUNDATION_PROJECT_ID}" \
            --organization "${CLOUD_FOUNDATION_ORGANIZATION_ID}" | \
            grep "${CLOUD_FOUNDATION_FOLDER_NAME}" | \
            awk '{print $3}')
        gcloud alpha resource-manager folders list --organization="${CLOUD_FOUNDATION_ORGANIZATION_ID}" >> ~/output.txt
        create_config
    fi

  # Per-test setup steps.
}

function teardown() {
    # Global teardown; executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        rm -f "${RANDOM_FILE}"
        gcloud alpha resource-manager folders delete \
            "${CLOUD_FOUNDATION_FOLDER_ID}"
        gsutil rb gs://${BUCKET_NAME}
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

@test "Verifying that new project exists" {
    run gcloud projects list
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${PROJECT_ID}" ]]
    [[ "$output" =~ "${PROJECT_NAME}" ]]
}

@test "Verifying server's service account" {
    run gcloud iam service-accounts list --project ${PROJECT_ID}
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${SERVER_SA_PREFIX}" ]]

    run gcloud projects get-iam-policy ${PROJECT_ID}
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "roles/cloudsql.client" ]]
    [[ "$output" =~ "roles/logging.logWriter" ]]
    [[ "$output" =~ "roles/storage.objectCreator" ]]
    [[ "$output" =~ "roles/storage.objectViewer" ]]

    run gcloud organizations get-iam-policy \
        ${CLOUD_FOUNDATION_ORGANIZATION_ID}
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "roles/appengine.appViewer" ]]
    [[ "$output" =~ "roles/bigquery.dataViewer" ]]
    [[ "$output" =~ "roles/browser" ]]
    [[ "$output" =~ "roles/cloudasset.viewer" ]]
    [[ "$output" =~ "roles/cloudsql.viewer" ]]
    [[ "$output" =~ "roles/compute.networkViewer" ]]
    [[ "$output" =~ "roles/compute.securityAdmin" ]]
    [[ "$output" =~ "roles/iam.securityReviewer" ]]
    [[ "$output" =~ "roles/servicemanagement.quotaViewer" ]]
    [[ "$output" =~ "roles/serviceusage.serviceUsageConsumer" ]]
}

@test "Verifying client's service account" {
    run gcloud iam service-accounts list --project ${PROJECT_ID}
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${CLIENT_SA_PREFIX}" ]]

    run gcloud projects get-iam-policy ${PROJECT_ID}
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "roles/logging.logWriter" ]]
    [[ "$output" =~ "roles/storage.objectViewer" ]]
}

@test "Verifying sql instance" {
    run gcloud sql instances list \
        --project ${PROJECT_ID}
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${SQL_NAME}" ]]
    [[ "$output" =~ "MYSQL_5_7" ]]
}

@test "Verifying sql database" {
    run gcloud sql databases list \
        --instance ${SQL_NAME} \
        --project ${PROJECT_ID}
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${SQL_DB_NAME}" ]]
}

@test "Verifying Forseti server" {
    run gcloud compute instances list \
        --project ${PROJECT_ID}
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${SERVER_NAME}" ]]

    run gcloud compute instances describe ${SERVER_NAME} \
        --project ${PROJECT_ID} \
        --zone ${SERVER_ZONE}
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${SQL_DB_NAME}" ]]
    [[ "$output" =~ "${SQL_NAME}" ]]
    [[ "$output" =~ "${BUCKET_NAME}" ]]
}

@test "Verifying Forseti client" {
    run gcloud compute instances list \
        --project ${PROJECT_ID}
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${CLIENT_NAME}" ]]

    run gcloud compute instances describe ${CLIENT_NAME} \
        --project ${PROJECT_ID} \
        --zone ${SERVER_ZONE}
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "forseti_conf_client.yaml" ]]
    [[ "$output" =~ "server_ip" ]]
}

@test "Deleting deployment" {
    run gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" -q \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}
