# Copyright 2017 Google Inc. All rights reserved.
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
""" This template creates an HA NAT gateway. """


SETUP_NATGATEWAY_SH = """#!/bin/bash
echo 1 > /proc/sys/net/ipv4/ip_forward
sysctl -w net.ipv4.ip_forward=1
echo "net.ipv4.ip_forward=1" | tee -a /etc/sysctl.conf > /dev/null
iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
apt-get -y install iptables-persistent
cat <<EOF > /usr/local/sbin/health-check-server.py
#!/usr/bin/python
from BaseHTTPServer import BaseHTTPRequestHandler,HTTPServer
import subprocess
PORT_NUMBER = 80
PING_HOST = "www.google.com"
def connectivityCheck():
  try:
    subprocess.check_call(["ping", "-c", "1", PING_HOST])
    return True
  except subprocess.CalledProcessError as e:
    return False
#This class will handle any incoming request
class myHandler(BaseHTTPRequestHandler):
  def do_GET(self):
    if self.path == '/health-check':
      if connectivityCheck():
        self.send_response(200)
      else:
        self.send_response(503)
    else:
      self.send_response(404)
try:
  server = HTTPServer(("", PORT_NUMBER), myHandler)
  print "Started httpserver on port " , PORT_NUMBER
  #Wait forever for incoming http requests
  server.serve_forever()
except KeyboardInterrupt:
  print "^C received, shutting down the web server"
  server.socket.close()
EOF
nohup python /usr/local/sbin/health-check-server.py >/dev/null 2>&1 &
#register a runtime config variable for a waiter to complete
CONFIG_NAME=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/runtime-config -H "Metadata-Flavor: Google")
VARIABLE_NAME=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/runtime-variable -H "Metadata-Flavor: Google")
gcloud beta runtime-config configs variables set $VARIABLE_NAME 1 --config-name $CONFIG_NAME 
"""

def get_network(properties):
    """ Gets a network name. """

    network_name = properties.get('network')
    is_self_link = '/' in network_name or '.' in network_name

    if is_self_link:
        network_url = network_name
    else:
        network_url = 'global/networks/{}'.format(network_name)

    return network_url


def get_subnetwork(context):
    """ Gets a subnetwork name. """

    subnet_name = context.properties.get('subnetwork')
    is_self_link = '/' in subnet_name or '.' in subnet_name

    if is_self_link:
        subnet_url = subnet_name
    else:
        subnet_url = 'projects/{}/regions/{}/subnetworks/{}'
        subnet_url = subnet_url.format(
            context.env['project'],
            context.properties['region'],
            subnet_name
        )

    return subnet_url


def get_healthcheck(name):
    """ Generate a healthcheck resource. """

    resource = {
        'name': name,
        'type': 'healthcheck.py',
        'properties':
            {
                'healthcheckType': 'HTTP',
                'port': 80,
                'requestPath': '/health-check',
                'healthyThreshold': 1,
                'unhealthyThreshold': 5,
                'checkIntervalSec': 30
            }
    }

    return resource


def get_firewall(context, network):
    """ Generate a firewall rule for the healthcheck. """

    # pylint: disable=line-too-long
    # See https://cloud.google.com/compute/docs/load-balancing/health-checks#health_check_source_ips_and_firewall_rules.
    name = context.env['name'] + '-healthcheck-firewall'
    resource = {
        'name': name,
        'type': 'firewall.py',
        'properties':
            {
                'network': network,
                'rules':
                    [
                        {
                            'name': name,
                            'allowed': [
                                {
                                    'IPProtocol': 'tcp',
                                    'ports': ['80'],
                                }
                            ],
                            'targetTags': [context.properties['natGatewayTag']],
                            'description':
                                'rule for allowing all health check traffic',
                            'sourceRanges': ['130.211.0.0/22',
                                             '35.191.0.0/16']
                        }
                    ]
            }
    }

    return resource


def get_external_internal_ip(ip_name,
                             external_ip_name,
                             internal_ip_name,
                             region,
                             subnet):

    """ Generate an external IP resource. """

    resource = {
        'name': ip_name,
        'type': 'ip_reservation.py',
        'properties':
            {
                'ipAddresses':
                    [
                        {
                            'name': external_ip_name,
                            'ipType': 'REGIONAL',
                            'region': region
                        },
                        {
                            'name': internal_ip_name,
                            'ipType': 'INTERNAL',
                            'region': region,
                            'subnetwork': subnet
                        }
                    ]
            }
    }

    return resource


def get_instance_template(context,
                          instance_template_name,
                          external_ip,
                          internal_ip,
                          network,
                          subnet):

    """ Generate an instance template resource. """

    resource = {
        'name': instance_template_name,
        'type': 'instance_template.py',
        'properties':
            {
                'natIP': external_ip,
                'network': network,
                'subnetwork': subnet,
                'networkIP': internal_ip,
                'diskImage': context.properties['imageType'],
                'machineType': context.properties['machineType'],
                'canIpForward': True,
                'diskType': context.properties['diskType'],
                'diskSizeGb': context.properties['diskSizeGb'],
                'tags': {
                    'items': [context.properties['natGatewayTag']]
                },
                'metadata':
                    {
                        'items':
                            [
                                {
                                    'key': 'startup-script',
                                    'value': SETUP_NATGATEWAY_SH
                                }
                            ]
                    },
            }
    }

    return resource


def get_route(context, route_name, internal_ip, network):
    """ Generate a route resource. """

    resource = {
        'name': route_name,
        'type': 'route.py',
        'properties':
            {
                'network': network,
                'routes':
                    [
                        {
                            'name': route_name + '-ip',
                            'routeType': 'ipaddress',
                            'nextHopIp': internal_ip,
                            'destRange': '0.0.0.0/0',
                            'priority': context.properties['routePriority'],
                            'tags': [context.properties['nattedVmTag']]
                        }
                    ]
            }
    }

    return resource


def get_managed_instance_group(name,
                               healthcheck,
                               instance_template_name,
                               base_instance_name,
                               zone):
    """ Generate a managed instance group resource. """

    resource = {
        'name': name,
        'type': 'compute.v1.instanceGroupManager',
        'properties':
            {
                'instanceTemplate':
                    '$(ref.' + instance_template_name + '.selfLink)',
                'baseInstanceName': base_instance_name,
                'zone': zone,
                'targetSize': 1,
                'autoHealingPolicies':
                    [
                        {
                            'healthCheck':
                                '$(ref.' + healthcheck + '.selfLink)',
                            'initialDelaySec': 120
                        }
                    ]
            }
    }

    return resource


def generate_config(context):
    """ Generate the deployment configuration. """

    resources = []
    prefix = context.env['name']
    hc_name = prefix + '-healthcheck'
    region = context.properties['region']
    network_name = get_network(context.properties)
    subnet_name = get_subnetwork(context)

    # Health check to be used by the managed instance groups.
    resources.append(get_healthcheck(hc_name))

    # Firewall rule that allows the healthcheck to work.
    resources.append(get_firewall(context, network_name))

    # Outputs:
    out = {}

    # Create a NAT gateway for each zone specified in the zones property.
    for zone in context.properties['zones']:

        # Reserve an internal/external static IP address.
        ip_name = prefix + '-ip-' + zone
        external_ip_name = prefix + '-ip-external-' + zone
        internal_ip_name = prefix + '-ip-internal-' + zone
        resources.append(
            get_external_internal_ip(
                ip_name,
                external_ip_name,
                internal_ip_name,
                region,
                subnet_name
            )
        )

        external_ip = '$(ref.{}.addresses.{}.address)'.format(
            ip_name,
            external_ip_name
        )

        internal_ip = '$(ref.{}.addresses.{}.address)'.format(
            ip_name,
            internal_ip_name
        )

        # Create a NAT gateway instance template.
        instance_template_name = prefix + '-insttempl-' + zone
        resources.append(
            get_instance_template(
                context,
                instance_template_name,
                external_ip,
                internal_ip,
                network_name,
                subnet_name
            )
        )

        # Create an Instance Group Manager for Healthcheck and AutoHealing.
        instance_group_manager_name = prefix + '-instgrpmgr-' + zone
        base_instance_name = prefix + '-gateway-' + zone
        resources.append(
            get_managed_instance_group(
                instance_group_manager_name,
                hc_name,
                instance_template_name,
                base_instance_name,
                zone
            )
        )

        # Create a route that will allow to use the NAT gateway VM as a
        # next hop.
        route_name = prefix + '-route-' + zone
        resources.append(
            get_route(context,
                      route_name,
                      internal_ip,
                      network_name)
        )

        # Set outputs grouped by the MIG name.
        out[base_instance_name] = {
            'instanceGroupManagerName': instance_group_manager_name,
            'instanceGroupmanagerSelflink': '$(ref.{}.selfLink)'.format(
                instance_group_manager_name
            ),
            'externalIP': external_ip,
            'internalIP': internal_ip,
            'instanceTemplateName': instance_template_name,
            'baseInstanceName': base_instance_name,
            'routeName': route_name,
            'zone': zone
        }

    outputs = [
        {
            'name': 'natGateways',
            'value': out
        },
        {
            'name': 'networkName',
            'value': network_name
        },
        {
            'name': 'subnetworkName',
            'value': subnet_name
        },
        {
            'name': 'natGatewayTag',
            'value': context.properties['natGatewayTag']
        },
        {
            'name': 'nattedVmTag',
            'value': context.properties['nattedVmTag']
        },
        {
            'name': 'region',
            'value': region
        },
        {
            'name': 'healthCheckName',
            'value': hc_name
        }
    ]

    return {'resources': resources, 'outputs': outputs}
