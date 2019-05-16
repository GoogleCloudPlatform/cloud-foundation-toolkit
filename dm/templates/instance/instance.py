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

    for network in properties.get('networks'):
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

            if 'natIP' in network:
                access_configs['natIP'] = network['natIP']

            network_interface['accessConfigs'] = [access_configs]

        netif_optional_props = ['subnetwork', 'networkIP']
        for prop in netif_optional_props:
            if prop in network:
                network_interface[prop] = network[prop]
        network_interfaces.append(network_interface)

    return network_interfaces


def generate_config(context):
    """ Entry point for the deployment resources. """

    zone = context.properties['zone']
    vm_name = context.properties.get('name', context.env['name'])
    machine_type = context.properties['machineType']

    boot_disk = create_boot_disk(context.properties, zone, vm_name)
    network_interfaces = get_network_interfaces(context.properties)
    instance = {
        'name': vm_name,
        'type': 'compute.v1.instance',
        'properties':{
            'zone': zone,
            'machineType': 'zones/{}/machineTypes/{}'.format(zone,
                                                             machine_type),
            'disks': [boot_disk],
            'networkInterfaces': network_interfaces
        }
    }

    for name in ['metadata', 'serviceAccounts', 'canIpForward', 'tags']:
        set_optional_property(instance['properties'], context.properties, name)

    outputs = [
        {
            'name': 'networkInterfaces',
            'value': '$(ref.{}.networkInterfaces)'.format(vm_name)
        },
        {
            'name': 'name',
            'value': '$(ref.{}.name)'.format(vm_name)
        },
        {
            'name': 'selfLink',
            'value': '$(ref.{}.selfLink)'.format(vm_name)
        }
    ]

    return {'resources': [instance], 'outputs': outputs}
