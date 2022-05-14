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
# Usage of update_branch_protection.sh
#   -b [branch | Default: main]
#   -o [GitHub organization | Default: terraform-google-modules]
#   -f [Repository name(s) contains filter | Default: NONE]
#

set -e

# Default Variables
ORG="terraform-google-modules"
BRANCH="main"
FILTER=".[].name"

# Process arguments
while getopts 'b:o:f:' arg; do
  case $arg in
    b) BRANCH=${OPTARG} ;;
    o) ORG=${OPTARG} ;;
    f) FILTER=".[]|select(.name|contains(\"$OPTARG\"))|.name" ;;
  esac
done

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
  echo "Updating $ORG/$REPO"

  # Retrieve any current checks
  CHECKS=`gh api \
  -H "Accept: application/vnd.github.v3+json" \
  /repos/$ORG/$REPO/branches/$BRANCH/protection/required_status_checks/contexts`

  # Update the branch protection, include any current checks
  jq -n '{"required_pull_request_reviews":{"required_approving_review_count":1},"required_status_checks":{"strict":true,"contexts":'"$CHECKS"'},"enforce_admins":true,"restrictions":{"teams":[],"users":[]}}' | \
  gh api \
  --method PUT \
  -H "Accept: application/vnd.github.v3+json" \
  /repos/$ORG/$REPO/branches/$BRANCH/protection \
  --silent --input -

  # Enable only squash commits
  gh repo edit $ORG/$REPO --enable-squash-merge --enable-merge-commit=false --enable-rebase-merge=false

done

