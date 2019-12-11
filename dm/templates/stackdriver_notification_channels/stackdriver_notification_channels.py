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
""" This template creates a Stackdriver Metric Descriptor. """


def get_condition_threshold(condition_properties):
    condition_threshold = {}
    properties = [
        'filter',
        'comparison',
        'duration',
        'thresholdValue',
        'trigger',
        'aggregations'
    ]
    for prop in condition_properties:
        if prop in properties:
            condition_threshold[prop] = condition_properties[prop]
    return condition_threshold


def get_policy_conditions(policy_conditions):
    conditions = []
    for condition in policy_conditions:
        conditions.append({
            'displayName': condition['displayName'],
            'conditionThreshold': get_condition_threshold(condition)
        })
    return conditions


def generate_config(context):
    """ Entry point for the deployment resources. """

    resources = []
    outputs = []
    properties = context.properties
    resource_name = context.env['name']
    project_id = properties.get('project', context.env['project'])
    notification_channels = properties["notificationChannels"]
    policies = properties["policies"]

    notification_channel_names = []

    for index, channel in enumerate(notification_channels):
        channel_type = channel['typeName']
        notification_channel_name = '{}-{}-{}'.format(resource_name, channel_type, index)
        notification_channel_names.append('$(ref.{}.name)'.format(notification_channel_name))
        resources.append({
            'name': notification_channel_name,
            'type': 'gcp-types/monitoring-v3:projects.notificationChannels',
            'properties': {
                'name': 'projects/{}'.format(project_id),
                'type': channel_type,
                'enabled': channel['channelEnabled'],
                'displayName': channel['displayName'],
                'labels': channel['labels']
            }
        })

    for index, policy in enumerate(policies):
        resources.append({
            'name': 'alerting-policy-{}-{}'.format(context.env['name'], index),
            'type': 'gcp-types/monitoring-v3:projects.alertPolicies',
            'properties': {
                'displayName': policy['name'],
                'documentation': {
                    'content': policy['documentationContent'],
                    'mimeType': policy['mimeType']
                },
                'combiner': policy['combiner'],
                'enabled': policy['policyEnabled'],
                'conditions': get_policy_conditions(policy['conditions']),
                'notificationChannels': notification_channel_names
            }
        })

    # Output variables:
    output_props = [
        'name',
        'type',
        'labels',
        'metricKind',
        'valueType',
        'unit',
        'description',
        'displayName',
        'metadata'
    ]

    for outprop in output_props:
        output = {}
        if outprop in properties:
            output['name'] = outprop
            output['value'] = '$(ref.{}.{})'.format(resource_name, outprop)
            outputs.append(output)

    return {'resources': resources, 'outputs': outputs}
