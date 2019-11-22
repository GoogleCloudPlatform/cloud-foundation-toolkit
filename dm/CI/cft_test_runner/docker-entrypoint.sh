#!/bin/bash
set -eu

readonly GIT_URL='https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit'
readonly CLONE_DIRNAME="/workspace"
readonly DM_root="/cloud-foundation-toolkit/dm"

readonly COLOR_RESET='\033[0m'
readonly COLOR_BOLD='\033[1m'
readonly COLOR_BG_BLUE='\033[44m'

echo_color() {
  echo -e "${COLOR_BOLD}${COLOR_BG_BLUE}$@${COLOR_RESET}"
}

echo_color 'Activating venv for testing'

cd "${DM_root}"

set +u  # Turn off because virtualenv uses undefined variables
. venv/bin/activate \
./src/cftenv 
set -u

export CLOUD_FOUNDATION_CONF=/etc/cloud-foundation-tests.conf

if [ -d "/workspace/dm" ]
then
    echo_color "/workspace/dm exists, no need to clone repo"
else
    echo_color "/workspace/dm is missing, cloning repo"
    readonly BRANCH_NAME=$2
    git clone "${GIT_URL}" "${CLONE_DIRNAME}"
    cd "${CLONE_DIRNAME}"
    git checkout "${BRANCH_NAME}"
fi

mv "/workspace/dm/templates"  "${DM_root}"

echo_color "Welcome your Majesty, ready to run some tests!"

# Running bats tests relative to dm folder for example: "./templates/project/tests/integration/project.bats"

cd "${DM_root}"

chmod 777 $1
exec bats $1
