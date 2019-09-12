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
""" This template creates a Compute Instance."""


def set_optional_property(receiver, source, property_name):
    """ If set, copies the given property value from one object to another. """

    if property_name in source:
        receiver[property_name] = source[property_name]


def create_boot_disk(properties, zone, instance_name):
    """ Create a boot disk configuration. """

    disk_name = instance_name
    boot_disk = {
        'deviceName': disk_name,
        'type': 'PERSISTENT',
        'boot': True,
        'autoDelete': True,
        'initializeParams': {
            'sourceImage': properties['diskImage']
        }
    }

    disk_params = boot_disk['initializeParams']
    set_optional_property(disk_params, properties, 'diskSizeGb')

    disk_type = properties.get('diskType')
    if disk_type:
        disk_params['diskType'] = 'zones/{}/diskTypes/{}'.format(zone,
                                                                 disk_type)

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
        if properties.get('hasExternalIp'):
            network['accessConfigs'] = [{
                "type": "ONE_TO_ONE_NAT",
            }]
            if properties.get('natIP'):
                network['accessConfigs'][0]['natIP'] = properties.get('natIP')

    for network in networks:
        if not '.' in network['network'] and not '/' in network['network']:
            network_name = 'global/networks/{}'.format(network['network'])
        else:
            network_name = network['network']

        network_interface = {
            'network': network_name,
        }

        netif_optional_props = ['subnetwork',
                                'networkIP', 'aliasIpRanges', 'accessConfigs']
        for prop in netif_optional_props:
            if network.get(prop):
                network_interface[prop] = network[prop]
        network_interfaces.append(network_interface)

    return network_interfaces


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    zone = properties['zone']
    vm_name = properties.get('name', context.env['name'])
    project_id = properties.get('project', context.env['project'])
    machine_type = properties['machineType']

    network_interfaces = get_network_interfaces(properties)
    instance = {
        'name': context.env['name'],
        # https://cloud.google.com/compute/docs/reference/rest/v1/instances
        'type': 'gcp-types/compute-v1:instances',
        'properties': {
            'name': vm_name,
            'zone': zone,
            'project': project_id,
            'machineType': 'zones/{}/machineTypes/{}'.format(zone,
                                                             machine_type),
            'networkInterfaces': network_interfaces
        }
    }

    optional_properties = [
        'description',
        'scheduling',
        'disks',
        'minCpuPlatform',
        'guestAccelerators',
        'deletionProtection',
        'hostname',
        'shieldedInstanceConfig',
        'shieldedInstanceIntegrityPolicy',
        'labels',
        'metadata',
        'serviceAccounts',
        'canIpForward',
        'tags',
    ]
    for name in optional_properties:
        set_optional_property(instance['properties'], properties, name)

    if not properties.get('disks'):
        instance['properties']['disks'] = [
            create_boot_disk(properties, zone, vm_name)]

    outputs = [
        {
            'name': 'networkInterfaces',
            'value': '$(ref.{}.networkInterfaces)'.format(context.env['name'])
        },
        {
            'name': 'name',
            'value': '$(ref.{}.name)'.format(context.env['name'])
        },
        {
            'name': 'selfLink',
            'value': '$(ref.{}.selfLink)'.format(context.env['name'])
        }
    ]

    if len(network_interfaces) == 1:
        outputs.append({
            'name': 'internalIp',
            'value': '$(ref.{}.networkInterfaces[0].networkIP)'.format(context.env['name'])
        })

        if 'accessConfigs' in network_interfaces[0]:
            access_configs = network_interfaces[0]['accessConfigs']
            for i, row in enumerate(access_configs, 0):
                if row['type'] == 'ONE_TO_ONE_NAT':
                    outputs.append({
                        'name': 'externalIp',
                        'value':
                            '$(ref.{}.networkInterfaces[0].accessConfigs[{}].natIP)'.format(
                                context.env['name'], i)
                    })
                    break

    return {'resources': [instance], 'outputs': outputs}
