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
"""This template creates a custom route."""


def generate_config(context):
    """ Entry point for the deployment resources. """

    network_name = generate_network_url(context.properties)

    resources = []
    out = {}
    for i, route in enumerate(context.properties['routes'], 1000):

        # Set the common route properties.
        properties = {
            'network': network_name,
            'tags': route['tags'],
            'priority': route.get('priority',
                                  i),
            'destRange': route['destRange']
        }

        # Check the route type and fill out the following fields:
        if route['routeType'] == 'ipaddress':
            properties['nextHopIp'] = route.get('nextHopIp')
        elif route['routeType'] == 'instance':
            instance_name = route.get('instanceName')
            zone = route.get('zone', '')
            properties['nextHopInstance'] = generate_instance_url(
                context.env['project'],
                zone,
                instance_name
            )
        elif route['routeType'] == 'gateway':
            gateway_name = route.get('gatewayName')
            properties['nextHopGateway'] = generate_gateway_url(
                context.env['project'],
                gateway_name
            )
        elif route['routeType'] == 'vpntunnel':
            vpn_tunnel_name = route.get('vpnTunnelName')
            region = route.get('region', '')
            properties['nextHopVpnTunnel'] = generate_vpn_tunnel_url(
                context.env['project'],
                region,
                vpn_tunnel_name
            )

        resources.append(
            {
                'name': route['name'],
                'type': 'compute.v1.route',
                'properties': properties
            }
        )

        out[route['name']] = {
            'selfLink': '$(ref.' + route['name'] + '.selfLink)',
            'nextHopNetwork': network_name
        }

    outputs = [{'name': 'routes', 'value': out}]

    return {'resources': resources, 'outputs': outputs}


def generate_network_url(properties):
    """ Gets the network name. """

    network_name = properties.get('network')
    is_self_link = '/' in network_name or '.' in network_name

    if is_self_link:
        network_url = network_name
    else:
        network_url = 'global/networks/{}'.format(network_name)

    return network_url


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
    return 'projects/{}/regions/{}/vpnTunnels/{}'.format(
        project,
        region,
        vpn_tunnel
    )
