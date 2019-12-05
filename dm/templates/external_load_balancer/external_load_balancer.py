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
""" This template creates an external load balancer. """

import copy
from hashlib import sha1
import json


def set_optional_property(destination, source, prop_name):
    """ Copies the property value if present. """

    if prop_name in source:
        destination[prop_name] = source[prop_name]


def get_backend_service(properties, backend_spec, res_name, project_id):
    """ Creates the backend service. """

    name = backend_spec.get('resourceName', res_name)
    backend_name = backend_spec.get('name', name)
    backend_properties = {
        'name': backend_name,
        'project': project_id,
        'loadBalancingScheme': 'EXTERNAL',
        'protocol': get_protocol(properties),
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
        'portName',
        'enableCDN',
        'affinityCookieTtlSec'
    ]

    for prop in optional_properties:
        set_optional_property(backend_properties, backend_spec, prop)

    return [backend_resource], [
        {
            'name': 'backendServiceName',
            'value': backend_name,
        },
        {
            'name': 'backendServiceSelfLink',
            'value': '$(ref.{}.selfLink)'.format(name),
        },
    ]


def get_forwarding_rule(properties, target, res_name, project_id):
    """ Creates the forwarding rule. """

    name = '{}-forwarding-rule'.format(res_name)
    rule_properties = {
        'name': properties.get('name', res_name),
        'project': project_id,
        'loadBalancingScheme': 'EXTERNAL',
        'target': '$(ref.{}.selfLink)'.format(target['name']),
        'IPProtocol': 'TCP',
    }

    rule_resource = {
        'name': name,
        'type': 'forwarding_rule.py',
        'properties': rule_properties,
        'metadata': {
            'dependsOn': [target['name']],
        },
    }

    optional_properties = [
        'description',
        'IPAddress',
        'ipVersion',
        'portRange',
    ]

    for prop in optional_properties:
        set_optional_property(rule_properties, properties, prop)

    return [rule_resource], [
        {
            'name': 'forwardingRuleName',
            'value': rule_properties['name'],
        },
        {
            'name': 'forwardingRuleSelfLink',
            'value': '$(ref.{}.selfLink)'.format(name),
        },
        {
            'name': 'IPAddress',
            'value': '$(ref.{}.IPAddress)'.format(name),
        },
    ]


def get_backend_services(properties, res_name, project_id):
    """ Creates all backend services to be used by the load balancer. """

    backend_resources = []
    backend_outputs_map = {
        'backendServiceName': [],
        'backendServiceSelfLink': []
    }
    backend_specs = properties['backendServices']

    for backend_spec in backend_specs:
        backend_res_name = '{}-backend-service-{}'.format(res_name, sha1(json.dumps(backend_spec).encode('utf-8')).hexdigest()[:10])
        resources, outputs = get_backend_service(properties, backend_spec, backend_res_name, project_id)
        backend_resources += resources
        # Merge outputs with the same name.
        for output in outputs:
            backend_outputs_map[output['name']].append(output['value'])

    backend_outputs = []
    for key, value in backend_outputs_map.items():
        backend_outputs.append({'name': key + 's', 'value': value})

    return backend_resources, backend_outputs


def get_ref(name, prop='selfLink'):
    """ Creates reference to a property of a given resource. """

    return '$(ref.{}.{})'.format(name, prop)


def update_refs_recursively(properties):
    """ Replaces service names with the service selflinks recursively. """

    for prop in properties:
        value = properties[prop]
        if prop == 'defaultService' or prop == 'service':
            is_regular_name = not '.' in value and not '/' in value
            if is_regular_name:
                properties[prop] = get_ref(value)
        elif isinstance(value, dict):
            update_refs_recursively(value)
        elif isinstance(value, list):
            for item in value:
                if isinstance(item, dict):
                    update_refs_recursively(item)


def get_url_map(properties, res_name, project_id):
    """ Creates a UrlMap resource. """

    spec = copy.deepcopy(properties)
    spec['project'] = project_id
    spec['name'] = properties.get('name', res_name)
    update_refs_recursively(spec)

    resource = {
        'name': res_name,
        'type': 'url_map.py',
        'properties': spec,
    }

    self_link = '$(ref.{}.selfLink)'.format(res_name)

    return self_link, [resource], [
        {
            'name': 'urlMapName',
            'value': '$(ref.{}.name)'.format(res_name)
        },
        {
            'name': 'urlMapSelfLink',
            'value': self_link
        }
    ]


def get_target_proxy(properties, res_name, project_id, bs_resources):
    """ Creates a target proxy resource. """

    protocol = get_protocol(properties)

    depends = []
    if 'HTTP' in protocol:
        urlMap = copy.deepcopy(properties['urlMap'])
        if 'name' not in urlMap and 'name' in properties:
            urlMap['name'] = '{}-url-map'.format(properties['name'])
        target, resources, outputs = get_url_map(
            urlMap,
            '{}-url-map'.format(res_name),
            project_id
        )
        depends.append(resources[0]['name'])
    else:
        depends.append(bs_resources[0]['name'])
        target = get_ref(bs_resources[0]['name'])
        resources = []
        outputs = []

    name = '{}-target'.format(res_name)
    proxy = {
        'name': name,
        'type': 'target_proxy.py',
        'properties': {
            'name': '{}-target'.format(properties.get('name', res_name)),
            'project': project_id,
            'protocol': protocol,
            'target': target,
        },
        'metadata': {
            'dependsOn': [depends],
        },
    }

    for prop in ['proxyHeader', 'quicOverride']:
        set_optional_property(proxy['properties'], properties, prop)

    outputs.extend(
        [
            {
                'name': 'targetProxyName',
                'value': '$(ref.{}.name)'.format(name)
            },
            {
                'name': 'targetProxySelfLink',
                'value': '$(ref.{}.selfLink)'.format(name)
            },
            {
                'name': 'targetProxyKind',
                'value': '$(ref.{}.kind)'.format(name)
            }
        ]
    )

    if 'ssl' in properties:
        ssl_spec = properties['ssl']
        proxy['properties']['ssl'] = ssl_spec
        creates_new_certificate = not 'url' in ssl_spec['certificate']
        if creates_new_certificate:
            outputs.extend(
                [
                    {
                        'name': 'certificateName',
                        'value': '$(ref.{}.certificateName)'.format(name)
                    },
                    {
                        'name': 'certificateSelfLink',
                        'value': '$(ref.{}.certificateSelfLink)'.format(name)
                    }
                ]
            )

    return [proxy] + resources, outputs


def get_protocol(properties):
    """ Finds what network protocol to use. """

    is_web = 'urlMap' in properties
    is_secure = 'ssl' in properties

    if is_web:
        if is_secure:
            return 'HTTPS'
        return 'HTTP'

    if is_secure:
        return 'SSL'
    return 'TCP'


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    project_id = properties.get('project', context.env['project'])

    # Forwarding rule + target proxy + backend service = ELB
    bs_resources, bs_outputs = get_backend_services(properties, context.env['name'], project_id)
    target_resources, target_outputs = get_target_proxy(properties, context.env['name'], project_id, bs_resources)
    rule_resources, rule_outputs = get_forwarding_rule(
        properties,
        target_resources[0],
        context.env['name'],
        project_id
    )

    return {
        'resources': bs_resources + target_resources + rule_resources,
        'outputs': bs_outputs + target_outputs + rule_outputs,
    }
