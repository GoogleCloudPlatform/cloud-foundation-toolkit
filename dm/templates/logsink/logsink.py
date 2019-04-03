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
""" This template creates a logsink (logging sink). """


def create_pubsub(context, logsink_name):
    """ Create the pubsub destination. """

    dest_properties = []
    if 'pubsubProperties' in context.properties:
        dest_prop = context.properties['pubsubProperties']
        access_control = dest_prop.get('accessControl', [])
        access_control.append(
            {
                'role': 'roles/pubsub.admin',
                'members': ['$(ref.' + logsink_name + '.writerIdentity)']
            }
        )

        dest_prop['accessControl'] = access_control
        dest_properties = [
            {
                'name': context.properties['destinationName'],
                'type': 'pubsub.py',
                'properties': dest_prop
            }
        ]

    return dest_properties


def create_bq_dataset(context, logsink_name):
    """ Create the BQ dataset destination. """

    dest_properties = []
    if 'bqProperties' in context.properties:
        dest_prop = context.properties['bqProperties']
        dest_prop['name'] = context.properties['destinationName']
        access = dest_prop.get('access', [])
        access.append(
            {
                'role': 'OWNER',
                'userByEmail': '$(ref.' + logsink_name + '.writerIdentity)'
            }
        )

        dest_prop['access'] = access
        dest_properties = [
            {
                'name': context.properties['destinationName'],
                'type': 'bigquery_dataset.py',
                'properties': dest_prop
            }
        ]

    return dest_properties


def create_storage(context, logsink_name):
    """ Create the bucket destination. """

    dest_properties = []
    if 'storageProperties' in context.properties:
        bucket_name = context.properties['destinationName']
        dest_prop = context.properties['storageProperties']
        dest_prop['name'] = bucket_name
        bindings = dest_prop.get('bindings', [])
        bindings.append(
            {
                'role': 'roles/storage.admin',
                'members': ['$(ref.' + logsink_name + '.writerIdentity)']
            }
        )

        # Do not set any IAM during the creation of the bucket since
        # we are going to set it afterwards
        if 'bindings' in dest_prop:
            del dest_prop['bindings']

        dest_properties = [
            {
                # Create the GCS Bucket
                'name': bucket_name,
                'type': 'gcs_bucket.py',
                'properties': dest_prop
            },
            {
                # Give the logsink writerIdentity permissions to the bucket
                'name': bucket_name + '-logging-storage-iampolicy',
                'action': 'gcp-types/storage-v1:storage.buckets.setIamPolicy',
                'properties':
                    {
                        'bucket': '$(ref.' + bucket_name + '.name)',
                        'project': context.env['project'],
                        'bindings': bindings
                    }
            }
        ]

    return dest_properties


def generate_config(context):
    """ Entry point for the deployment resources. """

    project_id = context.env['project']
    name = context.properties.get('name', context.env['name'])

    properties = {
        'name': name,
        'uniqueWriterIdentity': context.properties['uniqueWriterIdentity'],
        'sink': name
    }

    if 'orgId' in context.properties:
        source_id = str(context.properties.get('orgId'))
        source_type = 'organizations'
        properties['organization'] = str(source_id)
    elif 'billingAccountId' in context.properties:
        source_id = context.properties.get('billingAccountId')
        source_type = 'billingAccounts'
        del properties['sink']
    elif 'folderId' in context.properties:
        source_id = str(context.properties.get('folderId'))
        source_type = 'folders'
        properties['folder'] = str(source_id)
    elif 'projectId' in context.properties:
        source_id = context.properties.get('projectId')
        source_type = 'projects'

    properties['parent'] = '{}/{}'.format(source_type, source_id)

    dest_properties = []
    if context.properties['destinationType'] == 'pubsub':
        dest_properties = create_pubsub(context, name)
        destination = 'pubsub.googleapis.com/projects/{}/topics/{}'.format(
            project_id,
            context.properties['destinationName']
        )
    elif context.properties['destinationType'] == 'storage':
        dest_properties = create_storage(context, name)
        destination = 'storage.googleapis.com/{}'.format(
            context.properties['destinationName']
        )
    elif context.properties['destinationType'] == 'bigquery':
        # NOTE: There is a issue where BQ does not accept the uniqueWriter
        # returned by the logsink to be used in the userByEmail property.
        # Until that is resolved, this property is not supported.
        # dest_properties = create_bq_dataset(context, name)
        destination = 'bigquery.googleapis.com/projects/{}/datasets/{}'.format(
            project_id,
            context.properties['destinationName']
        )

    properties['destination'] = destination

    sink_filter = context.properties.get('filter')
    if sink_filter:
        properties['filter'] = sink_filter

    base_type = 'gcp-types/logging-v2:'
    resource = {
        'name': name,
        'type': base_type + source_type + '.sinks',
        'properties': properties
    }

    resources = [resource]

    if dest_properties:
        resources.extend(dest_properties)
        if context.properties['destinationType'] == 'storage':
            # GCS Bucket needs to be created first before the sink whereas
            # pub/sub and BQ do not. This might change in the future.
            resource['metadata'] = {
                'dependsOn': [context.properties['destinationName']]
            }

    return {
        'resources':
            resources,
        'outputs':
            [
                {
                    'name': 'writerIdentity',
                    'value': '$(ref.{}.writerIdentity)'.format(name)
                }
            ]
    }
