# Copyright 2018 Google Inc. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""
This template creates an organization policy to allow VMs to have public
IPs.
"""


def generate_config(context):
    """ Entry point for the deployment resources. """

    project = context.properties['projectId']
    resources = []

    for policy in context.properties['policies']:
        constraint_name = policy['constraint'].split('/')[1]
        resources.append(
            {
                'name': 'set-{}-{}'.format(project, constraint_name),
                'action': 'gcp-types/cloudresourcemanager-v1:cloudresourcemanager.projects.setOrgPolicy',  # pylint: disable=line-too-long
                'metadata': {'runtimePolicy': ['CREATE']},
                'properties': {
                    'resource': 'projects/{}'.format(project),
                    'policy': policy
                }
            }
        )
        resources.append(
            {
                'name': 'clear-{}-{}'.format(project, constraint_name),
                'action': 'gcp-types/cloudresourcemanager-v1:cloudresourcemanager.projects.clearOrgPolicy',  # pylint: disable=line-too-long
                'metadata': {'runtimePolicy': ['DELETE']},
                'properties': {
                    'resource': 'projects/{}'.format(project),
                    'constraint': policy['constraint']
                }
            }
        )
    return {'resources': resources}
