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
""" This template creates a forwarding rule. """

REGIONAL_GLOBAL_TYPE_NAMES = {
    True: 'compute.v1.forwardingRule',
    False: 'compute.v1.globalForwardingRule'
}


def set_optional_property(destination, source, prop_name):
    """ Copies the property value, if present. """

    if prop_name in source:
        destination[prop_name] = source[prop_name]


def get_forwarding_rule_outputs(res_name, region):
    """ Creates outputs for the forwarding rule. """

    outputs = [
        {
            'name': 'name',
            'value': '$(ref.{}.name)'.format(res_name)
        },
        {
            'name': 'selfLink',
            'value': '$(ref.{}.selfLink)'.format(res_name)
        },
        {
            'name': 'IPAddress',
            'value': '$(ref.{}.IPAddress)'.format(res_name)
        }
    ]

    if region:
        outputs.append({'name': 'region', 'value': region})

    return outputs


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    name = properties.get('name', context.env['name'])
    is_regional = 'region' in properties
    region = properties.get('region')
    rule_properties = {'name': name}

    resource = {
        'name': name,
        'type': REGIONAL_GLOBAL_TYPE_NAMES[is_regional],
        'properties': rule_properties
    }

    optional_properties = [
        'description',
        'IPAddress',
        'IPProtocol',
        'portRange',
        'ports',
        'region',
        'target',
        'loadBalancingScheme',
        'subnetwork',
        'network',
        'backendService',
        'ipVersion'
    ]

    for prop in optional_properties:
        set_optional_property(rule_properties, properties, prop)

    outputs = get_forwarding_rule_outputs(name, region)

    return {'resources': [resource], 'outputs': outputs}
