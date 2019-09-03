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
"""This template creates a custom IAM Organization role."""


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    org_id = str(properties['orgId'])
    included_permissions = properties['includedPermissions']

    role = {
        'includedPermissions': included_permissions,
        # Default the stage to General Availability.
        'stage': properties.get('stage')
    }

    title = properties.get('title')
    if title:
        role['title'] = title

    description = properties.get('description')
    if description:
        role['description'] = description

    resources = [
        {
            'name': context.env['name'],
            # https://cloud.google.com/iam/reference/rest/v1/organizations.roles
            'type': 'gcp-types/iam-v1:organizations.roles',
            'properties':
                {
                    'parent': 'organizations/' + org_id,
                    'roleId': properties['roleId'],
                    'role': role
                }
        }
    ]

    return {'resources': resources}
