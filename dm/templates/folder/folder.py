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
    This template creates a folder under an organization or under a
    parent folder.
"""


def generate_config(context):
    """ Entry point for the deployment resources. """

    resources = []
    out = {}
    for folder in context.properties.get('folders', []):

        create_folder = folder['name']

        parent = folder.get('orgId', folder.get('folderId'))

        resources.append(
            {
                'name': create_folder,
                'type': 'gcp-types/cloudresourcemanager-v2:folders',
                'properties':
                    {
                        'parent': parent,
                        'displayName': folder['displayName']
                    }
            }
        )

        out[create_folder] = {
            'name': '$(ref.{}.name)'.format(create_folder),
            'parent': '$(ref.{}.parent)'.format(create_folder),
            'displayName': '$(ref.{}.displayName)'.format(create_folder),
            'createTime': '$(ref.{}.createTime)'.format(create_folder),
            'lifecycleState': '$(ref.{}.lifecycleState)'.format(create_folder)
        }

    outputs = [{'name': 'folders', 'value': out}]

    return {'resources': resources, 'outputs': outputs}
