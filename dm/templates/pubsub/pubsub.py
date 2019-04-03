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

def create_subscription(resource_name, spec, topic_resource_name, spec_index):
    """ Create a pull/push subscription from the simplified spec. """

    subscription = {
        'name': '{}-subscription-{}'.format(resource_name, spec_index),
        'type': 'pubsub.v1.subscription',
        'properties':{
            'subscription': spec['name'],
            'topic': '$(ref.{}.name)'.format(topic_resource_name)
        }
    }

    push_endpoint = spec.get('pushEndpoint')
    if push_endpoint is not None:
        subscription['properties']['pushConfig'] = {
            'pushEndpoint': push_endpoint
        }

    ack_deadline_seconds = spec.get('ackDeadlineSeconds')
    if ack_deadline_seconds is not None:
        subscription['properties']['ackDeadlineSeconds'] = ack_deadline_seconds

    set_access_control(subscription, spec)

    return subscription

def create_iam_policy(bindings_spec):
    """ Create an IAM policy for the resource. """

    return {
        'gcpIamPolicy': {
            'bindings': bindings_spec
        }
    }

def set_access_control(resource, context):
    """ If necessary, define access control for the resource """

    access_control = context.get('accessControl')
    if access_control is not None:
        resource['accessControl'] = create_iam_policy(access_control)

def create_pubsub(resource_name, pubsub_spec):
    """ Create a topic with subscriptions. """

    topic_name = pubsub_spec.get('topic', resource_name)
    topic_resource_name = '{}-topic'.format(resource_name)
    topic = {
        'name': topic_resource_name,
        'type': 'pubsub.v1.topic',
        'properties':{
            'topic': topic_name
        }
    }

    set_access_control(topic, pubsub_spec)

    subscription_specs = pubsub_spec.get('subscriptions', [])
    subscriptions = [create_subscription(resource_name, spec,
                                         topic_resource_name, index)
                     for (index, spec)
                     in enumerate(subscription_specs, 1)]

    return [topic] + subscriptions

def create_topic_outputs(topic_resource):
    """ Create outputs for the topic. """

    return [
        {
            'name': 'topicName',
            'value': '$(ref.{}.name)'.format(topic_resource['name'])
        }
    ]

def generate_config(context):
    """ Entry point for the deployment resources. """

    resource_name = context.env['name']
    pubsub_resources = create_pubsub(resource_name, context.properties)
    pubsub_outputs = create_topic_outputs(pubsub_resources[0])

    return {
        'resources': pubsub_resources,
        'outputs': pubsub_outputs
    }
