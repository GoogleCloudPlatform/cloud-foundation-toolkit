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
""" This template creates a Pub/Sub (publish-subscribe) service. """

from hashlib import sha1
import json


def set_optional_property(destination, source, prop_name):
    """ Copies the property value if present. """

    if prop_name in source:
        destination[prop_name] = source[prop_name]

def create_subscription(resource_name, project_id, spec):
    """ Create a pull/push subscription from the simplified spec. """

    suffix = 'subscription-{}'.format(sha1((resource_name + json.dumps(spec)).encode('utf-8')).hexdigest()[:10])

    subscription = {
        'name': '{}-{}'.format(resource_name, suffix),
        # https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.subscriptions
        'type': 'gcp-types/pubsub-v1:projects.subscriptions',
        'properties':{
            'subscription': spec.get('name', suffix),
            'name': 'projects/{}/subscriptions/{}'.format(project_id, spec.get('name', suffix)),
            'topic': '$(ref.{}.name)'.format(resource_name)
        }
    }
    resources_list = [subscription]

    optional_properties = [
        'labels',
        'pushConfig',
        'ackDeadlineSeconds',
        'retainAckedMessages',
        'messageRetentionDuration',
        'expirationPolicy',
    ]

    for prop in optional_properties:
        set_optional_property(subscription['properties'], spec, prop)

    push_endpoint = spec.get('pushEndpoint')
    if push_endpoint is not None:
        subscription['properties']['pushConfig'] = {
            'pushEndpoint': push_endpoint,
        }

    return resources_list

def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    name = properties.get('name', properties.get('topic', context.env['name']))
    project_id = properties.get('project', context.env['project'])

    topic = {
        'name': context.env['name'],
        # https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.topics
        'type': 'gcp-types/pubsub-v1:projects.topics',
        'properties':{
            'topic': name,
            'name': 'projects/{}/topics/{}'.format(project_id, name),
        }
    }
    resources_list = [topic]

    optional_properties = [
        'labels',
    ]

    for prop in optional_properties:
        set_optional_property(topic['properties'], properties, prop)


    subscription_specs = properties.get('subscriptions', [])

    for spec in subscription_specs:
        resources_list = resources_list + create_subscription(context.env['name'], project_id, spec)

    return {
        'resources': resources_list,
        'outputs': [
            {
                'name': 'topicName',
                'value': '$(ref.{}.name)'.format(context.env['name'])
            }
        ],
    }
