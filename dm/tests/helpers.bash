#!/bin/bash

# This file is meant to hold common variables and functions to be used by the
# testing suite (bats).
#
# Tests need to run against the user's own organization/projects/etc, so the
# most basic configs are read and exported from the user's own
# `~/.cloud-foundation-test.conf`.
#
# An example for this config is placed under `tests/cloud-foundation-tests.conf`. Users should
# move this file to `~/.cloud-foundation-test.conf` and tweak according to their own GCP
# organizational structure

CLOUD_FOUNDATION_CONF=${CLOUD_FOUNDATION_CONF-~/.cloud-foundation-tests.conf}

if [[ -z "${CLOUD_FOUNDATION_ORGANIZATION_ID}" || -z "${CLOUD_FOUNDATION_BILLING_ACCOUNT_ID}" || -z "${CLOUD_FOUNDATION_PROJECT_ID}" ]]; then
    if [[ ! -e ${CLOUD_FOUNDATION_CONF} ]]; then
        echo "Please setup your environment variables or Cloud Foundation config file"
        echo "Default location for config: ~/.cloud-foundation-tests.conf. Example:"
        echo "====================="
        cat tests/cloud-foundation-tests.conf.example
        echo "====================="
        exit 1
    fi
    source ${CLOUD_FOUNDATION_CONF}
fi
