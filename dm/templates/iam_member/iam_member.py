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

    folder_id = context.properties.get('folderId')
    org_id = context.properties.get('organizationId')
    project_id = context.properties.get('projectId', context.env['project'])

    resources = []
    for ii, role in  enumerate(context.properties['roles']):
        for i, member in enumerate(role['members']):
            suffix =  sha1('{}-{}'.format(role['role'], member)).hexdigest()[:10]
            policy_get_name = '{}-{}'.format(context.env['name'], suffix)

            if org_id:
                resources.append({
                    'name': '{}-organization'.format(policy_get_name),
                    # TODO - Virtual type documentation needed
                    'type': 'gcp-types/cloudresourcemanager-v1:virtual.organizations.iamMemberBinding',
                    'properties': {
                        'resource': org_id,
                        'role': role['role'],
                        'member': member,
                    }
                })
            elif folder_id:
                resources.append({
                    'name': '{}-folder'.format(policy_get_name),
                    # TODO - Virtual type documentation needed
                    'type': 'gcp-types/cloudresourcemanager-v2:virtual.folders.iamMemberBinding',
                    'properties': {
                        'resource': folder_id,
                        'role': role['role'],
                        'member': member,
                    }
                })
            else:
                resources.append({
                    'name': '{}-project'.format(policy_get_name),
                    # TODO - Virtual type documentation needed
                    'type': 'gcp-types/cloudresourcemanager-v1:virtual.projects.iamMemberBinding',
                    'properties': {
                        'resource': project_id,
                        'role': role['role'],
                        'member': member,
                    }
                })

    return {"resources": resources}
