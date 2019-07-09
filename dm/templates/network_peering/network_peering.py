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
""" This template creates a VPC network peering. """


def generate_config(context):
    """ Entry point for the deployment resources. """

    resources = []
    properties = context.properties
    name = properties.get('name', context.env['name'])
    project_id = properties.get('project', context.env['project'])

    peer_create = {
        'name': context.env['name'] + '-create-peer',
        # https://cloud.google.com/compute/docs/reference/rest/v1/networks/addPeering
        'action': 'gcp-types/compute-v1:compute.networks.addPeering',
        'metadata': {
            'runtimePolicy': ['CREATE',
                             ]
        },
        'properties':
            {
                'name': name,
                'project': project_id,
                'network': properties['network'],
                'peerNetwork': properties['peerNetwork'],
                'autoCreateRoutes': properties.get('autoCreateRoutes')
            }
    }

    peer_delete = {
        'name': context.env['name'] + '-delete-peer',
        # https://cloud.google.com/compute/docs/reference/rest/v1/networks/removePeering
        'action': 'gcp-types/compute-v1:compute.networks.removePeering',
        'metadata': {
            'runtimePolicy': ['DELETE',
                             ]
        },
        'properties':
            {
                'name': name,
                'project': project_id,
                'network': properties['network'],
                'peerNetwork': properties['peerNetwork']
            }
    }

    resources.append(peer_create)
    resources.append(peer_delete)

    # As peerings are added/removed to/from a network, adding a peer does not
    # expose ouputs. For peering state and state_details, query the `peerings`
    # output on the network resource.

    return {'resources': resources}
