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
""" This template creates a Cloud Build trigger. """


def generate_config(context):
    """ Entry point for the deployment resources. """

    resources = []
    properties = context.properties
    name = context.env['name']
    project_id = properties.get('project', context.env['project'])
    build_def = properties.get('build')
    build_filename = properties.get('filename')
    build_trigger_template = properties.get('triggerTemplate')
    build_github = properties.get('github')
    build_trigger_id = '$(ref.' + name + '.id)'
    build_trigger_createTime = '$(ref.' + name + '.createTime)'

    # build trigger create action
    build_trigger_create = {
        'name': name,
        # https://cloud.google.com/cloud-build/docs/api/reference/rest/v1/projects.triggers/create
        'action': 'gcp-types/cloudbuild-v1:cloudbuild.projects.triggers.create',
        'metadata': {
            'runtimePolicy': ['CREATE'],
        },
        'properties': {
            'name': name.replace('_', '-'),
            'projectId': project_id,
        }
    }

    # build trigger update action
    build_trigger_update = {
        'name': name + '-update',
        # https://cloud.google.com/cloud-build/docs/api/reference/rest/v1/projects.triggers/patch
        'action': 'gcp-types/cloudbuild-v1:cloudbuild.projects.triggers.patch',
        'metadata': {
            'runtimePolicy': ['UPDATE_ON_CHANGE'],
        },
        'properties': {
            'name': name.replace('_', '-'),
            'projectId': project_id,
            'id': build_trigger_id,
            'triggerId': build_trigger_id,
        }
    }

    optional_properties = [
        'description',
        'disabled',
        'substitutions',
        'ignoredFiles',
        'includedFiles',
        'tags'
    ]

    for prop in optional_properties:
        if prop in properties:
            build_trigger_create['properties'][prop] = properties[prop]
            build_trigger_update['properties'][prop] = properties[prop]

    if build_def:
        build_trigger_create['properties']['build'] = build_def
        build_trigger_update['properties']['build'] = build_def
    elif build_filename:
        build_trigger_create['properties']['filename'] = build_filename
        build_trigger_update['properties']['filename'] = build_filename

    if build_trigger_template:
        build_trigger_create['properties']['triggerTemplate'] = build_trigger_template
        build_trigger_update['properties']['triggerTemplate'] = build_trigger_template
    elif build_github:
        build_trigger_create['properties']['github'] = build_github
        build_trigger_update['properties']['github'] = build_github

    resources.append(build_trigger_create)
    resources.append(build_trigger_update)

    # build trigger delete action
    build_trigger_delete = {
        'name': name + '-delete',
        # https://cloud.google.com/cloud-build/docs/api/reference/rest/v1/projects.triggers/delete
        'action': 'gcp-types/cloudbuild-v1:cloudbuild.projects.triggers.delete',
        'metadata': {
            'runtimePolicy': ['DELETE'],
        },
        'properties': {
            'projectId': project_id,
            'triggerId': build_trigger_id
        }
    }

    resources.append(build_trigger_delete)

    # Output variables
    outputs = [
        {
            'name': 'id',
            'value': build_trigger_id
        },
        {
            'name': 'createTime',
            'value': build_trigger_createTime
        }
    ]

    return {'resources': resources, 'outputs': outputs}
