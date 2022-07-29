#!/usr/bin/env bash

# Copyright 2022 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Prerequisites:
#   Install and Authenticate Github CLI (https://github.com/cli/cli)
#
# Usage of merge_dependabot_prs.sh
#   -o [GitHub organization | Default: terraform-google-modules]
#   -f [Repository name(s) contains filter | Default: NONE]
#   -l [label to apply to failed checks | Default: dependabot-checks-failed]
#

# Default Variables
ORG="terraform-google-modules"
FILTER=".[].name"
LABEL="dependabot-checks-failed"

# Process arguments
while getopts 'o:f:l:n' arg; do
  case $arg in
    o) ORG=${OPTARG} ;;
    f) FILTER=".[]|select(.name|contains(\"$OPTARG\"))|.name" ;;
    l) LABEL=${OPTARG} ;;
  esac
done

# Initialize Variables
FPRS=()
PPRS=()

# Check gh is installed
if [ ! -x $(which gh) ]; then
  echo "GitHub CLI (gh) not found - Install from https://github.com/cli/cli"
  exit 1
fi

# Check gh is authenticated
gh auth status > /dev/null 2>&1
if [ ! $? -eq 0 ]; then
  echo "Please authenticate GitHub CLI with 'gh auth login' prior to running"
  exit 1
fi

# Retrieve list of repos in the Org
REPOS=`gh repo list $ORG --no-archived --json name -q $FILTER -L 1000`

# Confirm we retrieved repos
if [[ ! -n $REPOS ]]; then
  echo "No repos found"
  exit
fi

# Process the repos
for REPO in $REPOS; do
  REPO=$ORG/$REPO
  echo "Processing $REPO"

  # Retrieve Pull Requests
  PRS=`gh pr list -R $REPO -s open --json number -q '.[].number' --app dependabot`
  if [[ -z $PRS ]]; then
    echo "  No open Dependabot Pull Requests found for $REPO"
    continue
  fi

  # Process Pull Rquests
  for PR in $PRS; do
    echo "  Processing Dependabot Pull Request $PR"
    # Check status of Pull Request Checks
    gh pr checks $PR -R $REPO
    if [ $? -eq 0 ]; then
      PPRS+=($PR)
      # Remove the label, if exists
      gh pr edit $PR -R $REPO --remove-label $LABEL
      # Approve the Pull Request
      gh pr review $PR -R $REPO --approve -b "LGTM"
      # Squash Merge the Pull Request and Delete the Branch
      gh pr merge $PR -d -s -R $REPO
    else
      FPRS+=($PR)
      # Create the Label, if not exist
      gh label create $LABEL --color E99695 -R $REPO > /dev/null 2>&1
      # Add the Label to the Pull Request
      gh pr edit $PR -R $REPO --add-label $LABEL
    fi
  done
done

# List number of approved/merged PRs
if [[ ! -z $PPRS ]]; then
  echo -e "\u2714 approved and merged ${PPRS[@]} Pull Requests"
fi

# List Failed Checks PRs
if [[ ! -z $FPRS ]]; then
  echo "These Pull Requests have failed checks and have been labled with $LABEL:"
  for PR in ${FPRS[@]}; do
    echo -e "\u274c $PR: https://github.com/$REPO/pull/$PR"
  done
  echo ""
  echo "View all $ORG PRs in Github with $LABEL at: https://github.com/pulls?q=is%3Aopen+is%3Apr+author%3Aapp%2Fdependabot+archived%3Afalse+org%3A$ORG+label%3A$LABEL"
fi

