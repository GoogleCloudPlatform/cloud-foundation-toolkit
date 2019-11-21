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
""" This template creates a network, optionally with subnetworks. """


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
    network_self_link = '$(ref.{}.selfLink)'.format(context.env['name'])

    network_resource = {
        # https://cloud.google.com/compute/docs/reference/rest/v1/networks/insert
        'type': 'gcp-types/compute-v1:networks',
        'name': context.env['name'],
        'properties':
            {
                'name': name,
                'autoCreateSubnetworks': properties.get('autoCreateSubnetworks', False)
            }
    }
    optional_properties = [
        'description',
        'routingConfig',
        'project',
    ]
    for prop in optional_properties:
        append_optional_property(network_resource, properties, prop)
    resources = [network_resource]

    # Subnetworks:
    out = {}
    for i, subnetwork in enumerate(
        properties.get('subnetworks', []), 1
    ):
        subnetwork['network'] = network_self_link
        if properties.get('project'):
            subnetwork['project'] = properties.get('project')

        subnetwork_name = 'subnetwork-{}'.format(i)
        resources.append(
            {
                'name': subnetwork_name,
                'type': 'subnetwork.py',
                'properties': subnetwork
            }
        )

        out[subnetwork_name] = {
            'selfLink': '$(ref.{}.selfLink)'.format(subnetwork_name),
            'ipCidrRange': '$(ref.{}.ipCidrRange)'.format(subnetwork_name),
            'region': '$(ref.{}.region)'.format(subnetwork_name),
            'network': '$(ref.{}.network)'.format(subnetwork_name),
            'gatewayAddress': '$(ref.{}.gatewayAddress)'.format(subnetwork_name)
        }

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
                    'value': network_self_link
                },
                {
                    'name': 'subnetworks',
                    'value': out
                }
            ]
    }
