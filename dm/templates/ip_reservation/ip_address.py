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

""" This template creates an IP address. """


def get_address_type(ip_type):
    """ Return the address type to reserve. """

    if ip_type in ['GLOBAL', 'REGIONAL']:
        return 'EXTERNAL'

    return 'INTERNAL'

def get_resource_type(ip_type):
    """ Return the address resource type. """

    if ip_type == 'GLOBAL':
        # https://cloud.google.com/compute/docs/reference/rest/v1/globalAddresses
        return 'gcp-types/compute-v1:globalAddresses'

    # https://cloud.google.com/compute/docs/reference/rest/v1/addresses
    return 'gcp-types/compute-v1:addresses'


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    resource_type = get_resource_type(context.properties['ipType'])
    address_type = get_address_type(context.properties['ipType'])
    name = context.properties.get('name', context.env['name'])
    project_id = properties.get('project', context.env['project'])

    res_properties = {
        'addressType': address_type,
        'resourceType': 'addresses',
        'project': project_id,
    }

    optional_properties = [
        'subnetwork',
        'address',
        'description',
        'region',
        'networkTier',
        'prefixLength',
        'ipVersion',
        'purpose',
    ]

    for prop in optional_properties:
        if prop in context.properties:
            res_properties[prop] = str(context.properties[prop])

    resources = [
        {
            'name': name,
            'type': resource_type,
            'properties': res_properties
        }
    ]

    outputs = [
        {
            'name': 'selfLink',
            'value': '$(ref.{}.selfLink)'.format(name)
        },
        {
            'name': 'address',
            'value': '$(ref.{}.address)'.format(name)
        },
        {
            'name': 'status',
            'value': '$(ref.{}.status)'.format(name)
        }
    ]

    return {'resources': resources, 'outputs': outputs}
