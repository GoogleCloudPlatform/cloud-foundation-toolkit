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

""" This template creates a BigQuery table. """


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    name = properties.get('name', context.env['name'])
    project_id = properties.get('project', context.env['project'])

    properties = {
        'tableReference':
            {
                'tableId': name,
                'datasetId': context.properties['datasetId'],
                'projectId': project_id
            },
        'datasetId': context.properties['datasetId'],
        'projectId': project_id,
    }

    optional_properties = [
        'description',
        'friendlyName',
        'expirationTime',
        'schema',
        'timePartitioning',
        'externalDataConfiguration',
        'view'
    ]

    for prop in optional_properties:
        if prop in context.properties:
            if prop == 'schema':
                properties[prop] = {'fields': context.properties[prop]}
            else:
                properties[prop] = context.properties[prop]

    resources = [
        {
            # https://cloud.google.com/bigquery/docs/reference/rest/v2/tables
            'type': 'gcp-types/bigquery-v2:tables',
            'name': context.env['name'],
            'properties': properties
        }
    ]

    if 'dependsOn' in context.properties:
        resources[0]['metadata'] = {'dependsOn': context.properties['dependsOn']}

    outputs = [
        {
            'name': 'selfLink',
            'value': '$(ref.{}.selfLink)'.format(context.env['name'])
        },
        {
            'name': 'etag',
            'value': '$(ref.{}.etag)'.format(context.env['name'])
        },
        {
            'name': 'creationTime',
            'value': '$(ref.{}.creationTime)'.format(context.env['name'])
        },
        {
            'name': 'lastModifiedTime',
            'value': '$(ref.{}.lastModifiedTime)'.format(context.env['name'])
        },
        {
            'name': 'location',
            'value': '$(ref.{}.location)'.format(context.env['name'])
        },
        {
            'name': 'numBytes',
            'value': '$(ref.{}.numBytes)'.format(context.env['name'])
        },
        {
            'name': 'numLongTermBytes',
            'value': '$(ref.{}.numLongTermBytes)'.format(context.env['name'])
        },
        {
            'name': 'numRows',
            'value': '$(ref.{}.numRows)'.format(context.env['name'])
        },
        {
            'name': 'type',
            'value': '$(ref.{}.type)'.format(context.env['name'])
        }
    ]

    return {'resources': resources, 'outputs': outputs}
