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


def generate_config(context):
    """ Entry point for the deployment resources. """

    name = context.properties.get('name') or context.env['name']
    network_self_link = '$(ref.{}.selfLink)'.format(name)
    auto_create_subnetworks = context.properties.get(
        'autoCreateSubnetworks',
        False
    )

    resources = [
        {
            'type': 'compute.v1.network',
            'name': name,
            'properties':
                {
                    'name': name,
                    'autoCreateSubnetworks': auto_create_subnetworks
                }
        }
    ]

    # Subnetworks:
    out = {}
    for subnetwork in context.properties.get('subnetworks', []):
        subnetwork['network'] = network_self_link
        resources.append(
            {
                'name': subnetwork['name'],
                'type': 'subnetwork.py',
                'properties': subnetwork
            }
        )

        out[subnetwork['name']] = {
            'selfLink': '$(ref.{}.selfLink)'.format(subnetwork['name']),
            'ipCidrRange': '$(ref.{}.ipCidrRange)'.format(subnetwork['name']),
            'region': '$(ref.{}.region)'.format(subnetwork['name']),
            'network': '$(ref.{}.network)'.format(subnetwork['name']),
            'gatewayAddress': '$(ref.{}.gatewayAddress)'.format(
                subnetwork['name']
            )
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
