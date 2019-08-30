#!/bin/bash
set -eu

readonly GIT_URL='https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit'
readonly CLONE_DIRNAME="$(mktemp -d)"
readonly BRANCH_NAME="cft-dm-dev"

readonly COLOR_RESET='\033[0m'
readonly COLOR_BOLD='\033[1m'
readonly COLOR_BG_BLUE='\033[44m'

echo_color() {
  echo -e "${COLOR_BOLD}${COLOR_BG_BLUE}$@${COLOR_RESET}"
}

echo_color "Cloning repo"
git clone "${GIT_URL}" "${CLONE_DIRNAME}"
cd "${CLONE_DIRNAME}"
git checkout "${BRANCH_NAME}"

echo_color 'Initializing CFT DM templates'

cd dm/templates

# cat healthcheck/examples/healthcheck.yaml | yq .resources[0].properties > project.json; cat healthcheck/healthcheck.py.schema | yq . > project.py.schema.json; ajv validate -s project.py.schema.json -d project.json

EXAMPLE_COUNT=`cat $@ | yq '.resources | length'`
EXAMPLE_COUNT=$(($EXAMPLE_COUNT-1))

while [ $EXAMPLE_COUNT -ge 0 ]; 
do
    echo_color "Example $EXAMPLE_COUNT"
    cat $@ | yq .resources[$EXAMPLE_COUNT].properties > example.json;
    cat example.json
    export SCHEMA_PATH=`cat $@ | yq -r .imports[0].path | awk '{print $1".schema"}'` 
    echo_color $SCHEMA_PATH
    cat $SCHEMA_PATH | yq . > example.py.schema.json;
    echo_color "Schema validation"
    ajv validate -s example.py.schema.json -d example.json
    EXAMPLE_COUNT=$(($EXAMPLE_COUNT-1))

done
