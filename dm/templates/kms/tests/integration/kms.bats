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
    # test specific variables
    export KEYRING_NAME="test-keyring-${RAND}"
    export REGION="global"
    export KEY_NAME="test-key-${RAND}"
    export SA_NAME="test-kms-${RAND}"
    export SA_FQDN="${SA_NAME}@${CLOUD_FOUNDATION_PROJECT_ID}.iam.gserviceaccount.com"
    export ROLE="roles/cloudkms.admin"
    export KEY_PURPOSE="ENCRYPT_DECRYPT"
    # export NEXT_ROTATION_TIME=$(date -d '2 months' '+%Y-%m-%dT%H:%M:%S.%NZ')
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
        # Create service accounts to test IAM bindings.
        gcloud iam service-accounts create "${SA_NAME}" \
            --display-name "Test KMS Service Account"
    fi

    # Per-test setup steps.
}

function teardown() {
    # Global teardown; executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
        delete_config
        rm -f "${RANDOM_FILE}"
        # Delete service account after tests had been completed.
        gcloud --quiet iam service-accounts delete "${SA_FQDN}"
    fi

    # Per-test teardown steps.
}

########## TESTS ##########

@test "Creating deployment ${DEPLOYMENT_NAME} from ${CONFIG}" {
    run gcloud deployment-manager deployments create "${DEPLOYMENT_NAME}" \
        --config "${CONFIG}" \
        --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$status" -eq 0 ]]
}

sleep 5

@test "Verifying the KeyRing ${KEYRING_NAME} was created " {
    run gcloud kms keyrings list --location ${REGION} \
        --format="value(name.scope(keyRings))" \
        --filter="${KEYRING_NAME}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${KEYRING_NAME}" ]]
}

@test "KEY ${KEY_NAME} is created in KeyRing ${KEYRING_NAME} " {
    run gcloud kms keys list --location ${REGION} --keyring="${KEYRING_NAME}" \
        --format="value(name.scope(cryptoKeys))" \
        --filter="${KEY_NAME}"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${KEY_NAME}" ]]
}

@test "CryptoKey's PURPOSE is set to ${KEY_PURPOSE} " {
    run gcloud kms keys describe ${KEY_NAME} --location ${REGION} \
        --keyring="${KEYRING_NAME}" \
        --format="value(purpose)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${KEY_PURPOSE}" ]]
}

@test "protectionLevel on key ${KEY_NAME} = SOFTWARE " {
    run gcloud kms keys describe ${KEY_NAME} --location ${REGION} \
        --keyring="${KEYRING_NAME}" \
        --format="value(versionTemplate.protectionLevel)"
    [[ "$status" -eq 0 ]]
    [[ "$output" -eq "SOFTWARE" ]]
}

@test "Enc algorithm on key ${KEY_NAME} = GOOGLE_SYMMETRIC_ENCRYPTION" {
    run gcloud kms keys describe ${KEY_NAME} --location ${REGION} \
        --keyring="${KEYRING_NAME}" \
        --format="value(versionTemplate.algorithm)"
    [[ "$status" -eq 0 ]]
    [[ "$output" -eq "GOOGLE_SYMMETRIC_ENCRYPTION" ]]
}

@test "Verify whether ${SA_NAME} has role ${ROLE} " {
    run gcloud kms keys get-iam-policy ${KEY_NAME} --location ${REGION} \
        --keyring="${KEYRING_NAME}" \
        --format="value(bindings.role)"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~  "${ROLE}" ]]
}

@test "Verify if ${SA_NAME} has access to ${KEY_NAME} " {
    run gcloud kms keys get-iam-policy ${KEY_NAME} --location ${REGION} \
        --keyring="${KEYRING_NAME}" \
        --format="value(bindings.members[0])"
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "${SA_FQDN}" ]]
}

########### NOTE ##################
# There is no Delete Deployment step because KeyRings, Keys cannot be deleted.
##################################
