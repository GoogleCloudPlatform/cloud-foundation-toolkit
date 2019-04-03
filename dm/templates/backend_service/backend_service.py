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
""" This template creates a backend service. """

REGIONAL_GLOBAL_TYPE_NAMES = {
    True: 'compute.v1.regionBackendService',
    False: 'compute.v1.backendService'
}


def set_optional_property(destination, source, prop_name):
    """ Copies the property value if present. """

    if prop_name in source:
        destination[prop_name] = source[prop_name]


def get_backend_service_outputs(res_name, backend_name, region):
    """ Creates outputs for the backend service. """

    outputs = [
        {
            'name': 'name',
            'value': backend_name
        },
        {
            'name': 'selfLink',
            'value': '$(ref.{}.selfLink)'.format(res_name)
        }
    ]

    if region:
        outputs.append({'name': 'region', 'value': region})

    return outputs


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    res_name = context.env['name']
    name = properties.get('name', res_name)
    is_regional = 'region' in properties
    region = properties.get('region')
    backend_properties = {'name': name}

    resource = {
        'name': res_name,
        'type': REGIONAL_GLOBAL_TYPE_NAMES[is_regional],
        'properties': backend_properties
    }

    optional_properties = [
        'description',
        'backends',
        'timeoutSec',
        'protocol',
        'region',
        'portName',
        'enableCDN',
        'sessionAffinity',
        'affinityCookieTtlSec',
        'loadBalancingScheme',
        'connectionDraining',
        'cdnPolicy'
    ]

    for prop in optional_properties:
        set_optional_property(backend_properties, properties, prop)

    if 'healthCheck' in properties:
        backend_properties['healthChecks'] = [properties['healthCheck']]

    outputs = get_backend_service_outputs(res_name, name, region)

    return {'resources': [resource], 'outputs': outputs}
