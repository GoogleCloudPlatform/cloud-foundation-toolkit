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

from hashlib import sha1


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    project_id = properties.get('project', context.env['project'])
    network = properties.get('network')
    if network:
        if not ('/' in network or '.' in network):
            network = 'global/networks/{}'.format(network)
    else:
        network = 'projects/{}/global/networks/{}'.format(
            project_id,
            properties.get('networkName', 'default')
        )

    resources = []
    out = {}
    for i, rule in enumerate(properties['rules'], 1000):
        res_name = sha1(rule['name'].encode('utf-8')).hexdigest()[:10]

        rule['network'] = network
        rule['priority'] = rule.get('priority', i)
        rule['project'] = project_id
        resources.append(
            {
                'name': res_name,
                'type': 'gcp-types/compute-v1:firewalls',
                'properties': rule
            }
        )

        out[res_name] = {
            'selfLink': '$(ref.' + res_name + '.selfLink)',
            'creationTimestamp': '$(ref.' + res_name
                                 + '.creationTimestamp)',
        }

    outputs = [{'name': 'rules', 'value': out}]

    return {'resources': resources, 'outputs': outputs}
