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


def append_optional_property(res, properties, prop_name):
    """ If the property is set, it is added to the resource. """

    val = properties.get(prop_name)
    if val:
        res['properties'][prop_name] = val
    return

def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    name = properties.get('name', context.env['name'])
    project_id = properties.get('project', context.env['project'])

    bgp = properties.get('bgp', {'asn': properties.get('asn')})

    router = {
        'name': context.env['name'],
        # https://cloud.google.com/compute/docs/reference/rest/v1/routers
        'type': 'gcp-types/compute-v1:routers',
        'properties':
            {
                'name':
                    name,
                'project':
                    project_id,
                'region':
                    properties['region'],
                'network':
                    properties.get('networkURL', generate_network_uri(
                        project_id,
                        properties.get('network', ''))),
            }
    }

    if properties.get('bgp'):
        router['properties']['bgp'] = bgp

    optional_properties = [
        'description',
        'bgpPeers',
        'interfaces',
        'nats',
    ]

    for prop in optional_properties:
        append_optional_property(router, properties, prop)

    return {
        'resources': [router],
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


def generate_network_uri(project_id, network):
    """Format the network name as a network URI."""

    return 'projects/{}/global/networks/{}'.format(
        project_id,
        network
    )
