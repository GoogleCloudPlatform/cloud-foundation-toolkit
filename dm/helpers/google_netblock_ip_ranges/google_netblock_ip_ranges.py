# Copyright 2019 Google Inc. All rights reserved.
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
""" This template substitutes google netblock IP ranges into firewall rules."""

import yaml

def generate_config(context):
    google_netblock_ip_ranges = yaml.load(context.imports['google_netblock_ip_ranges.yaml'])
    properties = context.properties
    resource_type = properties['template']
    properties.pop('template', None)
    name = properties.get('name', context.env['name'])
    properties.pop('name', None)
    rules = []
    for rule in properties['rules']:
        rule_sub = rule
        if 'sourceRanges' in rule:
            rule_source = []
            for index, src_range in enumerate(rule['sourceRanges']):
                if 'google_netblock_ip_ranges' in src_range:
                    rule_source.extend(eval(src_range))
                else:
                    rule_source.append(src_range)
            rule_sub['sourceRanges'] = rule_source
        if 'destinationRanges' in rule:
            rule_destination = []
            for index, dst_range in enumerate(rule['destinationRanges']):
                if 'google_netblock_ip_ranges' in dst_range:
                    rule_destination.extend(eval(dst_range))
                else:
                    rule_destination.append(dst_range)
            rule_sub['destinationRanges'] = rule_destination
        rules.append(rule_sub)
    properties.update({'rules': rules})
    resources = {
        'name': name,
        'type': resource_type,
        'properties': properties
    }
    
    return {'resources': [resources]}
