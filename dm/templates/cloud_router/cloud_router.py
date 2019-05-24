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
""" This template creates a Cloud Router. """


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    name = properties.get('name', context.env['name'])
    project_id = properties.get('project', context.env['project'])

    resources = [
        {
            'name': context.env['name'],
            'type': 'gcp-types/compute-v1:routers',
            'properties':
                {
                    'name':
                        name,
                    'bgp': {
                        'asn': context.properties['asn']
                    },
                    'network':
                        generate_network_url(
                            project_id,
                            context.properties['network']
                        ),
                    'region':
                        context.properties['region'],
                    'project':
                        project_id
                }
        }
    ]

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
                    'name': 'selfLink',
                    'value': '$(ref.' + context.env['name'] + '.selfLink)'
                },
                {
                    'name':
                        'creationTimestamp',
                    'value':
                        '$(ref.' + context.env['name'] + '.creationTimestamp)'
                }
            ]
    }


def generate_network_url(project_id, network):
    """Format the resource name as a resource URI."""

    return 'projects/{}/global/networks/{}'.format(
        project_id,
        network
    )
