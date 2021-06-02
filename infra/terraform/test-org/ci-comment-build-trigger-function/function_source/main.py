# Copyright 2020 Google LLC
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

import base64
import sys
import os
import json
import logging
import requests

from google.cloud.devtools.cloudbuild_v1 import CloudBuildClient as cloudbuild
from google.cloud.devtools.cloudbuild_v1.types import BuildStep, Build, BuildOptions
from google.protobuf import duration_pb2 as duration

CFT_TOOLS_DEFAULT_IMAGE = 'gcr.io/cloud-foundation-cicd/cft/developer-tools'
CFT_TOOLS_DEFAULT_IMAGE_VERSION = '0.12'
DISABLED_MODULES = ["terraform-example-foundation"]


def main(event, context):
    """ Triggers a new downstream build based on a PubSub message originating from a parent cloudbuild """
    # if cloud build project is not set, exit
    if not os.getenv('CLOUDBUILD_PROJECT'):
        logging.warn('Cloud Build project not set')
        sys.exit(1)
    # if no data in PubSub event, log and exit
    if 'data' not in event:
        logging.info('Unable to find data in PubSub event')
        sys.exit(1)
    # decode data in PubSub event
    data = json.loads(base64.b64decode(event['data']).decode('utf-8'))
    # if the parent build originated from CF, ignore
    if data['substitutions'].get('_IS_TRIGGERED_BY_CF', False):
        logging.warn('Triggered by CF, Ignoring')
        return
    logging.info('Parent build not triggered by CF')
    # if parent build is not a lint build, ignore
    if 'lint' not in data['tags']:
        logging.warn('Parent build is not a lint build')
        return
    # if parent build has not started, or is in any other state, ignore
    if data['status'] != 'WORKING':
        logging.warn('Parent build is not in WORKING status')
        return
    logging.info('Parent build is in WORKING status')
    # if repo ref for the parent build has disabled PR bot, ignore
    if data['substitutions']['REPO_NAME'] in DISABLED_MODULES:
        logging.warn('Comment bot is disabled for this repo')
        return
    if data['substitutions'].get('_DOCKER_TAG_VERSION_DEVELOPER_TOOLS', False):
        logging.info(
            f'Found _DOCKER_TAG_VERSION_DEVELOPER_TOOLS. Setting tools image version to {data["substitutions"]["_DOCKER_TAG_VERSION_DEVELOPER_TOOLS"]}'
        )
        CFT_TOOLS_DEFAULT_IMAGE_VERSION = data['substitutions'][
            '_DOCKER_TAG_VERSION_DEVELOPER_TOOLS'
        ]
    # Cloud Build seems to have a bug where if a build is re run through Github UI, it will not set _PR_NUMBER or _HEAD_REPO_URL
    # workaround using the GH API to infer PR number and _HEAD_REPO_URL
    PR_NUMBER = data['substitutions'].get('_PR_NUMBER', False)
    _HEAD_REPO_URL = data['substitutions'].get('_HEAD_REPO_URL', False)
    # default clone repo step
    get_repo_args = [
        '-c',
        'git clone $$REPO_URL . && git checkout $$COMMIT_SHA && git status',
    ]
    if not (PR_NUMBER or _HEAD_REPO_URL):
        logging.warn('Unable to infer PR number via Cloud Build. Trying via GH API')
        # get list of github PRs that have this SHA
        response = requests.get(
            f'https://api.github.com/search/issues?q={data["substitutions"]["COMMIT_SHA"]}'
        )
        response.raise_for_status()
        response_obj = response.json()
        # if more than one PR, ignore
        if response_obj['total_count'] != 1:
            logging.info(f'Multiple associated PRs found. Exiting...')
            return
        # if only one PR, its safe to assume that is associated with parent build's PR
        logging.info(f'One associated PR found: {response_obj["items"][0]["number"]}')
        PR_NUMBER = response_obj['items'][0]['number']
        # get target repo URL
        pr_url = response_obj['items'][0]['html_url']
        _HEAD_REPO_URL = pr_url[: pr_url.find('/pull')]
        # fetch PR at head using PR number
        get_repo_args = [
            '-c',
            'git clone $$REPO_URL . && git fetch origin pull/$$_PR_NUMBER/head:$$_PR_NUMBER && git checkout $$_PR_NUMBER && git show --name-only',
        ]

    # prepare env vars
    env = [
        f'_PR_NUMBER={PR_NUMBER}',
        f'REPO_NAME={data["substitutions"]["REPO_NAME"]}',
        f'REPO_URL={_HEAD_REPO_URL}',
        f'COMMIT_SHA={data["substitutions"]["COMMIT_SHA"]}',
    ]
    get_repo_step = BuildStep(
        name='gcr.io/cloud-builders/git',
        env=env,
        args=get_repo_args,
        id='get_repo',
        entrypoint='bash',
    )
    # lint comment step
    lint_args = [
        '-c',
        'source /usr/local/bin/task_helper_functions.sh && printenv && post_lint_status_pr_comment',
    ]
    lint_step = BuildStep(
        name=f'{CFT_TOOLS_DEFAULT_IMAGE}:{CFT_TOOLS_DEFAULT_IMAGE_VERSION}',
        env=env,
        args=lint_args,
        id='lint_comment',
        entrypoint='/bin/bash',
    )
    # substitutions
    sub = {
        '_IS_TRIGGERED_BY_CF': '1',
    }
    # create and trigger build
    build = Build(
        steps=[get_repo_step, lint_step],
        options=BuildOptions(substitution_option='ALLOW_LOOSE'),
        substitutions=sub,
        timeout=duration.Duration(seconds=1200),
    )
    response = cloudbuild().create_build(os.getenv('CLOUDBUILD_PROJECT'), build)
    logging.info(response)
