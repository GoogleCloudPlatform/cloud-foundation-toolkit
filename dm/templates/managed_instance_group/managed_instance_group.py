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
""" This template creates a managed instance group. """

import copy

REGIONAL_LOCAL_IGM_TYPES = {
    # https://cloud.google.com/compute/docs/reference/rest/v1/regionInstanceGroupManagers
    True: 'gcp-types/compute-v1:regionInstanceGroupManagers',
    # https://cloud.google.com/compute/docs/reference/rest/v1/instanceGroupManagers
    False: 'gcp-types/compute-v1:instanceGroupManagers'
}


def set_optional_property(receiver, source, property_name):
    """ If set, copies the given property value from one object to another. """

    if property_name in source:
        receiver[property_name] = source[property_name]


def create_instance_template(properties, name_prefix):
    """ Creates an instance template resource. """

    name = name_prefix + '-it'

    instance_template = {
        'type': 'instance_template.py',
        'name': name,
        'properties': properties
    }

    self_link = '$(ref.{}.selfLink)'.format(name)

    return self_link, [instance_template], [
        {
            'name': 'instanceTemplateSelfLink',
            'value': self_link
        }
    ]


def get_instance_template(properties, name_prefix):
    """ If an instance template exists, returns a link to that template.
    If no instance template exists:
        (a) creates that template;
        (b) returns a link to it; and
        (c) returns resources/outputs that were required to create the template.
    """

    if 'url' in properties:
        return properties['url'], [], []

    return create_instance_template(properties, name_prefix)


def create_autoscaler(context, autoscaler_spec, igm):
    """ Creates an autoscaler. """

    igm_properties = igm['properties']

    autoscaler_properties = autoscaler_spec.copy()
    name =  '{}-autoscaler'.format(context.env['name'])

    autoscaler_properties['project'] = context.properties.get('project', context.env['project'])

    autoscaler_resource = {
        'type': 'autoscaler.py',
        'name': name,
        'properties': autoscaler_properties
    }

    # Use IGM's targetSize as maxNumReplicas
    autoscaler_properties['maxNumReplicas'] = igm_properties['targetSize']

    # And rename minSize to minNumReplicas
    min_size = autoscaler_properties.pop('minSize')
    autoscaler_properties['minNumReplicas'] = min_size

    autoscaler_properties['target'] = '$(ref.{}.selfLink)'.format(context.env['name'])

    for location in ['zone', 'region']:
        set_optional_property(autoscaler_properties, igm_properties, location)

    autoscaler_output = {
        'name': 'autoscalerSelfLink',
        'value': '$(ref.{}.selfLink)'.format(name)
    }

    return [autoscaler_resource], [autoscaler_output]


def get_autoscaler(context, igm):
    """ Creates an autoscaler, if necessary. """

    autoscaler_spec = context.properties.get('autoscaler')
    if autoscaler_spec:
        return create_autoscaler(context, autoscaler_spec, igm)

    return [], []


def get_igm_outputs(name, igm_properties):
    """ Creates Instance Group Manaher (IGM) resource outputs. """

    location_prop = 'region' if 'region' in igm_properties else 'zone'

    return [
        {
            'name': 'selfLink',
            'value': '$(ref.{}.selfLink)'.format(name)
        },
        {
            'name': 'name',
            'value': name
        },
        {
            'name': 'instanceGroupSelfLink',
            'value': '$(ref.{}.instanceGroup)'.format(name)
        },
        {
            'name': location_prop,
            'value': igm_properties[location_prop]
        }
    ]


def dereference_name(reference):
    """ Extracts resource name from Deployment Manager reference string. """

    # Extracting a name from `$(ref.NAME.property)` value results a string
    # which starts with `yaml%`. Remove the prefix.
    return reference.split('.')[1].replace('yaml%', '')


def is_reference(candidate):
    """ Checks if provided value is Deployment Manager reference string. """

    return candidate.strip().startswith('$(ref.')


def create_health_checks_assignment(healthchecks, igm_resource, project):
    """ Create resource for IGMs health checks assignment. """

    igm_properties = igm_resource['properties']
    igm_name = igm_properties['name']

    properties = {
        'instanceGroupManager': igm_name,
        'autoHealingPolicies': healthchecks,
        'project': project
    }

    dependencies = []
    metadata = {'dependsOn': dependencies}
    # Have to use a type-provider for health checks assignment
    # https://cloud.google.com/compute/docs/reference/rest/beta/regionInstanceGroupManagers/setAutoHealingPolicies
    # https://cloud.google.com/compute/docs/reference/rest/beta/instanceGroupManagers/setAutoHealingPolicies
    type_provider = 'gcp-types/compute-beta'
    action = '{}:compute.{}GroupManagers.setAutoHealingPolicies'.format(
        type_provider,
        'regionInstance' if 'region' in igm_properties else 'instance'
    )

    assign_healthcheck_resource = {
        'action': action,
        'name': igm_resource['name'] + '-set-hc',
        'properties': properties,
        'metadata': metadata
    }

    for healthcheck in healthchecks:
        if is_reference(healthcheck['healthCheck']):
            hc_resource_name = dereference_name(healthcheck['healthCheck'])
            dependencies.append(hc_resource_name)

    if dependencies:
        # instanceGroupManager must have a dependsOn metadata for all the
        # healthchecks it's going to use, so when the time comes, it's deleted
        # first
        igm_resource['metadata'] = copy.deepcopy(metadata)

    # setAutoHealingPolicies depends both on the health checks and IGM
    # resource
    dependencies.append(igm_resource['name'])

    for location in ['region', 'zone']:
        set_optional_property(properties, igm_properties, location)

    return assign_healthcheck_resource


def get_health_checks(properties, igm_resource, project):
    """ Assign health checks to IGM, if there're any. """

    if 'healthChecks' in properties:
        healthcheck_resources = create_health_checks_assignment(
            properties['healthChecks'],
            igm_resource,
            project
        )
        return [healthcheck_resources]

    return []


def get_igm(context, template_link):
    """ Creates the IGM resource with its outputs. """

    properties = context.properties
    name = properties.get('name', context.env['name'])
    project_id = properties.get('project', context.env['project'])
    is_regional = 'region' in properties

    igm_properties = {
        'name': name,
        'project': project_id,
        'instanceTemplate': template_link,
    }

    igm = {
        'name': context.env['name'],
        'type': REGIONAL_LOCAL_IGM_TYPES[is_regional],
        'properties': igm_properties
    }

    known_properties = [
        'description',
        'distributionPolicy',
        'namedPorts',
        'zone',
        'region',
        'targetSize',
        'baseInstanceName'
    ]

    for prop in known_properties:
        set_optional_property(igm_properties, properties, prop)

    outputs = get_igm_outputs(context.env['name'], igm_properties)

    return [igm], outputs


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    project_id = properties.get('project', context.env['project'])

    # Instance template
    properties['instanceTemplate']['project'] = project_id
    template = get_instance_template(properties['instanceTemplate'], context.env['name'])
    template_link, template_resources, template_outputs = template

    # Instance group manager
    igm_resources, igm_outputs = get_igm(context, template_link)
    igm = igm_resources[0]

    # Autoscaler
    autoscaler = get_autoscaler(context, igm)
    autoscaler_resources, autoscaler_outputs = autoscaler

    # Health checks
    healthcheck_resources = get_health_checks(properties, igm, project_id)

    return {
        'resources':
            igm_resources + template_resources + autoscaler_resources +
            healthcheck_resources,
        'outputs':
            igm_outputs + template_outputs + autoscaler_outputs
    }
