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
This template creates a Runtime Configurator with the associated resources.
"""


def generate_config(context):
    """ Entry point for the deployment resources. """

    resources = []
    properties = context.properties
    project_id = properties.get('projectId', context.env['project'])
    name = properties.get('config', context.env['name'])
    parent = 'projects/{}/configs/{}'.format(project_id, name)

    # The runtimeconfig resource.
    runtime_config = {
        'name': name,
        'type': 'runtimeconfig.v1beta1.config',
        'properties': {
            'config': name,
            'description': properties['description']
        }
    }

    resources.append(runtime_config)

    # The runtimeconfig variable resources.
    for variable in properties.get('variables', []):
        variable['parent'] = parent
        variable['config'] = name
        variable_res = {
            'name': variable['variable'],
            'type': 'variable.py',
            'properties': variable
        }
        resources.append(variable_res)

    # The runtimeconfig waiter resources.
    for waiter in properties.get('waiters', []):
        waiter['parent'] = parent
        waiter['config'] = name
        waiter_res = {
            'name': waiter['waiter'],
            'type': 'waiter.py',
            'properties': waiter
        }
        resources.append(waiter_res)

    outputs = [{'name': 'configName', 'value': '$(ref.{}.name)'.format(name)}]

    return {'resources': resources, 'outputs': outputs}
