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
"""This template creates a custom IAM Project role."""


def generate_config(context):
    """ Entry point for the deployment resources. """

    project_id = context.env['project']
    included_permissions = context.properties['includedPermissions']

    role = {
        'includedPermissions': included_permissions,
        # Default the stage to General Availability.
        'stage': 'GA'
    }

    title = context.properties.get('title')
    if title:
        role['title'] = title

    description = context.properties.get('description')
    if description:
        role['description'] = description

    resources = [
        {
            'name': context.env['name'],
            'type': 'gcp-types/iam-v1:projects.roles',
            'properties':
                {
                    'parent': 'projects/' + project_id,
                    'roleId': context.properties['roleId'],
                    'role': role
                }
        }
    ]

    return {'resources': resources}
