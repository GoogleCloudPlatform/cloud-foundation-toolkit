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


def set_optional_property(destination, source, prop_name):
    """ Copies the property value if present. """

    if prop_name in source:
        destination[prop_name] = source[prop_name]


def get_backend_service(properties, backend_spec):
    """ Creates the backend service. """

    backend_properties = {
        'loadBalancingScheme': 'EXTERNAL',
        'protocol': get_protocol(properties)
    }

    backend_name = backend_spec['name']
    backend_resource = {
        'name': backend_name,
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
            'value': '$(ref.{}.selfLink)'.format(backend_name),
        },
    ]


def get_forwarding_rule(properties, target, name):
    """ Creates the forwarding rule. """

    rule_properties = {
        'loadBalancingScheme': 'EXTERNAL',
        'target': '$(ref.{}.selfLink)'.format(target['name']),
        'IPProtocol': 'TCP'
    }

    rule_resource = {
        'name': name,
        'type': 'forwarding_rule.py',
        'properties': rule_properties
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
            'value': name,
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


def get_backend_services(properties):
    """ Creates all backend services to be used by the load balancer. """

    backend_resources = []
    backend_outputs_map = {
        'backendServiceName': [],
        'backendServiceSelfLink': []
    }
    backend_specs = properties['backendServices']

    for backend_spec in backend_specs:
        resources, outputs = get_backend_service(properties, backend_spec)
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


def get_url_map(properties, res_name):
    """ Creates a UrlMap resource. """

    name = properties.get('name', res_name + '-url-map')
    spec = copy.deepcopy(properties)
    update_refs_recursively(spec)

    resource = {'name': name, 'type': 'url_map.py', 'properties': spec}

    self_link = '$(ref.{}.selfLink)'.format(name)

    return self_link, [resource], [
        {
            'name': 'urlMapName',
            'value': '$(ref.{}.name)'.format(name)
        },
        {
            'name': 'urlMapSelfLink',
            'value': self_link
        }
    ]


def get_target_proxy(properties, res_name):
    """ Creates a target proxy resource. """

    protocol = get_protocol(properties)

    if 'HTTP' in protocol:
        target, resources, outputs = get_url_map(
            properties['urlMap'],
            res_name
        )
    else:
        backend_services = properties['backendServices']
        service_name = backend_services[0]['name']
        target = get_ref(service_name)
        resources = []
        outputs = []

    name = '{}-target'.format(res_name)
    proxy = {
        'name': name,
        'type': 'target_proxy.py',
        'properties': {
            'protocol': protocol,
            'target': target
        }
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
    name = properties.get('name', context.env['name'])

    # Forwarding rule + target proxy + backend service = ELB
    bs_resources, bs_outputs = get_backend_services(properties)
    target_resources, target_outputs = get_target_proxy(properties, name)
    rule_resources, rule_outputs = get_forwarding_rule(
        properties,
        target_resources[0],
        name
    )

    return {
        'resources': bs_resources + target_resources + rule_resources,
        'outputs': bs_outputs + target_outputs + rule_outputs
    }
