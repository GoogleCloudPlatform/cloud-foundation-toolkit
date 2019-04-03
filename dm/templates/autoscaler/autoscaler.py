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
""" This template creates an autoscaler. """

REGIONAL_LOCAL_AUTOSCALER_TYPES = {
    True: 'compute.v1.regionAutoscaler',
    False: 'compute.v1.autoscaler'
}

def set_optional_property(receiver, source, property_name):
    """ If set, copies the given property value from one object to another. """

    if property_name in source:
        receiver[property_name] = source[property_name]

def set_autoscaler_location(autoscaler, is_regional, location):
    """ Sets location-dependent properties of the autoscaler. """

    name = autoscaler['name']
    location_prop_name = 'region' if is_regional else 'zone'

    autoscaler['type'] = REGIONAL_LOCAL_AUTOSCALER_TYPES[is_regional]
    autoscaler['properties'][location_prop_name] = location
    location_output = {
        'name': location_prop_name,
        'value': '$(ref.{}.{})'.format(name, location_prop_name)
    }

    return location_output

def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    name = properties.get('name', context.env['name'])
    target = properties['target']

    policy = {}

    autoscaler = {
        'type': None, # Will be set up at a later stage.
        'name': name,
        'properties': {
            'autoscalingPolicy': policy,
            'target': target
        }
    }

    policy_props = ['coolDownPeriodSec',
                    'minNumReplicas',
                    'maxNumReplicas',
                    'customMetricUtilizations',
                    'loadBalancingUtilization',
                    'cpuUtilization']

    for prop in policy_props:
        set_optional_property(policy, properties, prop)

    is_regional = 'region' in properties
    location = properties['region'] if is_regional else properties['zone']
    location_output = set_autoscaler_location(autoscaler, is_regional, location)

    set_optional_property(autoscaler['properties'], properties, 'description')

    return {
        'resources': [autoscaler],
        'outputs': [
            {
                'name': 'name',
                'value': name
            },
            {
                'name': 'selfLink',
                'value': '$(ref.{}.selfLink)'.format(name)
            }
        ] + [location_output]
    }
