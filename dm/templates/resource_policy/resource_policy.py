# Copyright 2019 Google Inc. All rights reserved.
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
""" This template creates a Resource Policy. """


def generate_config(context):
    """ Entry point for the deployment resources. """

    resources = []
    properties = context.properties
    project_id = properties.get('project', context.env['project'])
    name = properties.get('name', context.env['name'])
    region = properties['region']
    resource_name = properties['resource']
    policy = properties['snapshotSchedulePolicy']

    resource = {
        'name': name,
        # https://cloud.google.com/compute/docs/reference/rest/v1/resourcePolicies/insert
        'type': 'gcp-types/compute-v1:resourcePolicies',
        'properties': {
            'project': project_id,
            'name': resource_name,
            'region': region,
            'snapshotSchedulePolicy': policy
        }
    }

    resources.append(resource)

    return {
        'resources':
            resources,
        'outputs':
            [
                {
                    'name': 'name',
                    'value': name
                }
            ]
    }
