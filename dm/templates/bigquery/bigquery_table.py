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

    name = context.properties['name']

    properties = {
        'tableReference':
            {
                'tableId': name,
                'datasetId': context.properties['datasetId'],
                'projectId': context.env['project']
            },
        'datasetId': context.properties['datasetId']
    }

    optional_properties = [
        'description',
        'friendlyName',
        'expirationTime',
        'schema',
        'timePartitioning',
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
            'type': 'bigquery.v2.table',
            'name': name,
            'properties': properties,
            'metadata': {
                'dependsOn': [context.properties['datasetId']]
            }
        }
    ]

    outputs = [
        {
            'name': 'selfLink',
            'value': '$(ref.{}.selfLink)'.format(name)
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
        },
        {
            'name': 'location',
            'value': '$(ref.{}.location)'.format(name)
        },
        {
            'name': 'numBytes',
            'value': '$(ref.{}.numBytes)'.format(name)
        },
        {
            'name': 'numLongTermBytes',
            'value': '$(ref.{}.numLongTermBytes)'.format(name)
        },
        {
            'name': 'numRows',
            'value': '$(ref.{}.numRows)'.format(name)
        },
        {
            'name': 'type',
            'value': '$(ref.{}.type)'.format(name)
        }
    ]

    return {'resources': resources, 'outputs': outputs}
