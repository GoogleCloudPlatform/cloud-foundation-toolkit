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

""" This template creates an Interconnect resource. """


def generate_config(context):
    """ Entry point for the deployment resources. """
    resources = []
    intercon = {
        'name': context.env['name'],
        'type': 'compute.v1.interconnects',
        'properties':
            {
                'name':
                    context.properties.get('name', context.env['name']),
                'customerName':
                    context.properties['customerName'],
                'interconnectType':
                    context.properties['interconnectType'],
                'location':
                    context.properties['location'],
                'requestedLinkCount':
                    context.properties['requestedLinkCount']
            }
    }

    optional_props = [
        'adminEnabled',
        'description',
        'linkType',
        'nocContactEmail'
    ]

    for prop in optional_props:
        if prop in context.properties:
            intercon['properties'][prop] = context.properties[prop]

    resources.append(intercon)

    return {
        'resources':
            resources,
        'outputs':
            [
                {
                    'name': 'name',
                    'value': context.env['name']
                },
                {
                    'name': 'selfLink',
                    'value': '$(ref.{}.selfLink)'.format(context.env['name'])
                }
            ]
    }
