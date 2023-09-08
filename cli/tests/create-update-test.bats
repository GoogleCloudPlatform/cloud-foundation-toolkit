#!/usr/bin/env bats

source ../../dm/tests/helpers.bash

########## HELPER FUNCTIONS ##########

function setup() {
    # Global setup; this is executed once per test file.
    if [ ${BATS_TEST_NUMBER} -eq 1 ]; then
    	echo "setup"
    fi
    # Per-test setup steps.
}

function teardown() {
    # Global teardown; this is executed once per test file.
    if [[ "$BATS_TEST_NUMBER" -eq "${#BATS_TEST_NAMES[@]}" ]]; then
    	echo "teardown"
    fi
    # Per-test teardown steps.
}

########## TESTS ##########

@test "Creating deployments" {
    ../bin/cft create ./create --project "${CLOUD_FOUNDATION_PROJECT_ID}"
}


@test "Verifying that resources were created in deployment" {
    run gcloud compute networks list --filter="name:cftcli-test-network" --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "cftcli-test-network" ]]
}

@test "Verifying subnets were created" {
    run gcloud compute networks subnets list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "cftcli-test-subnetwork-1" ]]
    [[ "$output" =~ "cftcli-test-subnetwork-2" ]]
    [[ "$output" =~ "cftcli-test-subnetwork-3" ]]
}

@test "Verifying firewall rules were created" {
    run gcloud compute firewall-rules list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "allow-proxy-from-inside" ]]
    [[ "$output" =~ "443" ]]
    [[ "$output" =~ "allow-dns-from-inside" ]]
}

@test "Update deployments" {
    ../bin/cft update ./update --project "${CLOUD_FOUNDATION_PROJECT_ID}"
}

@test "Verifying one subnet was removed by update operation" {
    run gcloud compute networks subnets list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "cftcli-test-subnetwork-1" ]]
    [[ "$output" =~ "cftcli-test-subnetwork-2" ]]
    [[ ! "$output" =~ "cftcli-test-subnetwork-3" ]]
}

@test "Verifying one firewall rule was removed one updated" {
    run gcloud compute firewall-rules list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ "$output" =~ "allow-proxy-from-inside" ]]
    [[ ! "$output" =~ "443" ]]
    [[ ! "$output" =~ "allow-dns-from-inside" ]]
}

@test "Deleting deployments" {
     ../bin/cft delete ./create --project "${CLOUD_FOUNDATION_PROJECT_ID}"

    run gcloud compute networks list --filter="name:cftcli-test-network" --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ ! "$output" =~ "cftcli-test-network" ]]

    run gcloud compute networks subnets list --project "${CLOUD_FOUNDATION_PROJECT_ID}"
    [[ ! "$output" =~ "cftcli-test-subnetwork-1" ]]
    [[ ! "$output" =~ "cftcli-test-subnetwork-2" ]]
}
