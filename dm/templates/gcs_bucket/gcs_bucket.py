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
""" This template creates a Google Cloud Storage bucket. """

from hashlib import sha1


def generate_config(context):
    """ Entry point for the deployment resources. """

    resources = []
    properties = context.properties
    project_id = properties.get('project', context.env['project'])
    bucket_name = properties.get('name', context.env['name'])

    # output variables
    bucket_selflink = '$(ref.{}.selfLink)'.format(context.env['name'])
    bucket_uri = 'gs://' + bucket_name + '/'

    bucket = {
        'name': context.env['name'],
        # https://cloud.google.com/storage/docs/json_api/v1/buckets
        'type': 'gcp-types/storage-v1:buckets',
        'properties': {
            'project': project_id,
            'name': bucket_name
        }
    }

    requesterPays = context.properties.get('requesterPays')
    if requesterPays is not None:
      bucket['properties']['billing'] = {'requesterPays': requesterPays}

    optional_props = [
        'acl',
        'iamConfiguration',
        'retentionPolicy',
        'encryption',
        'defaultEventBasedHold',
        'cors',
        'defaultObjectAcl',
        'billing',
        'location',
        'versioning',
        'storageClass',
        'predefinedAcl',
        'predefinedDefaultObjectAcl',
        'logging',
        'lifecycle',
        'labels',
        'website'
    ]

    for prop in optional_props:
        if prop in properties:
            bucket['properties'][prop] = properties[prop]

    if not properties.get('iamConfiguration', {}).get('bucketPolicyOnly', {}).get('enabled', False):
        if 'predefinedAcl' not in bucket['properties']:
            bucket['properties']['predefinedAcl'] = 'private'
        if 'predefinedDefaultObjectAcl' not in bucket['properties']:
            bucket['properties']['predefinedDefaultObjectAcl'] = 'private'

    resources.append(bucket)

    # If IAM policy bindings are defined, apply these bindings.
    storage_provider_type = 'gcp-types/storage-v1:virtual.buckets.iamMemberBinding'
    bindings = properties.get('bindings', [])

    if 'dependsOn' in properties:
        dependson = { 'metadata': { 'dependsOn': properties['dependsOn'] } }
        dependson_root = properties['dependsOn']
    else:
        dependson = {}
        dependson_root = []

    if bindings:
        for role in bindings:
            for member in role['members']:
                suffix = sha1('{}-{}'.format(role['role'], member).encode('utf-8')).hexdigest()[:10]
                policy_get_name = '{}-{}'.format(context.env['name'], suffix)
                policy_name = '{}-iampolicy'.format(policy_get_name)
                iam_policy_resource = {
                    'name': policy_name,
                    # TODO - Virtual type documentation needed
                    'type': (storage_provider_type),
                    'properties':
                        {
                            'bucket': '$(ref.{}.name)'.format(context.env['name']),
                            'role': role['role'],
                            'member': member,
                        }
                }
                iam_policy_resource.update(dependson)
                resources.append(iam_policy_resource)
                dependson = { 'metadata': { 'dependsOn': [policy_name] + dependson_root } }

    if properties.get('billing', {}).get('requesterPays'):
        for resource in resources:
            resource['properties']['userProject'] = properties.get('userProject', context.env['project'])

    return {
        'resources':
            resources,
        'outputs':
            [
                {
                    'name': 'name',
                    'value': bucket_name
                },
                {
                    'name': 'selfLink',
                    'value': bucket_selflink
                },
                {
                    'name': 'url',
                    'value': bucket_uri
                }
            ]
    }
