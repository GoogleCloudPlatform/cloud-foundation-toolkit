# Copyright 2020 PrimarySite Limited All rights reserved.
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
""" This template creates a backend bucket. """

def set_optional_property(destination, source, prop_name):
    """ Copies the property value if present. """

    if prop_name in source:
        destination[prop_name] = source[prop_name]


def get_backend_bucket_outputs(res_name, backend_name):
    """ Creates outputs for the backend bucket. """

    outputs = [
        {
            'name': 'name',
            'value': backend_name
        },
        {
            'name': 'selfLink',
            'value': '$(ref.{0}.selfLink)'.format(res_name)
        }
    ]

    return outputs


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    res_name = context.env['name']
    name = properties.get('name', res_name)
    project_id = properties.get('project', context.env['project'])
    # bucketName is a required property so we should key error on this if missing
    bucket_name = properties['bucketName']
    backend_properties = {
        'name': name,
        'project': project_id,
        'bucketName': bucket_name,
    }

    resource = {
        'name': res_name,
        'type': 'gcp-types/compute-v1:backendBuckets',
        'properties': backend_properties,
    }

    optional_properties = [
        'description',
        'enableCdn',
        'cdnPolicy'
    ]

    for prop in optional_properties:
        set_optional_property(backend_properties, properties, prop)

    outputs = get_backend_bucket_outputs(res_name, name)

    return {'resources': [resource], 'outputs': outputs}
