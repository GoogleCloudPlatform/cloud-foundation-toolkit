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
""" Creates a runtimeConfig waiter resource. """


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    project_id = properties.get('project', context.env['project'])
    config_name = properties.get('config')

    props = {
        'waiter': properties.get('name', properties.get('waiter')),
        'parent': properties['parent'],
        'timeout': properties['timeout'],
        'success': properties['success'],
        # TODO: uncomment after gcp type is fixed
        # 'project': project_id,
    }

    optional_properties = ['failure']
    props.update({
        p: properties[p]
        for p in optional_properties if p in properties
    })

    resources = [{
        'name': context.env['name'],
        # https://cloud.google.com/deployment-manager/runtime-configurator/reference/rest/v1beta1/projects.configs.waiters
        'type': 'gcp-types/runtimeconfig-v1beta1:projects.configs.waiters',
        'properties': props,
        'metadata': {
            'dependsOn': [config_name]
        }
    }]

    outputs = [{
        'name': 'createTime',
        'value': '$(ref.{}.createTime)'.format(context.env['name'])
    }]

    return {'resources': resources, 'outputs': outputs}
