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
""" This template creates a Bastion host. """

import copy

IMAGE = 'projects/ubuntu-os-cloud/global/images/family/ubuntu-1804-lts'
DISABLE_SUDO_SCRIPT = """sudo cat /etc/sudoers | \\
    sed 's/%.*$//' | \\
    sed 's/#includedir.*$//' | \\
    sudo EDITOR=tee visudo"""
SSH = {'IPProtocol': 'tcp', 'ports': [22]}


def set_optional_property(destination, source, prop_name):
    """ Copies the property value if present. """

    if prop_name in source:
        destination[prop_name] = source[prop_name]


def find_metadata_item(metadata_items, key_name):
    """ Finds a metadata entry by the key name. """

    for item in metadata_items:
        if item['key'] == key_name:
            return item

    return None


def disable_sudo(bastion_props):
    """ Adds startup-script metadata that disables sudo. """

    metadata = bastion_props.get('metadata', {'items': []})
    startup_item = find_metadata_item(metadata['items'], 'startup-script')
    if not startup_item:
        startup_item = {'key': 'startup-script', 'value': ''}
        metadata['items'].append(startup_item)
    new_script = DISABLE_SUDO_SCRIPT + '\n' + startup_item['value']
    startup_item['value'] = new_script
    bastion_props['metadata'] = metadata


def get_ssh_firewall_rule(
        name,
        optional_properties,
        output_name,
        output_self_link
):
    """ Creates a new firewall rule with outputs. """

    ssh_props = {'allowed': [copy.deepcopy(SSH)]}

    ssh_rule = {
        'name': name,
        'type': 'compute.v1.firewall',
        'properties': ssh_props
    }

    for key, value in optional_properties.items():
        if value:
            ssh_props[key] = value

    return [ssh_rule], [
        {
            'name': output_name,
            'value': name
        },
        {
            'name': output_self_link,
            'value': '$(ref.{}.selfLink)'.format(name)
        },
    ]


def create_bastion_in_ssh_rule(bastion, firewall_settings):
    """ Creates a firewall rule for inbound SSH traffic. """

    to_bastion_rule = firewall_settings.get('sshToBastion')

    if to_bastion_rule:
        bastion_host_tag = to_bastion_rule['tag']

        # Append the Bastion tag, if it is not there yet.
        existing_tags = bastion['properties'].get('tags', {}).get('items', [])
        if not bastion_host_tag in existing_tags:
            existing_tags.append(bastion_host_tag)
            bastion['properties']['tags'] = {'items': existing_tags}

        rule_setup = {
            'sourceTags': to_bastion_rule.get('sourceTags'),
            'targetTags': [bastion_host_tag],
            'sourceRanges': to_bastion_rule.get('sourceRanges'),
            'priority': to_bastion_rule.get('priority'),
            'network': bastion['properties']['network']
        }

        return get_ssh_firewall_rule(
            to_bastion_rule['name'],
            rule_setup,
            'sshToBastionRuleName',
            'sshToBastionRuleSelfLink'
        )

    return [], []


def create_bastion_out_ssh_rule(bastion, firewall_settings):
    """ Creates a firewall rule for the Bastion outbound SSH traffic. """

    from_bastion_rule = firewall_settings.get('sshFromBastion')
    if from_bastion_rule:
        bastion_target_tag = from_bastion_rule.get('tag')

        # Calculate the firewall rule's source tags.
        if 'sshToBastion' in firewall_settings:
            bastion_host_tags = [firewall_settings['sshToBastion']['tag']]
        else:
            # Fall back to the instance tags collection.
            bastion_host_tags = bastion['properties'].get('tags', {})
            bastion_host_tags = bastion_host_tags.get('items', [])
            if bastion_host_tags:
                bastion_host_tags = copy.deepcopy(bastion_host_tags)
            else:
                msg = 'To enable SSH traffic from the Bastion host, at least one network tag must be assigned to it.'  # pylint: disable=line-too-long
                raise ValueError(msg)

        rule_setup = {
            'sourceTags': bastion_host_tags,
            'targetTags': [bastion_target_tag],
            'priority': from_bastion_rule.get('priority'),
            'network': bastion['properties']['network'],
        }

        return get_ssh_firewall_rule(
            from_bastion_rule['name'],
            rule_setup,
            'sshFromBastionRuleName',
            'sshFromBastionRuleSelfLink'
        )

    return [], []


def create_firewall_rules(bastion, firewall_settings):
    """ Creates in/out SSH rules for the Bastion host. """

    ssh_in_resources, ssh_in_outputs = create_bastion_in_ssh_rule(
        bastion,
        firewall_settings
    )

    ssh_out_resources, ssh_out_outputs = create_bastion_out_ssh_rule(
        bastion,
        firewall_settings
    )

    return (
        ssh_in_resources + ssh_out_resources,
        ssh_in_outputs + ssh_out_outputs
    )


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    name = properties.get('name', context.env['name'])

    bastion_props = {
        'zone': properties['zone'],
        'network': properties['network'],
        'machineType': properties['machineType'],
        'diskImage': IMAGE
    }

    bastion = {'name': name, 'type': 'instance.py', 'properties': bastion_props}

    optional_props = ['diskSizeGb', 'metadata', 'tags']

    for prop in optional_props:
        set_optional_property(bastion_props, properties, prop)

    if properties.get('disableSudo'):
        disable_sudo(bastion_props)

    firewall_settings = properties.get('createFirewallRules')
    if firewall_settings:
        extra_resources, extra_outputs = create_firewall_rules(
            bastion,
            firewall_settings
        )
    else:
        extra_resources = []
        extra_outputs = []

    outputs = [
        {
            'name': 'name',
            'value': name
        },
        {
            'name': 'selfLink',
            'value': '$(ref.{}.selfLink)'.format(name)
        },
        {
            'name': 'internalIp',
            'value': '$(ref.{}.internalIp)'.format(name)
        },
        {
            'name': 'externalIp',
            'value': '$(ref.{}.externalIp)'.format(name)
        }
    ]

    return {
        'resources': [bastion] + extra_resources,
        'outputs': outputs + extra_outputs
    }
