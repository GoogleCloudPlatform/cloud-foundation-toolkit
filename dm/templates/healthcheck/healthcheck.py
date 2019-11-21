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
"""This template creates an instance healthcheck."""


def set_if_exists(healthcheck, properties, prop):
    """
    If prop exists in properties, set the healthcheck's property to it.
    Input:  [dict] healthcheck: a dictionary representing a healthcheck object
            [dict] properties: a dictionary of the user supplied values
            [string] prop: the value to check if exists within properties

    """
    if prop in properties:
        healthcheck[prop] = properties[prop]

def generate_config(context):
    """ Entry point for the deployment resources. """

    resources = []
    properties = context.properties
    healthcheck_name = properties.get('name', context.env['name'])
    healthcheck_type = properties['healthcheckType']
    healthcheck_version = properties.get('version', 'v1')

    project_id = properties.get('project', context.env['project'])

    # Deployment Manager resource types per healthcheck type.
    healthcheck_type_dictionary = {
        'HTTP':
            {
                'v1': 'gcp-types/compute-v1:httpHealthChecks',
                'beta': 'gcp-types/compute-beta:httpHealthChecks'
            },
        'HTTPS':
            {
                'v1': 'gcp-types/compute-v1:httpsHealthChecks',
                'beta': 'gcp-types/compute-beta:httpsHealthChecks'
            },
        'SSL':
            {
                'v1': 'gcp-types/compute-v1:healthChecks',
                'beta': 'gcp-types/compute-beta:healthChecks'
            },
        'TCP':
            {
                'v1': 'gcp-types/compute-v1:healthChecks',
                'beta': 'gcp-types/compute-beta:healthChecks'
            },
        'HTTP2': {
            'beta': 'gcp-types/compute-beta:healthChecks'
        }
    }

    # Deployment Manager object types associated with each type of healthcheck.
    healthcheck_object_dictionary = {
        'HTTP': 'httpHealthCheck',
        'HTTPS': 'httpsHealthCheck',
        'SSL': 'sslHealthCheck',
        'TCP': 'tcpHealthCheck',
        'HTTP2': 'http2HealthCheck'
    }

    # Create a generic healthcheck object.
    healthcheck = {
        'name':
            context.env['name'],
        'type':
            healthcheck_type_dictionary[healthcheck_type][healthcheck_version]
    }

    # Create the generic healthcheck properties separately.
    healthcheck_properties = {
        'checkIntervalSec': properties['checkIntervalSec'],
        'timeoutSec': properties['timeoutSec'],
        'unhealthyThreshold': properties['unhealthyThreshold'],
        'healthyThreshold': properties['healthyThreshold'],
        'kind': 'compute#healthCheck',
        'type': healthcheck_type,
        'project': project_id,
        'name': healthcheck_name,
    }

    set_if_exists(healthcheck_properties, properties, 'description')

    # Create a specific healthcheck object.
    specific_healthcheck_type = healthcheck_object_dictionary[healthcheck_type]
    specific_healthcheck = {
        'proxyHeader': properties.get('proxyHeader',
                                      'NONE'),
    }

    set_if_exists(specific_healthcheck, properties, 'port')

    # Check for beta-specific properties.
    # Add them to the specific healthcheck object.
    if healthcheck_version == 'beta':
        for prop in ['portSpecification', 'portName', 'response']:
            set_if_exists(specific_healthcheck, properties, prop)

    # Check for HTTP/S/2-specific properties.
    # Add them to the generic healthcheck.
    if healthcheck_type in ['HTTP', 'HTTPS', 'HTTP2']:
        for prop in ['requestPath', 'host']:
            set_if_exists(healthcheck_properties, properties, prop)

    # Check for TCP/SSL-specific properties.
    # Add them to the specific healthcheck object.
    if healthcheck_type in ['TCP', 'SSL']:
        for prop in ['request', 'response']:
            set_if_exists(specific_healthcheck, properties, prop)

    healthcheck_properties[specific_healthcheck_type] = specific_healthcheck
    healthcheck['properties'] = healthcheck_properties
    resources.append(healthcheck)

    outputs = [
        {
            'name': 'name',
            'value': '$(ref.{}.name)'.format(context.env['name'])
        },
        {
            'name': 'selfLink',
            'value': '$(ref.{}.selfLink)'.format(context.env['name'])
        },
        {
            'name': 'creationTimestamp',
            'value': '$(ref.{}.creationTimestamp)'.format(context.env['name'])
        }
    ]

    return {'resources': resources, 'outputs': outputs}
