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
""" This template creates a Cloud Run """

import collections
import random
import string

DMBundle = collections.namedtuple('DMBundle', 'resource outputs')

SUFFIX_LENGTH = 5
CHAR_CHOICE = string.digits + string.ascii_lowercase


def get_random_string(length):
    """ Creates a random string of characters of the specified length. """

    return ''.join([random.choice(CHAR_CHOICE) for _ in range(length)])


def generate_config(context):
    """ Creates the Cloud SQL instance, databases, and user. """

    properties = context.properties
    res_name = properties.get('name', context.env['name'])
    project_id = properties.get('project', context.env['project'])

    cloud_run = {
        'name': res_name,
        'type': project_id + '/cloud-run-custom-type:namespaces.services',
        'properties': {
            'oauth_token': properties.get('oauth_token'),
            'parent': 'namespaces/' + project_id,
            'kind': 'Service',
            'apiVersion': 'serving.knative.dev/v1',
            'metadata': {
                'name': properties.get('metadataName') + '-' + get_random_string(6),
            },
            'spec': {
                'template': {
                    'spec': {
                        'containers': properties.get('containers')
                    }
                }
            }
        }
    }

    outputs = [
        {
            'name': 'selfLink',
            'value': '$(ref.{}.metadata.selfLink)'.format(res_name)
        }
    ]

    return {'resources': [cloud_run], 'outputs': outputs}
