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
""" This template creates a Stackdriver Metric Descriptor. """


def generate_config(context):
    """ Entry point for the deployment resources. """

    resources = []
    outputs = []
    properties = context.properties
    name = properties.get('name', context.env['name'])
    metric_descriptor = {
        'name': name,
        'type': 'gcp-types/monitoring-v3:projects.metricDescriptors',
        'properties': {}
    }

    required_properties = [
        'type',
        'metricKind',
        'valueType',
        'unit'
    ]

    for prop in required_properties:
        if prop in properties:
            metric_descriptor['properties'][prop] = properties[prop]

    # Optional properties:
    optional_properties = ['displayName', 'labels', 'description', 'metadata']

    for prop in optional_properties:
        if prop in properties:
            metric_descriptor['properties'][prop] = properties[prop]

    resources.append(metric_descriptor)

    # Output variables:
    output_props = [
        'name',
        'type',
        'labels',
        'metricKind',
        'valueType',
        'unit',
        'description',
        'displayName',
        'metadata'
    ]

    for outprop in output_props:
        output = {}
        if outprop in properties:
            output['name'] = outprop
            output['value'] = '$(ref.{}.{})'.format(name, outprop)
            outputs.append(output)

    return {'resources': resources, 'outputs': outputs}
