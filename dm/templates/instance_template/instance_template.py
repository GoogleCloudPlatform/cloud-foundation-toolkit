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

    networks = properties.get('networks', [{
        "name": properties.get('network'),
        "hasExternalIp": properties.get('hasExternalIp'),
        "natIP": properties.get('natIP'),
        "subnetwork": properties.get('subnetwork'),
        "networkIP": properties.get('networkIP'),
    }])

    for network in networks:
        if not '.' in network['name'] and not '/' in network['name']:
            network_name = 'global/networks/{}'.format(network['name'])
        else:
            network_name = network['name']

        network_interface = {
            'network': network_name,
        }

        if network['hasExternalIp']:
            access_configs = {
                'name': 'External NAT',
                'type': 'ONE_TO_ONE_NAT'
            }

            if network.get('natIP'):
                access_configs['natIP'] = network['natIP']

            network_interface['accessConfigs'] = [access_configs]

        netif_optional_props = ['subnetwork', 'networkIP']
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
    boot_disk = create_boot_disk(properties)
    network_interfaces = get_network_interfaces(context.properties)
    instance_template = {
        'name': name,
        'type': 'compute.v1.instanceTemplate',
        'properties':
            {
                'properties':
                    {
                        'machineType': machine_type,
                        'disks': [boot_disk],
                        'networkInterfaces': network_interfaces
                    }
            }
    }

    template_spec = instance_template['properties']['properties']

    optional_props = [
        'metadata',
        'tags',
        'canIpForward',
        'labels',
        'serviceAccounts',
        'scheduling'
    ]

    for prop in optional_props:
        set_optional_property(template_spec, properties, prop)

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
                    'value': '$(ref.{}.selfLink)'.format(name)
                }
            ]
    }
