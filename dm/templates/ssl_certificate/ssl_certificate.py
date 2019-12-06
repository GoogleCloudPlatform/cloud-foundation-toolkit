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
""" This template creates an SSL certificate. """


def set_optional_property(destination, source, prop_name):
    """ Copies the property value if present. """

    if prop_name in source:
        destination[prop_name] = source[prop_name]


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    name = properties.get('name', context.env['name'])
    project_id = properties.get('project', context.env['project'])

    ssl_props = {
        'name': name,
        'project': project_id,
    }
 
    if properties.get('betaFeaturesEnabled', False):
        gcptype = 'gcp-types/compute-beta:sslCertificates'
    else:
        gcptype = 'gcp-types/compute-v1:sslCertificates'

    resource = {
        'name': context.env['name'],
        # https://cloud.google.com/compute/docs/reference/rest/v1/sslCertificates
        'type': gcptype,
        'properties': ssl_props,
    }

    for prop in [
            'privateKey',
            'certificate',
            'description',
            'managed',
            'selfManaged',
            'type']:
        set_optional_property(ssl_props, properties, prop)

    return {
        'resources': [resource],
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
