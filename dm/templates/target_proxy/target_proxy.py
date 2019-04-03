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
""" This template creates HTTP(S) and TCP/SSL proxy resources. """

import copy

HTTP_BASE = True
TCP_BASE = False


def set_optional_property(destination, source, prop_name):
    """ Copies the property value if present. """

    if prop_name in source:
        destination[prop_name] = source[prop_name]


def get_certificate(properties, res_name):
    """
    Gets a link to an existing or newly created SSL Certificate
    resource.
    """

    if 'url' in properties:
        return properties['url'], [], []

    name = properties.get('name', '{}-ssl-cert'.format(res_name))

    resource = {
        'name': name,
        'type': 'ssl_certificate.py',
        'properties': copy.copy(properties)
    }

    self_link = '$(ref.{}.selfLink)'.format(name)
    outputs = [
        {
            'name': 'certificateName',
            'value': '$(ref.{}.name)'.format(name)
        },
        {
            'name': 'certificateSelfLink',
            'value': self_link
        }
    ]

    return self_link, [resource], outputs


def get_insecure_proxy(is_http, name, properties, optional_properties):
    """ Creates a TCP or HTTP Proxy resource. """

    if is_http:
        type_name = 'compute.v1.targetHttpProxy'
        target_prop = 'urlMap'
    else:
        type_name = 'compute.alpha.targetTcpProxy'
        target_prop = 'service'

    resource_props = {}
    resource = {'type': type_name, 'name': name, 'properties': resource_props}

    resource_props[target_prop] = properties['target']

    for prop in optional_properties:
        set_optional_property(resource_props, properties, prop)

    return [resource], []


def get_secure_proxy(is_http, name, properties, optional_properties):
    """ Creates an SSL or HTTPS Proxy resource. """

    if is_http:
        create_base_proxy = get_http_proxy
        target_type = 'compute.v1.targetHttpsProxy'
    else:
        create_base_proxy = get_tcp_proxy
        target_type = 'compute.v1.targetSslProxy'

    # Base proxy settings:
    resources, outputs = create_base_proxy(properties, name)
    resource = resources[0]
    resource['type'] = target_type
    resource_prop = resource['properties']
    for prop in optional_properties:
        set_optional_property(resource_prop, properties, prop)

    # SSL settings:
    ssl = properties['ssl']
    url, ssl_resources, ssl_outputs = get_certificate(ssl['certificate'], name)
    resource_prop['sslCertificates'] = [url]
    set_optional_property(resource_prop, ssl, 'sslPolicy')

    return resources + ssl_resources, outputs + ssl_outputs


def get_http_proxy(properties, name):
    """ Creates the HTTP Proxy resource. """

    return get_insecure_proxy(HTTP_BASE, name, properties, ['description'])


def get_tcp_proxy(properties, name):
    """ Creates the TCP Proxy resource. """

    optional_properties = ['description', 'proxyHeader']
    return get_insecure_proxy(TCP_BASE, name, properties, optional_properties)


def get_https_proxy(properties, name):
    """ Creates the HTTPS Proxy resource. """

    return get_secure_proxy(HTTP_BASE, name, properties, ['quicOverride'])


def get_ssl_proxy(properties, name):
    """ Creates the SSL Proxy resource. """

    return get_secure_proxy(TCP_BASE, name, properties, [])


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    name = properties.get('name', context.env['name'])
    protocol = properties['protocol']

    if protocol == 'SSL':
        resources, outputs = get_ssl_proxy(properties, name)
    elif protocol == 'TCP':
        resources, outputs = get_tcp_proxy(properties, name)
    elif protocol == 'HTTPS':
        resources, outputs = get_https_proxy(properties, name)
    else:
        resources, outputs = get_http_proxy(properties, name)

    return {
        'resources':
            resources,
        'outputs':
            outputs + [
                {
                    'name': 'name',
                    'value': name
                },
                {
                    'name': 'selfLink',
                    'value': '$(ref.{}.selfLink)'.format(name)
                },
                {
                    'name': 'kind',
                    'value': '$(ref.{}.kind)'.format(name)
                },
            ]
    }
