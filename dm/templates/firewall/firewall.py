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
""" This template creates firewall rules for a network. """


def get_network(properties):
    """ Gets a network name. """

    network_name = properties.get('network')
    if network_name:
        is_self_link = '/' in network_name or '.' in network_name

        if is_self_link:
            network_url = network_name
        else:
            network_url = 'global/networks/{}'.format(network_name)

    return network_url


def generate_config(context):
    """ Entry point for the deployment resources. """

    network = context.properties.get('network')

    resources = []
    out = {}
    for i, rule in enumerate(context.properties['rules'], 1000):
        # Use VPC if specified in the properties. Otherwise, specify
        # the network URL in the config. If the network is not specified in
        # the config, the API defaults to 'global/networks/default'.
        if network and not rule.get('network'):
            rule['network'] = get_network(context.properties)

        rule['priority'] = rule.get('priority', i)
        resources.append(
            {
                'name': rule['name'],
                'type': 'compute.beta.firewall',
                'properties': rule
            }
        )

        out[rule['name']] = {
            'selfLink': '$(ref.' + rule['name'] + '.selfLink)',
            'creationTimestamp': '$(ref.' + rule['name']
                                 + '.creationTimestamp)',
        }

    outputs = [{'name': 'rules', 'value': out}]

    return {'resources': resources, 'outputs': outputs}
