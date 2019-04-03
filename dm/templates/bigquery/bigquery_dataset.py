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
""" This template creates a BigQuery dataset. """


def generate_config(context):
    """ Entry point for the deployment resources. """

    # You can modify the roles you wish to whitelist.
    whitelisted_roles = ['READER', 'WRITER', 'OWNER']

    name = context.properties['name']

    properties = {
        'datasetReference':
            {
                'datasetId': name,
                'projectId': context.env['project']
            },
        'location': context.properties['location']
    }

    optional_properties = ['description', 'defaultTableExpirationMs']

    for prop in optional_properties:
        if prop in context.properties:
            properties[prop] = context.properties[prop]

    if 'access' in context.properties:
        # Validate access roles.
        for access_role in context.properties['access']:
            if 'role' in access_role:
                role = access_role['role']
                if role not in whitelisted_roles:
                    raise ValueError(
                        'Role supplied \"{}\" for dataset \"{}\" not '
                        ' within the whitelist: {} '.format(
                            role,
                            context.properties['name'],
                            whitelisted_roles
                        )
                    )

        properties['access'] = context.properties['access']

        if context.properties.get('setDefaultOwner', False):
            # Build the default owner for the dataset.
            base = '@cloudservices.gserviceaccount.com'
            default_dataset_owner = context.env['project_number'] + base

            # Build the default access for the owner.
            owner_access = {
                'role': 'OWNER',
                'userByEmail': default_dataset_owner
            }
            properties['access'].append(owner_access)

    resources = [
        {
            'type': 'bigquery.v2.dataset',
            'name': name,
            'properties': properties
        }
    ]

    outputs = [
        {
            'name': 'selfLink',
            'value': '$(ref.{}.selfLink)'.format(name)
        },
        {
            'name': 'datasetId',
            'value': name
        },
        {
            'name': 'etag',
            'value': '$(ref.{}.etag)'.format(name)
        },
        {
            'name': 'creationTime',
            'value': '$(ref.{}.creationTime)'.format(name)
        },
        {
            'name': 'lastModifiedTime',
            'value': '$(ref.{}.lastModifiedTime)'.format(name)
        }
    ]

    return {'resources': resources, 'outputs': outputs}
