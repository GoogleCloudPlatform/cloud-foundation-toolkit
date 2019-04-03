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
""" This template grants IAM roles to a user on a shared VPC subnetwork. """


def generate_config(context):
    """ Entry point for the deployment resources. """

    resources = []
    out = {}
    for subnet in context.properties['subnets']:
        subnet_id = subnet['subnetId']
        policy_name = 'iam-subnet-policy-{}'.format(subnet_id)

        policies_to_add = [
            {
                'role': subnet['role'],
                'members': subnet['members']
            }
        ]

        resources.append(
            {
                'name': policy_name,
                'type': 'gcp-types/compute-beta:compute.subnetworks.setIamPolicy',  # pylint: disable=line-too-long
                'properties':
                    {
                        'name': subnet_id,
                        'project': context.env['project'],
                        'region': subnet['region'],
                        'bindings': policies_to_add
                    }
            }
        )

        out[policy_name] = {
            'etag': '$(ref.' + policy_name + '.etag)'
        }

    outputs = [{'name': 'policies', 'value': out}]

    return {'resources': resources, 'outputs': outputs}
