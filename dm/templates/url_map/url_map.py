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
""" This template creates a URL map. """


def set_optional_property(destination, source, prop_name):
    """ Copies the property value if present. """

    if prop_name in source:
        destination[prop_name] = source[prop_name]


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    name = properties.get('name', context.env['name'])
    project_id = properties.get('project', context.env['project'])

    resource = {
        'name': context.env['name'],
        # https://cloud.google.com/compute/docs/reference/rest/v1/urlMaps
        'type': 'gcp-types/compute-v1:urlMaps',
        'properties': {
            'name': name,
            'project': project_id,
        },
    }

    optional_properties = [
        'defaultService',
        'defaultUrlRedirect',
        'description',
        'hostRules',
        'pathMatchers',
        'tests',
    ]

    for prop in optional_properties:
        set_optional_property(resource['properties'], properties, prop)

    return {
        'resources': [resource],
        'output':
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
