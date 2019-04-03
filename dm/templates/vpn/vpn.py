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
"""This template creates a VPN tunnel, gateway, and forwarding rules."""


def generate_config(context):
    """ Entry point for the deployment resources. """

    network = generate_network_url(
        context.env['project'],
        context.properties['network']
    )
    target_vpn_gateway = context.env['name'] + '-tvpng'
    static_ip = context.env['name'] + '-ip'
    esp_rule = context.env['name'] + '-esp-rule'
    udp_500_rule = context.env['name'] + '-udp-500-rule'
    udp_4500_rule = context.env['name'] + '-udp-4500-rule'
    vpn_tunnel = context.env['name'] + '-vpn'
    router_vpn_binding = context.env['name'] + '-router-vpn-binding'

    resources = [
        {
            # The target VPN gateway resource.
            'name': target_vpn_gateway,
            'type': 'compute.v1.targetVpnGateway',
            'properties':
                {
                    'network': network,
                    'region': context.properties['region']
                }
        },
        {
            # The reserved address resource.
            'name': static_ip,
            'type': 'compute.v1.address',
            'properties': {
                'region': context.properties['region']
            }
        },
        {
            # The forwarding rule resource for the ESP traffic.
            'name': esp_rule,
            'type': 'compute.v1.forwardingRule',
            'properties':
                {
                    'IPAddress': '$(ref.' + static_ip + '.address)',
                    'IPProtocol': 'ESP',
                    'region': context.properties['region'],
                    'target': '$(ref.' + target_vpn_gateway + '.selfLink)'
                }
        },
        {
            # The forwarding rule resource for the UDP traffic on port 4500.
            'name': udp_4500_rule,
            'type': 'compute.v1.forwardingRule',
            'properties':
                {
                    'IPAddress': '$(ref.' + static_ip + '.address)',
                    'IPProtocol': 'UDP',
                    'portRange': 4500,
                    'region': context.properties['region'],
                    'target': '$(ref.' + target_vpn_gateway + '.selfLink)'
                }
        },
        {
            # The forwarding rule resource for the UDP traffic on port 500
            'name': udp_500_rule,
            'type': 'compute.v1.forwardingRule',
            'properties':
                {
                    'IPAddress': '$(ref.' + static_ip + '.address)',
                    'IPProtocol': 'UDP',
                    'portRange': 500,
                    'region': context.properties['region'],
                    'target': '$(ref.' + target_vpn_gateway + '.selfLink)'
                }
        },
        {
            # The VPN tunnel resource.
            'name': vpn_tunnel,
            'type': 'compute.v1.vpnTunnel',
            'properties':
                {
                    'description':
                        'A vpn tunnel',
                    'ikeVersion':
                        2,
                    'peerIp':
                        context.properties['peerAddress'],
                    'region':
                        context.properties['region'],
                    'router':
                        generate_router_url(
                            context.env['project'],
                            context.properties['region'],
                            context.properties['router']
                        ),
                    'sharedSecret':
                        context.properties['sharedSecret'],
                    'targetVpnGateway':
                        '$(ref.' + target_vpn_gateway + '.selfLink)'
                },
            'metadata': {
                'dependsOn': [esp_rule,
                              udp_500_rule,
                              udp_4500_rule]
            }
        },
        {
            # An action that is executed after the vpn_tunnel function.
            # It calls the method patch by ID on the descriptor document
            # https://www.googleapis.com/discovery/v1/apis/compute/v1/rest.
            'name': router_vpn_binding,
            'action': 'gcp-types/compute-v1:compute.routers.patch',
            'properties':
                {
                    'router':
                        context.properties['router'],
                    'region':
                        context.properties['region'],
                    'project':
                        context.env['project'],
                    'name':
                        context.properties['router'],
                    'asn':
                        context.properties['asn'],
                    'interfaces':
                        [
                            {
                                'ipRange':
                                    '169.254.1.1/31',
                                'linkedVpnTunnel':
                                    '$(ref.' + vpn_tunnel + '.selfLink)',
                                'name':
                                    'if-1'
                            }
                        ]
                }
        }
    ]

    return {
        'resources':
            resources,
        'outputs':
            [
                {
                    'name': 'targetVpnGateway',
                    'value': target_vpn_gateway
                },
                {
                    'name': 'staticIp',
                    'value': static_ip
                },
                {
                    'name': 'espRule',
                    'value': esp_rule
                },
                {
                    'name': 'udp500Rule',
                    'value': udp_500_rule
                },
                {
                    'name': 'udp4500Rule',
                    'value': udp_4500_rule
                },
                {
                    'name': 'vpnTunnel',
                    'value': vpn_tunnel
                }
            ]
    }


def generate_network_url(project_id, network):
    """Format the resource name as a resource URI."""
    return 'projects/{}/global/networks/{}'.format(project_id, network)


def generate_router_url(project_id, region, router):
    """Format the resource name as a resource URI."""
    return 'projects/{}/regions/{}/routers/{}'.format(
        project_id,
        region,
        router
    )
