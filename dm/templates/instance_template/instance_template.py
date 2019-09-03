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
""" This template creates an Instance Template. """


def set_optional_property(receiver, source, property_name, rename_to=None):
    """ If set, copies the given property value from one object to another
        and optionally rename it.
    """

    rename_to = rename_to or property_name
    if property_name in source:
        receiver[rename_to] = source[property_name]


def create_boot_disk(properties):
    """ Creates the boot disk configuration. """

    boot_disk = {
        'deviceName': 'boot',
        'type': 'PERSISTENT',
        'boot': True,
        'autoDelete': True,
        'initializeParams': {
            'sourceImage': properties['diskImage']
        }
    }

    for prop in ['diskSizeGb', 'diskType']:
        set_optional_property(boot_disk['initializeParams'], properties, prop)

    return boot_disk


def get_network_interfaces(properties):
    """ Get the configuration that connects the instance to an existing network
        and assigns to it an ephemeral public IP if specified.
    """
    network_interfaces = []

    networks = properties.get('networks', [])
    if len(networks) == 0 and properties.get('network'):
        network = {
            "network": properties.get('network'),
            "subnetwork": properties.get('subnetwork'),
            "networkIP": properties.get('networkIP'),
        }
        networks.append(network)
        if (properties.get('hasExternalIp')):
            network['accessConfigs'] = [{
                "type": "ONE_TO_ONE_NAT",
            }]
        if properties.get('natIP'):
          network['accessConfigs'][0]["natIp"] = properties.get('natIP')

    for network in networks:
        if not '.' in network['network'] and not '/' in network['network']:
            network_name = 'global/networks/{}'.format(network['network'])
        else:
            network_name = network['network']

        network_interface = {
            'network': network_name,
        }

        netif_optional_props = ['subnetwork', 'networkIP', 'aliasIpRanges', 'accessConfigs']
        for prop in netif_optional_props:
            if network.get(prop):
                network_interface[prop] = network[prop]
        network_interfaces.append(network_interface)

    return network_interfaces


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    name = properties.get('name', context.env['name'])
    machine_type = properties['machineType']
    network_interfaces = get_network_interfaces(context.properties)
    project_id = properties.get('project', context.env['project'])
    instance_template = {
        'name': context.env['name'],
        # https://cloud.google.com/compute/docs/reference/rest/v1/instanceTemplates
        'type': 'gcp-types/compute-v1:instanceTemplates',
        'properties':
            {
                'name': name,
                'project': project_id,
                'properties':
                    {
                        'machineType': machine_type,
                        'networkInterfaces': network_interfaces
                    }
            }
    }

    template_spec = instance_template['properties']['properties']

    optional_props = [
        'metadata',
        'disks',
        'scheduling',
        'tags',
        'canIpForward',
        'labels',
        'serviceAccounts',
        'scheduling',
        'shieldedInstanceConfig',
        'minCpuPlatform',
        'guestAccelerators',
    ]

    for prop in optional_props:
        set_optional_property(template_spec, properties, prop)
    if not template_spec.get('disks'):
        template_spec['disks'] = [create_boot_disk(properties)]

    set_optional_property(
        template_spec,
        properties,
        'instanceDescription',
        'description'
    )

    set_optional_property(
        instance_template['properties'],
        properties,
        'templateDescription',
        'description'
    )

    set_optional_property(
        instance_template['properties'],
        properties,
        'sourceInstance'
    )

    set_optional_property(
        instance_template['properties'],
        properties,
        'sourceInstanceParams'
    )

    return {
        'resources': [instance_template],
        'outputs':
            [
                {
                    'name': 'name',
                    'value': name
                },
                {
                    'name': 'selfLink',
                    'value': '$(ref.{}.selfLink)'.format(context.env['name'])
                }
            ]
    }
