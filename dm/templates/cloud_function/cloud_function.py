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
"""
    Creates a Cloud Function from a local file system, a Cloud Storage bucket,
    or a Cloud Source Repository, and then assigns an HTTPS, Storage, or Pub/Sub
    trigger to that Cloud Function.
"""

NO_RESOURCES_OR_OUTPUTS = [], []

def get_source_url_output(function_name, context):
    """ Generates the Cloud Function output with a link to the source archive.
    """

    return {
        'name':  'sourceArchiveUrl',
        'value': '$(ref.{}.sourceArchiveUrl)'.format(function_name, context.env['name'])
    }

def append_cloud_storage_sources(function, project, context):
    """ Adds source code from the Cloud Storage. """

    properties = context.properties
    upload_path = properties.get('sourceArchiveUrl')

    resources = []
    outputs = [get_source_url_output(function['name'], context)]
    
    if not upload_path:
        msg = "sourceArchiveUrl must be provided"
        raise Exception(msg)

    function['properties']['sourceArchiveUrl'] = upload_path

    return resources, outputs

def append_cloud_repository_sources(function, context):
    """ Adds the source code from the cloud repository. """

    repo = context.properties.get('sourceRepository', {
        'url': context.properties.get('sourceRepositoryUrl')
    })
    function['properties']['sourceRepository'] = repo

    name = function['name']
    output = {
        'name': 'sourceRepositoryUrl',
        'value': '$(ref.{}.sourceRepository.deployedUrl)'.format(context.env['name'])
    }

    return [], [output]

def append_source_code(function, project, context):
    """ Append a reference to the Cloud Function's source code. """

    properties = context.properties

    if 'sourceRepository' in properties or 'sourceRepositoryUrl' in properties:
        return append_cloud_repository_sources(function, context)

    if 'sourceUploadUrl' in properties:
        append_optional_property(function, properties, 'sourceUploadUrl')
        return [], []

    if 'sourceArchiveUrl' in properties or 'localUploadPath' in properties:
        return append_cloud_storage_sources(function, project, context)

    raise ValueError('At least one of source properties must be provided')

def append_trigger_topic(function, properties):
    """ Appends the Pub/Sub event trigger. """

    topic = properties['triggerTopic']

    function['properties']['eventTrigger'] = {
        'eventType': 'providers/cloud.pubsub/eventTypes/topic.publish',
        'resource': topic
    }

    return NO_RESOURCES_OR_OUTPUTS

def append_trigger_http(function, context):
    """ Appends the HTTPS trigger and returns the generated URL. """

    function['properties']['httpsTrigger'] = {}
    output = {
        'name': 'httpsTriggerUrl',
        'value': '$(ref.{}.httpsTrigger.url)'.format(context.env['name'])
    }

    return [], [output]

def append_trigger_storage(function, context):
    """ Appends the Storage trigger. """

    bucket = context.properties['triggerStorage']['bucketName']
    event = context.properties['triggerStorage']['event']

    project_id = context.env['project']
    function['properties']['eventTrigger'] = {
        'eventType': 'google.storage.object.' + event,
        'resource': 'projects/{}/buckets/{}'.format(project_id, bucket)
    }

    return NO_RESOURCES_OR_OUTPUTS

def append_trigger(function, context):
    """ Adds the Trigger section and returns all the associated new
    resources and outputs.
    """

    if 'triggerTopic' in context.properties:
        return append_trigger_topic(function, context.properties)
    elif 'triggerStorage' in context.properties:
        return append_trigger_storage(function, context)

    return append_trigger_http(function, context)

def append_optional_property(function, properties, prop_name):
    """ If the property is set, it is added to the function body. """

    val = properties.get(prop_name)
    if val:
        function['properties'][prop_name] = val
    return

def create_function_resource(context):
    """ Creates the Cloud Function resource. """

    properties = context.properties
    name = properties.get('name', context.env['name'])
    project_id = properties.get('project', context.env['project'])
    location = properties.get('location', properties.get('region'))

    function = {
        # https://cloud.google.com/functions/docs/reference/rest/v1/projects.locations.functions
        'type': 'gcp-types/cloudfunctions-v1:projects.locations.functions',
        'name': context.env['name'],
        'properties':
            {
                'parent': 'projects/{}/locations/{}'.format(project_id, location),
                'function': name,
                # 'name': 'projects/{}/locations/{}/functions/{}'.format(project_id, location, name),
            },
    }

    optional_properties = ['entryPoint',
                           'labels',
                           'environmentVariables',
                           'timeout',
                           'runtime',
                           'maxInstances',
                           'availableMemoryMb',
                           'description']

    for prop in optional_properties:
        append_optional_property(function, properties, prop)

    trigger_resources, trigger_outputs = append_trigger(function, context)
    code_resources, code_outputs = append_source_code(function, project_id, context)

    if code_resources:
        function['metadata'] = {
            'dependsOn': [dep['name'] for dep in code_resources]
        }

    return (trigger_resources + code_resources + [function],
            trigger_outputs + code_outputs + [
                {
                    'name':  'region',
                    'value': context.properties['region']
                },
                {
                    'name': 'name',
                    'value': '$(ref.{}.name)'.format(context.env['name'])
                }
            ])

def generate_config(context):
    """ Entry point for the deployment resources. """

    resources, outputs = create_function_resource(context)

    return {
        'resources': resources,
        'outputs': outputs
    }
