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
""" This template creates a Cloud Task queue. """


def generate_config(context):
    """ Entry point for the deployment resources. """

    resources = []
    properties = context.properties
    name = properties.get('name', context.env['name'])
    project_id = properties.get('projectId', context.env['project'])
    location = properties['location']
    parent = 'projects/{}/locations/{}'.format(project_id, location)

    queue = {
        'name': name,
        'type': '{}/cloudtasks:projects.locations.queues'.format(project_id),
        'properties': {
            'name': '{}/queues/{}'.format(parent, name),
            'parent': parent,
            'appEngineHttpQueue': properties['appEngineHttpQueue']
        }
    }

    optional_properties = ['rateLimits', 'retryConfig']

    for prop in optional_properties:
        if prop in properties:
            queue['properties'][prop] = properties[prop]

    resources.append(queue)

    outputs = [
        {
            'name': 'name',
            'value': '$(ref.{}.name)'.format(name)
        },
        {
            'name': 'state',
            'value': '$(ref.{}.state)'.format(name)
        }
    ]

    return {'resources': resources, 'outputs': outputs}
