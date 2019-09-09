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

""" This template creates an Interconnect Attachment. """


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    name = properties.get('name', context.env['name'])
    project_id = properties.get('project', context.env['project'])

    resources = []
    attach = {
        'name': context.env['name'],
        # https://cloud.google.com/compute/docs/reference/rest/v1/interconnectAttachments
        'type': 'gcp-types/compute-v1:interconnectAttachments',
        'properties':
            {
                'project': project_id,
                'name': name,
                'router':
                    context.properties['router'],
                'region':
                    context.properties['region'],
                'type':
                    context.properties['type']
            }
    }

    optional_props = [
        'adminEnabled',
        'bandwidth',
        'candidateSubnets',
        'description',
        'edgeAvailabilityDomain',
        'interconnect',
        'partnerAsn',
        'partnerMetadata',
        'vlanTag8021q',
    ]

    for prop in optional_props:
        if prop in context.properties:
            attach['properties'][prop] = context.properties[prop]

    resources.append(attach)

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
                    'value': '$(ref.{}.selfLink)'.format(context.env['name'])
                }
            ]
    }
