# Copyright 2019 Google Inc. All rights reserved.
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
"""This template creates a custom route."""


from hashlib import sha1
import json


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    project_id = properties.get('project', context.env['project'])

    # Set the common route properties.
    res_properties = {
        'network': properties['network'],
        'project': project_id,
        'tags': properties['tags'],
        'priority': properties['priority'],
        'destRange': properties['destRange']
    }

    # Check the route type and fill out the following fields:
    if properties.get('routeType') == 'instance':
        instance_name = properties.get('instanceName')
        zone = properties.get('zone', '')
        res_properties['nextHopInstance'] = generate_instance_url(
            project_id,
            zone,
            instance_name
        )
    elif properties.get('routeType') == 'gateway':
        gateway_name = properties.get('gatewayName')
        res_properties['nextHopGateway'] = generate_gateway_url(
            project_id,
            gateway_name
        )
    elif properties.get('routeType') == 'vpntunnel':
        vpn_tunnel_name = properties.get('vpnTunnelName')
        region = properties.get('region', '')
        res_properties['nextHopVpnTunnel'] = generate_vpn_tunnel_url(
            project_id,
            region,
            vpn_tunnel_name
        )

    optional_properties = [
        'nextHopIp',
        'nextHopInstance',
        'nextHopNetwork',
        'nextHopGateway',
        'nextHopVpnTunnel',
    ]

    for prop in optional_properties:
        if prop in properties:
            res_properties[prop] = properties[prop]

    name = properties['name']
    resources = [
        {
            'name': name,
            # https://cloud.google.com/compute/docs/reference/rest/v1/routes
            'type': 'gcp-types/compute-v1:routes',
            'properties': res_properties
        }
    ]

    outputs = [
        {'name': 'selfLink', 'value': '$(ref.' + name + '.selfLink)'},
        {'name': 'nextHopNetwork', 'value': properties['network']},
    ]

    return {'resources': resources, 'outputs': outputs}


def generate_instance_url(project, zone, instance):
    """ Format the resource name as a resource URI. """

    is_self_link = '/' in instance or '.' in instance

    if is_self_link:
        instance_url = instance
    else:
        instance_url = 'projects/{}/zones/{}/instances/{}'
        instance_url = instance_url.format(project, zone, instance)

    return instance_url


def generate_gateway_url(project, gateway):
    """ Format the resource name as a resource URI. """
    return 'projects/{}/global/gateways/{}'.format(project, gateway)


def generate_vpn_tunnel_url(project, region, vpn_tunnel):
    """ Format the resource name as a resource URI. """
    is_self_link = '/' in vpn_tunnel or '.' in vpn_tunnel

    if is_self_link:
        tunnel_url = vpn_tunnel
    else:
        tunnel_url = 'projects/{}/regions/{}/vpnTunnels/{}'
        tunnel_url = tunnel_url.format(project, region, vpn_tunnel)
    return tunnel_url
