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
    # Test specific variables:
    export BUCKET_NAME="test-bucket-${RAND}"
    export SA_NAME="${BUCKET_NAME}@${CLOUD_FOUNDATION_PROJECT_ID}.iam.gserviceaccount.com"
    export SA_FQDN="serviceAccount:${SA_NAME}"
    export ROLE="roles/storage.objectViewer"
    export LIFECYCLE_ACTION_TYPE="SetStorageClass"
    export LIFECYCLE_STORAGE_CLASS="NEARLINE"
    export LIFECYCLE_AGE_DAYS="36500"
    export LIFECYCLE_OBJ_CREATED_BEFORE="2018-01-01"
    export LIFECYCLE_NUM_NEWERVERSION="5"
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
        # create service accounts to test IAM bindings
        gcloud iam service-accounts create "${BUCKET_NAME}" \
            --display-name "Test Service Account"
    fi

    # Per-test setup steps.
}

function teardown() {
    # Global teardown; executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        delete_config
        rm -f "${RANDOM_FILE}"
        # delete service account after tests are complete.
        gcloud --quiet iam service-accounts delete "${SA_NAME}"
    fi

    # Per-test teardown steps.
}

########## TESTS ##########

@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config "${CONFIG}" --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

@test "Verify if Storage Bucket ${BUCKET_NAME} is created " {
    res=$(gsutil ls | grep "${BUCKET_NAME}")
    [[ "$status" -eq 0 ]]
    [[ "$res" =~ "gs://${BUCKET_NAME}/" ]]
}

@test "storageClass on ${BUCKET_NAME} is set to STANDARD " {
    run gsutil defstorageclass get "gs://${BUCKET_NAME}/"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "STANDARD" ]]
}

@test "Versioning on ${BUCKET_NAME} is ENABLED " {
    run gsutil versioning get "gs://${BUCKET_NAME}/"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "Enabled" ]]
}

@test "Logging configuration on ${BUCKET_NAME} is not set " {
    run gsutil logging get "gs://${BUCKET_NAME}/"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "has no logging configuration" ]]
}

@test "Verify if SA ${SA_NAME} has role ${ROLE}" {
    role=$(gsutil iam get "gs://${BUCKET_NAME}/" | grep role)
    [[ "$status" -eq 0 ]]
    [[ "$role" =~ "${ROLE}" ]]
}

@test "Verify if SA ${SA_NAME} is a member of this bucket" {
    member=$(gsutil iam get "gs://${BUCKET_NAME}/" | grep serviceAccount)
    [[ "$status" -eq 0 ]]
    [[ "$member" =~ "${SA_NAME}" ]]
}

@test "lifeCycle configuration is set on ${BUCKET_NAME}" {
    run gsutil lifecycle get "gs://${BUCKET_NAME}/"
    [[ "$status" -eq 0 ]]
    [[ ! "$output" =~ "has no lifecycle configuration" ]]
}

@test "lifeCycle Action Type is ${LIFECYCLE_ACTION_TYPE}" {
    lc_type=$(gsutil lifecycle get "gs://${BUCKET_NAME}/" | \
        grep ${LIFECYCLE_ACTION_TYPE})
    [[ "$status" -eq 0 ]]
    [[ "$lc_type" =~ "${LIFECYCLE_ACTION_TYPE}" ]]
}

@test "lifeCycle StorageClass is set to ${LIFECYCLE_STORAGE_CLASS}" {
    lc_class=$(gsutil lifecycle get "gs://${BUCKET_NAME}/" | \
        grep ${LIFECYCLE_STORAGE_CLASS})
    [[ "$status" -eq 0 ]]
    [[ "$lc_class" =~ "${LIFECYCLE_STORAGE_CLASS}" ]]
}

@test "lifeCycle Condition has AGE set to ${LIFECYCLE_AGE_DAYS}" {
    lc_age=$(gsutil lifecycle get "gs://${BUCKET_NAME}/" | \
        grep age)
    [[ "$status" -eq 0 ]]
    [[ "$lc_age" =~ "${LIFECYCLE_AGE_DAYS}" ]]
}

@test "lifeCycle Objects CreatedBefore Date is ${LIFECYCLE_OBJ_CREATED_BEFORE}" {
    lc_date=$(gsutil lifecycle get "gs://${BUCKET_NAME}/" | \
        grep createdBefore)
    [[ "$status" -eq 0 ]]
    [[ "$lc_date" =~ "${LIFECYCLE_OBJ_CREATED_BEFORE}" ]]
}

@test "lifeCycle numNewerVersions is ${LIFECYCLE_NUM_NEWERVERSION}" {
    lc_ver=$(gsutil lifecycle get "gs://${BUCKET_NAME}/" | \
        grep numNewerVersions)
    [[ "$status" -eq 0 ]]
    [[ "$lc_ver" =~ "${LIFECYCLE_NUM_NEWERVERSION}" ]]
}

@test "Deleting deployment ${DEPLOYMENT_NAME}" {
    gcloud deployment-manager deployments delete "${DEPLOYMENT_NAME}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}" -q
    [[ "$status" -eq 0 ]]

    run gsutil ls
    [[ "$status" -eq 0 ]]
    [[ ! "$output" =~ "gs://${BUCKET_NAME}/" ]]
}
