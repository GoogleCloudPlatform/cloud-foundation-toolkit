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
"""This template grants IAM roles to a user on a shared VPC subnetwork."""


def _append_resource(subnets, project, name_id):
    """Append subnets to resources."""
    resources = []
    out = {}
    for subnet in subnets:
        policy_name = 'iam-subnet-policy-{}'.format(subnet[name_id])
        resources.append({
            'name': policy_name,
            # https://cloud.google.com/compute/docs/reference/rest/beta/subnetworks/setIamPolicy
            'type': 'gcp-types/compute-beta:compute.subnetworks.setIamPolicy',
            'properties': {
                'name': subnet[name_id],
                'project': project,
                'region': subnet['region'],
                'bindings': [{
                    'role': subnet['role'],
                    'members': subnet['members']
                }]
            }
        })

        out[policy_name] = {
            'etag': '$(ref.' + policy_name + '.etag)'
        }
    return resources, out


def generate_config(context):
    """Entry point for the deployment resources."""
    try:
        resources, out = _append_resource(
            context.properties['subnets'],  # Legacy syntax
            context.env['project'],
            'subnetId'
        )
    except KeyError:
        try:
            resources, out = _append_resource(
                context.properties['policy']['bindings'],  # Policy syntax
                context.env['project'],
                'resourceId'
            )
        except KeyError:
            resources, out = _append_resource(
                context.properties['bindings'],  # Bindings syntax
                context.env['project'],
                'resourceId'
            )
    outputs = [{'name': 'policies', 'value': out}]

    return {'resources': resources, 'outputs': outputs}
