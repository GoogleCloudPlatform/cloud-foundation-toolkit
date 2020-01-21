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


def get_certificate(properties, project_id, res_name):
    """
    Gets a link to an existing or newly created SSL Certificate
    resource.
    """

    if 'url' in properties:
        return properties['url'], [], []

    name = '{}-ssl-cert'.format(res_name)

    resource = {
        'name': name,
        'type': 'ssl_certificate.py',
        'properties': copy.copy(properties)
    }
    resource['properties']['name'] = properties.get('name', name)
    resource['properties']['project'] = project_id

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


def get_insecure_proxy(is_http, res_name, project_id, properties, optional_properties):
    """ Creates a TCP or HTTP Proxy resource. """

    if is_http:
        # https://cloud.google.com/compute/docs/reference/rest/v1/targetHttpProxies
        type_name = 'gcp-types/compute-v1:targetHttpProxies'
        target_prop = 'urlMap'
    else:
        # https://cloud.google.com/compute/docs/reference/rest/v1/targetTcpProxies
        type_name = 'gcp-types/compute-v1:targetTcpProxies'
        target_prop = 'service'

    resource_props = {
        'name': properties.get('name', res_name),
        'project': project_id,
    }
    resource = {'type': type_name, 'name': res_name, 'properties': resource_props}

    resource_props[target_prop] = properties['target']

    for prop in optional_properties:
        set_optional_property(resource_props, properties, prop)

    return [resource], []


def get_secure_proxy(is_http, res_name, project_id, properties, optional_properties):
    """ Creates an SSL or HTTPS Proxy resource. """

    if is_http:
        create_base_proxy = get_http_proxy
        # https://cloud.google.com/compute/docs/reference/rest/v1/targetHttpsProxies
        target_type = 'gcp-types/compute-v1:targetHttpsProxies'
    else:
        create_base_proxy = get_tcp_proxy
        # https://cloud.google.com/compute/docs/reference/rest/v1/targetSslProxies
        target_type = 'gcp-types/compute-v1:targetSslProxies'

    # Base proxy settings:
    resources, outputs = create_base_proxy(properties, res_name, project_id)
    resource = resources[0]
    resource['type'] = target_type
    resource_prop = resource['properties']
    for prop in optional_properties:
        set_optional_property(resource_prop, properties, prop)

    # SSL settings:
    ssl_resources = []
    ssl_outputs = []
    if 'certificate' in resource.get('ssl'):
        ssl = properties['ssl']
        url, ssl_resources, ssl_outputs = get_certificate(ssl['certificate'], project_id, res_name)
        resource_prop['sslCertificates'] = [url]
        set_optional_property(resource_prop, ssl, 'sslPolicy')
    if 'sslCertificates' in resource.get('ssl'):
        set_optional_property(resource_prop, ssl, 'sslCertificates')

    return resources + ssl_resources, outputs + ssl_outputs


def get_http_proxy(properties, res_name, project_id):
    """ Creates the HTTP Proxy resource. """

    return get_insecure_proxy(HTTP_BASE, res_name, project_id, properties, ['description'])


def get_tcp_proxy(properties, res_name, project_id):
    """ Creates the TCP Proxy resource. """

    optional_properties = ['description', 'proxyHeader']
    return get_insecure_proxy(TCP_BASE, res_name, project_id, properties, optional_properties)


def get_https_proxy(properties, res_name, project_id):
    """ Creates the HTTPS Proxy resource. """

    return get_secure_proxy(HTTP_BASE, res_name, project_id, properties, ['quicOverride'])


def get_ssl_proxy(properties, res_name, project_id):
    """ Creates the SSL Proxy resource. """

    return get_secure_proxy(TCP_BASE, res_name, project_id, properties, [])


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    name = properties.get('name', context.env['name'])
    project_id = properties.get('project', context.env['project'])
    protocol = properties['protocol']

    if protocol == 'SSL':
        resources, outputs = get_ssl_proxy(properties, context.env['name'], project_id)
    elif protocol == 'TCP':
        resources, outputs = get_tcp_proxy(properties, context.env['name'], project_id)
    elif protocol == 'HTTPS':
        resources, outputs = get_https_proxy(properties, context.env['name'], project_id)
    else:
        resources, outputs = get_http_proxy(properties, context.env['name'], project_id)

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
                    'value': '$(ref.{}.selfLink)'.format(context.env['name'])
                },
                {
                    'name': 'kind',
                    'value': '$(ref.{}.kind)'.format(context.env['name'])
                },
            ]
    }
