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
""" This template creates a Dataproc cluster. """

PRIMARY_GROUP_SCHEMA = {'numInstances': None, 'machineType': 'machineTypeUri'}

SECONDARY_GROUP_SCHEMA = {'numInstances': None, 'isPreemptible': None}

GROUP_SCHEMAS = {
    'master': PRIMARY_GROUP_SCHEMA,
    'worker': PRIMARY_GROUP_SCHEMA,
    'secondaryWorker': SECONDARY_GROUP_SCHEMA
}


def get_disk_config(properties):
    """ If any disk property is specified, creates the diskConfig section. """

    disk_schema = {
        'diskType': 'bootDiskType',
        'diskSizeGb': 'bootDiskSizeGb',
        'numLocalSsds': None
    }

    return read_configuration(properties, disk_schema)


def read_configuration(properties, schema):
    """ Creates a new config section by reading and renaming properties from
    the source section.
    """

    if any(name in properties for name in schema):
        config = {}
        for name, rename_to in schema.iteritems():
            add_optional_property(config, properties, name, rename_to)
        return config

    return None


def get_instance_group_config(properties, image, cluster_schema):
    """ Creates a cluster instance group. """

    config = read_configuration(properties, cluster_schema)

    disk_config = get_disk_config(properties)
    if disk_config:
        config['diskConfig'] = disk_config

    if image:
        config['imageUri'] = image

    return config


def add_optional_property(destination, source, property_name, rename_to=None):
    """ Copies each property defined in the source object to the destination
    object.
    """

    rename_to = rename_to or property_name
    if property_name in source:
        destination[rename_to] = source[property_name]


def get_gce_cluster_config(properties):
    """ Creates the configuration section for a cluster. """

    gce_schema = {
        'zone': 'zoneUri',
        'network': 'networkUri',
        'subnetwork': 'subnetworkUri',
        'serviceAccountEmail': 'serviceAccount',
        'serviceAccountScopes': None,
        'internalIpOnly': None,
        'networkTags': 'tags',
        'metadata': None
    }

    if 'network' in properties and 'subnetwork' in properties:
        msg = 'Specifying both "network" and "subnetwork" is not allowed.'
        raise ValueError(msg)

    return read_configuration(properties, gce_schema)


def set_instance_group_config(properties, cluster, image, instance_group):
    """ Assign instance group config to the cluster. """

    group_spec = properties.get(instance_group)
    group_schema = GROUP_SCHEMAS[instance_group]
    group_config = get_instance_group_config(group_spec, image, group_schema)
    config_name = instance_group + 'Config'
    cluster['properties']['config'][config_name] = group_config
    config_output_path = 'ref.{}.config.{}'.format(cluster['name'], config_name)

    return {
        'name': '{}InstanceNames'.format(instance_group),
        'value': '$({}.instanceNames)'.format(config_output_path)
    }


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    name = properties.get('name', context.env['name'])
    project_id = context.env['project']
    image = context.properties.get('image')
    region = properties['region']

    cluster_config = get_gce_cluster_config(properties)

    cluster = {
        'name': name,
        'type': 'dataproc.v1.cluster',
        'properties':
            {
                'clusterName': name,
                'projectId': project_id,
                'region': region,
                'config': {
                    'gceClusterConfig': cluster_config,
                }
            }
    }

    for prop in ['configBucket', 'softwareConfig', 'initializationActions']:
        add_optional_property(cluster['properties']['config'], properties, prop)

    outputs = [
        {
            'name': 'name',
            'value': name
        },
        {
            'name': 'configBucket',
            'value': '$(ref.{}.config.configBucket)'.format(name)
        }
    ]

    for instance_group in ['master', 'worker', 'secondaryWorker']:
        if instance_group in properties:
            instance_group_output = set_instance_group_config(
                properties,
                cluster,
                image,
                instance_group
            )
            outputs.append(instance_group_output)

    return {'resources': [cluster], 'outputs': outputs}
