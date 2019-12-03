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
""" This template creates an IAM policy member. """

from hashlib import sha1

mapper = {
    'organizationId': {
        'dm_type': 'gcp-types/cloudresourcemanager-v1:virtual.organizations.iamMemberBinding',
        'dm_resource_property': 'resource',
        'postfix': 'organization'},
    'folderId': {
        'dm_type': 'gcp-types/cloudresourcemanager-v2:virtual.folders.iamMemberBinding',
        'dm_resource_property': 'resource',
        'postfix': 'folder'},
    'projectId': {
        'dm_type': 'gcp-types/cloudresourcemanager-v1:virtual.projects.iamMemberBinding',
        'dm_resource_property': 'resource',
        'postfix': 'project'},
    'bucket': {
        'dm_type': 'gcp-types/storage-v1:virtual.buckets.iamMemberBinding',
        'dm_resource_property': 'bucket',
        'postfix': 'bucket'}
}


def get_type(context):
    for resource_type, resource_value in mapper.items():
        if resource_type in context.properties:
            resource_value.update({'id': context.properties[resource_type]})
            return resource_value

    # If nothing specified the default is projectID from context
    mapper['projectId'].update({'id': context.env['project']})
    return mapper['projectId']


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties

    base_resource = get_type(context)

    resources = []

    if 'dependsOn' in properties:
        dependson = {'metadata': {'dependsOn': properties['dependsOn']}}
        dependson_root = properties['dependsOn']
    else:
        dependson = {}
        dependson_root = []

    for role in properties['roles']:
        for member in role['members']:
            suffix = sha1(
                '{}-{}'.format(role['role'], member).encode('utf-8')).hexdigest()[:10]
            policy_get_name = '{}-{}'.format(context.env['name'], suffix)

            resource_name = '{}-{}'.format(policy_get_name,
                                           base_resource['postfix'])
            iam_resource = {
                'name': resource_name,
                # TODO - Virtual type documentation needed
                'type': base_resource['dm_type'],
                'properties': {
                    base_resource['dm_resource_property']: base_resource['id'],
                    'role': role['role'],
                    'member': member,
                }
            }
            iam_resource.update(dependson)
            resources.append(iam_resource)

            dependson = {'metadata': {'dependsOn': [
                resource_name] + dependson_root}}

    return {"resources": resources}
