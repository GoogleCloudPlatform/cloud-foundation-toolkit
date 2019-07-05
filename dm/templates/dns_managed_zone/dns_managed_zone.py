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
"""This template creates a managed zone resource in the Cloud DNS."""


def generate_config(context):
    """Entry point for the deployment resources."""
    # Backward Compatibility with the old property `zoneName`
    try:
        managed_zone_name = context.properties['zoneName']
    except KeyError:
        managed_zone_name = context.properties.get('name', context.env['name'])
    dnsname = context.properties['dnsName']
    managed_zone_description = context.properties['description']
    name_servers = '$(ref.' + context.env['name'] + '.nameServers)'
    project_id = context.properties.get('project', context.env['project'])

    resources = []
    outputs = [
        {
            'name': 'dnsName',
            'value': dnsname
        },
        {
            'name': 'managedZoneDescription',
            'value': managed_zone_description
        },
        {
            'name': 'nameServers',
            'value': name_servers
        },
        {
            'name': 'managedZoneName',
            'value': managed_zone_name
        }
    ]

    managed_zone = {
        'name': context.env['name'],
        # https://cloud.google.com/dns/docs/reference/v1/managedZones
        'type': 'gcp-types/dns-v1:managedZones',
        'properties': {
            'name': managed_zone_name,
            'dnsName': dnsname,
            'description': managed_zone_description,
            'project_id': project_id
        }
    }

    # making resources and outputs additional properties
    for prop in context.properties:
        if prop not in managed_zone['properties']:
            managed_zone['properties'][prop] = context.properties[prop]
            outputs.append(
                {
                    'name': prop,
                    'value': context.properties[prop]
                }
            )
    resources.append(managed_zone)

    return {'resources': resources, 'outputs': outputs}
