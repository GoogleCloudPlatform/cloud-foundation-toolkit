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
""" This template creates a Cloud Task resource. """


def generate_queue_name(context):
    """ Create queue name based on input. """

    if ('projects/' in context.properties['queueId'] or
            '$(ref.' in context.properties['queueId']):
        # Full queue name or reference
        queue_name = context.properties['queueId']
    else:
        # Format the queue name
        project_id = context.properties.get('projectId', context.env['project'])
        queue_name = 'projects/{}/locations/{}/queues/{}'.format(
            project_id,
            context.properties['location'],
            context.properties['queueId']
        )

    return queue_name


def generate_config(context):
    """ Entry point for the deployment resources. """

    resources = []
    properties = context.properties
    name = context.env['name']
    project_id = properties.get('projectId', context.env['project'])
    parent = generate_queue_name(context)

    task = {
        'name':
            name,
        'type':
            '{}/cloudtasks:projects.locations.queues.tasks'.format(project_id),
        'properties':
            {
                'parent': parent,
                'task':
                    {
                        'name':
                            '{}/tasks/{}'.format(parent,
                                                 name),
                        'appEngineHttpRequest':
                            properties['task']['appEngineHttpRequest']
                    }
            }
    }

    optional_properties = ['scheduleTime']

    for prop in optional_properties:
        if prop in properties['task']:
            task['properties']['task'][prop] = properties['task'][prop]

    resources.append(task)

    return {
        'resources': resources,
        'outputs': [
            {
                'name':'name',
                'value': '$(ref.{}.name)'.format(name)
            },
            {
                'name':'createTime',
                'value': '$(ref.{}.createTime)'.format(name)
            },
            {
                'name':'view',
                'value': '$(ref.{}.view)'.format(name)
            },
            {
                'name':'scheduleTime',
                'value': '$(ref.{}.scheduleTime)'.format(name)
            }
        ]
    }
