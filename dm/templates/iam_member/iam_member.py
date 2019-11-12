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

def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    folder_id = properties.get('folderId')
    org_id = properties.get('organizationId')
    project_id = properties.get('projectId', context.env['project'])

    resources = []

    if 'dependsOn' in properties:
        dependson = { 'metadata': { 'dependsOn': properties['dependsOn'] } }
        dependson_root = properties['dependsOn']
    else:
        dependson = {}
        dependson_root = []

    for role in properties['roles']:
        for member in role['members']:
            suffix = sha1('{}-{}'.format(role['role'], member)).hexdigest()[:10]
            policy_get_name = '{}-{}'.format(context.env['name'], suffix)

            if org_id:
                resourse_name = '{}-organization'.format(policy_get_name)
                iam_resource = {
                    'name': resourse_name,
                    # TODO - Virtual type documentation needed
                    'type': 'gcp-types/cloudresourcemanager-v1:virtual.organizations.iamMemberBinding',
                    'properties': {
                        'resource': org_id,
                        'role': role['role'],
                        'member': member,
                    }
                }
                iam_resource.update(dependson)
                resources.append(iam_resource)
            elif folder_id:
                resourse_name = '{}-folder'.format(policy_get_name)
                iam_resource = {
                    'name': resourse_name,
                    # TODO - Virtual type documentation needed
                    'type': 'gcp-types/cloudresourcemanager-v2:virtual.folders.iamMemberBinding',
                    'properties': {
                        'resource': folder_id,
                        'role': role['role'],
                        'member': member,
                    }
                }
                iam_resource.update(dependson)
                resources.append(iam_resource)
            else:
                resourse_name = '{}-project'.format(policy_get_name)
                iam_resource = {
                    'name': (resourse_name),
                    # TODO - Virtual type documentation needed
                    'type': 'gcp-types/cloudresourcemanager-v1:virtual.projects.iamMemberBinding',
                    'properties': {
                        'resource': project_id,
                        'role': role['role'],
                        'member': member,
                    }
                }
                iam_resource.update(dependson)
                resources.append(iam_resource)

            dependson = { 'metadata': { 'dependsOn': [resourse_name] + dependson_root } }

    return {"resources": resources}
