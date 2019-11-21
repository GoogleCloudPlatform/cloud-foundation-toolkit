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
"""This template creates a custom route."""


from hashlib import sha1
import json


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    project_id = properties.get('project', context.env['project'])

    network_name = generate_network_url(properties)

    resources = []
    out = {}
    for i, route in enumerate(properties['routes'], 1000):
        name = route.get('name')
        if not name:
            name = '{}-{}'.format(context.env['name'], sha1(json.dumps(route).encode('utf-8')).hexdigest()[:10])
        
        route_properties = {
            'name': name,
            'network': network_name,
            'project': project_id,
            'priority': route.get('priority', i),
        }
        for specified_properties in route:
            route_properties[specified_properties] = route[specified_properties]

        resources.append(
            {
                'name': name,
                'type': 'single_route.py',
                'properties': route_properties
            }
        )

        out[name] = {
            'selfLink': '$(ref.' + name + '.selfLink)',
            'nextHopNetwork': '$(ref.' + name + '.nextHopNetwork)',
        }

    outputs = [{'name': 'routes', 'value': out}]

    return {'resources': resources, 'outputs': outputs}


def generate_network_url(properties):
    """ Gets the network name. """

    network_name = properties.get('network')
    is_self_link = '/' in network_name or '.' in network_name

    if is_self_link:
        network_url = network_name
    else:
        network_url = 'global/networks/{}'.format(network_name)

    return network_url
