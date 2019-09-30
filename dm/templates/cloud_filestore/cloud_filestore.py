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
""" This template creates a Google Cloud Filestore instance. """


def generate_config(context):
    """ Entry point for the deployment resources. """

    resources = []
    properties = context.properties
    project_id = properties.get('project', context.env['project'])
    name = properties.get('name', context.env['name'])

    resource = {
        'name': context.env['name'],
        # https://cloud.google.com/filestore/docs/reference/rest/v1beta1/projects.locations.instances/create
        'type': 'gcp-types/file-v1beta1:projects.locations.instances',
        'properties': {
            'parent': 'projects/{}/locations/{}'.format(project_id, properties['location']),
            'instanceId': name,
        }
    }

    optional_props = [
        'description',
        'tier',
        'labels',
        'fileShares',
        'networks',
    ]

    for prop in optional_props:
        if prop in properties:
            resource['properties'][prop] = properties[prop]

    resources.append(resource)

    return {
        'resources':
            resources,
        'outputs':
            [
                {
                    'name': 'name',
                    'value': name
                },
                {
                    'name': 'fileShares',
                    'value': '$(ref.{}.fileShares)'.format(context.env['name'])
                },
                {
                    'name': 'networks',
                    'value': '$(ref.{}.networks)'.format(context.env['name'])
                }
            ]
    }
