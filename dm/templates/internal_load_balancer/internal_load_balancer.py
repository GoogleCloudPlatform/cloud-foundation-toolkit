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
""" This template creates an internal load balancer. """


def set_optional_property(destination, source, prop_name):
    """ Copies the property value if present. """

    if prop_name in source:
        destination[prop_name] = source[prop_name]


def get_backend_service(properties, project_id, res_name):
    """ Creates the backend service. """

    backend_spec = properties['backendService']
    for backend in backend_spec['backends']:
        backend.update({
            'balancingMode': 'CONNECTION'
        })

    name = '{}-bs'.format(res_name)
    backend_properties = {
        'name': backend_spec.get('name', properties.get('name', name)),
        'project': project_id,
        'loadBalancingScheme': 'INTERNAL',
        'protocol': properties['protocol'],
        'region': properties['region'],
    }

    backend_resource = {
        'name': name,
        'type': 'backend_service.py',
        'properties': backend_properties
    }

    optional_properties = [
        'description',
        'backends',
        'timeoutSec',
        'sessionAffinity',
        'connectionDraining',
        'backends',
        'healthCheck',
        'healthChecks',
    ]

    for prop in optional_properties:
        set_optional_property(backend_properties, backend_spec, prop)

    return [backend_resource], [
        {
            'name': 'backendServiceName',
            'value': backend_resource['properties']['name'],
        },
        {
            'name': 'backendServiceSelfLink',
            'value': '$(ref.{}.selfLink)'.format(name),
        },
    ]


def get_forwarding_rule(properties, backend, project_id, res_name):
    """ Creates the forwarding rule. """

    rule_properties = {
        'name': properties.get('name', res_name),
        'project': project_id,
        'loadBalancingScheme': 'INTERNAL',
        'IPProtocol': properties['protocol'],
        'backendService': '$(ref.{}.selfLink)'.format(backend['name']),
        'region': properties['region'],
    }

    rule_resource = {
        'name': res_name,
        'type': 'forwarding_rule.py',
        'properties': rule_properties,
    }

    optional_properties = [
        'description',
        'IPAddress',
        'ipVersion',
        'ports',
        'network',
        'subnetwork',
    ]

    for prop in optional_properties:
        set_optional_property(rule_properties, properties, prop)

    return [rule_resource], [
        {
            'name': 'forwardingRuleName',
            'value': res_name,
        },
        {
            'name': 'forwardingRuleSelfLink',
            'value': '$(ref.{}.selfLink)'.format(res_name),
        },
        {
            'name': 'IPAddress',
            'value': '$(ref.{}.IPAddress)'.format(res_name),
        },
        {
            'name': 'region',
            'value': properties['region']
        },
    ]


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    project_id = properties.get('project', context.env['project'])

    backend_resources, backend_outputs = get_backend_service(properties, project_id, context.env['name'])
    rule_resources, rule_outputs = get_forwarding_rule(
        properties,
        backend_resources[0],
        project_id,
        context.env['name']
    )

    return {
        'resources': rule_resources + backend_resources,
        'outputs': rule_outputs + backend_outputs
    }
