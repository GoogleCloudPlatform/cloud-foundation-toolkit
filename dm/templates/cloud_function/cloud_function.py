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

def get_source_url_output(function_name):
    """ Generates the Cloud Function output with a link to the source archive.
    """

    return {
        'name':  'sourceArchiveUrl',
        'value': '$(ref.{}.sourceArchiveUrl)'.format(function_name)
    }

def append_cloud_storage_sources(function, context):
    """ Adds source code from the Cloud Storage. """

    properties = context.properties
    upload_path = properties.get('sourceArchiveUrl')
    local_path = properties.get('localUploadPath')

    resources = []
    outputs = [get_source_url_output(function['name'])]

    if local_path:
        # The 'upload.py' file must be imported into the YAML file first.
        from upload import generate_upload_path, upload_source

        upload_path = upload_path or generate_upload_path()
        res = upload_source(function, context.imports, local_path, upload_path)
        source_resources, source_outputs = res
        resources += source_resources
        outputs += source_outputs
    elif not upload_path:
        msg = "Either localUploadPath or sourceArchiveUrl must be provided"
        raise Exception(msg)

    function['properties']['sourceArchiveUrl'] = upload_path

    return resources, outputs

def append_cloud_repository_sources(function, context):
    """ Adds the source code from the cloud repository. """

    append_optional_property(function,
                             context.properties,
                             'sourceRepositoryUrl')

    name = function['name']
    output = {
        'name': 'sourceRepositoryUrl',
        'value': '$(ref.{}.sourceRepository.repositoryUrl)'.format(name)
    }

    return [], [output]

def append_source_code(function, context):
    """ Append a reference to the Cloud Function's source code. """

    properties = context.properties
    if 'sourceArchiveUrl' in properties or 'localUploadPath' in properties:
        return append_cloud_storage_sources(function, context)
    elif 'sourceRepositoryUrl' in properties:
        return append_cloud_repository_sources(function, context)

    msg = """At least one of the following properties must be provided:
        - sourceRepositoryUrl
        - localUploadPath
        - sourceArchiveUrl"""
    raise ValueError(msg)

def append_trigger_topic(function, properties):
    """ Appends the Pub/Sub event trigger. """

    topic = properties['triggerTopic']

    function['properties']['eventTrigger'] = {
        'eventType': 'providers/cloud.pubsub/eventTypes/topic.publish',
        'resource': topic
    }

    return NO_RESOURCES_OR_OUTPUTS

def append_trigger_http(function):
    """ Appends the HTTPS trigger and returns the generated URL. """

    function['properties']['httpsTrigger'] = {}
    output = {
        'name': 'httpsTriggerUrl',
        'value': '$(ref.{}.httpsTrigger.url)'.format(function['name'])
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

    return append_trigger_http(function)

def append_optional_property(function, properties, prop_name):
    """ If the property is set, it is added to the function body. """

    val = properties.get(prop_name)
    if val:
        function['properties'][prop_name] = val
    return

def create_function_resource(resource_name, context):
    """ Creates the Cloud Function resource. """

    properties = context.properties
    region = properties['region']
    function_name = properties.get('name', resource_name)

    function = {
        'type': 'cloudfunctions.v1beta2.function',
        'name': function_name,
        'properties':
            {
                'location': region,
                'function': function_name,
            },
    }

    optional_properties = ['entryPoint',
                           'timeout',
                           'runtime',
                           'availableMemoryMb',
                           'description']

    for prop in optional_properties:
        append_optional_property(function, properties, prop)

    trigger_resources, trigger_outputs = append_trigger(function, context)
    code_resources, code_outputs = append_source_code(function, context)

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
                    'value': '$(ref.{}.name)'.format(function_name)
                }
            ])

def generate_config(context):
    """ Entry point for the deployment resources. """

    resource_name = context.env['name']
    resources, outputs = create_function_resource(resource_name, context)

    return {
        'resources': resources,
        'outputs': outputs
    }
