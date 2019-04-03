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
""" This template creates a Cloud Build resource. """


def generate_config(context):
    """ Entry point for the deployment resources. """

    resources = []
    outputs = []
    properties = context.properties
    name = context.env['name']
    build_steps = properties['steps']
    cloud_build = {
        'name': name,
        'action': 'gcp-types/cloudbuild-v1:cloudbuild.projects.builds.create',
        'properties': {
            'steps': build_steps
        },
        'metadata': {
            'runtimePolicy': ['UPDATE_ALWAYS']
        }
    }

    optional_properties = [
        'source',
        'timeout',
        'images',
        'artifacts',
        'logsBucket',
        'options',
        'substitutions',
        'tags',
        'secrets'
    ]

    for prop in optional_properties:
        if prop in properties:
            cloud_build['properties'][prop] = properties[prop]

    resources.append(cloud_build)

    # Output variables
    output_props = [
        'id',
        'status',
        'results',
        'createTime',
        'startTime',
        'finishTime',
        'logUrl',
        'sourceProvenance'
    ]

    for outprop in output_props:
        output_obj = {}
        output_obj['name'] = outprop
        output_obj['value'] = '$(ref.{}.{})'.format(name, outprop)
        outputs.append(output_obj)

    return {'resources': resources, 'outputs': outputs}
