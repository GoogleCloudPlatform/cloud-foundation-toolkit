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
""" This template creates a Cloud SQL Instance with databases and users. """

import collections
import random
import string

DMBundle = collections.namedtuple('DMBundle', 'resource outputs')

SUFFIX_LENGTH = 5
CHAR_CHOICE = string.digits + string.ascii_lowercase


def get_random_string(length):
    """ Creates a random string of characters of the specified length. """

    return ''.join([random.choice(CHAR_CHOICE) for _ in range(length)])


def set_optional_property(receiver, source, property_name):
    """ If set, copies the given property value from one object to another. """

    if property_name in source:
        receiver[property_name] = source[property_name]


def get_instance(res_name, project_id, properties):
    """ Creates a Cloud SQL instance. """

    name = res_name
    instance_properties = {
        'region': properties['region'],
        'project': project_id,
        'name': name
    }

    optional_properties = [
        'databaseVersion',
        'failoverReplica',
        'instanceType',
        'masterInstanceName',
        'maxDiskSize',
        'onPremisesConfiguration',
        'replicaConfiguration',
        'serverCaCert',
        'serviceAccountEmailAddress',
        'settings',
    ]

    for prop in optional_properties:
        set_optional_property(instance_properties, properties, prop)

    instance = {
        'name': name,
        'type': 'sqladmin.v1beta4.instance',
        'properties': instance_properties
    }

    if 'dependsOn' in properties:
        instance['metadata'] = {'dependsOn': properties['dependsOn']}

    outputs = [
        {
            'name': 'name',
            'value': '$(ref.{}.name)'.format(name)
        },
        {
            'name': 'selfLink',
            'value': '$(ref.{}.selfLink)'.format(name)
        },
        {
            'name': 'gceZone',
            'value': '$(ref.{}.gceZone)'.format(name)
        },
        {
            'name': 'connectionName',
            'value': '$(ref.{}.connectionName)'.format(name)
        },
    ]

    return DMBundle(instance, outputs)


def get_database(instance_name, project_id, properties):
    """ Creates a Cloud SQL database. """

    name = properties['name']
    res_name = name

    db_properties = {
        'name': name,
        'project': project_id,
        'instance': instance_name
    }

    optional_properties = [
        'charset',
        'collation',
        'instance',
    ]

    for prop in optional_properties:
        set_optional_property(db_properties, properties, prop)

    database = {
        'name': res_name,
        'type': 'sqladmin.v1beta4.database',
        'properties': db_properties
    }

    outputs = [
        {
            'name': 'name',
            'value': '$(ref.{}.name)'.format(res_name)
        },
        {
            'name': 'selfLink',
            'value': '$(ref.{}.selfLink)'.format(res_name)
        }
    ]

    return DMBundle(database, outputs)


def get_databases(instance_name, project_id, properties):
    """ Creates Cloud SQL databases for the given instance. """

    dbs = properties.get('databases')
    if dbs:
        return [get_database(instance_name, project_id, db) for db in dbs]

    return []


def get_user(instance_name, project_id, properties):
    """ Creates a Cloud SQL user. """

    name = properties['name']
    res_name = 'cloud-sql-{}'.format(name)
    user_properties = {
        'name': name,
        'project': project_id,
        'instance': instance_name,
    }

    for prop in ['host', 'password']:
        set_optional_property(user_properties, properties, prop)

    user = {
        'name': res_name,
        'type': 'sqladmin.v1beta4.user',
        'properties': user_properties
    }

    outputs = [{'name': 'name', 'value': name}]

    return DMBundle(user, outputs)


def get_users(instance_name, project_id, properties):
    """ Creates Cloud SQL users for the given instance. """

    users = properties.get('users')
    if users:
        return [get_user(instance_name, project_id, user) for user in users]

    return []


def create_sequentially(resources):
    """
    Sets up the resources' metadata in such a way that the resources are
    created sequentially.
    """

    if resources and len(resources) > 1:
        previous = resources[0]
        for current in resources[1:]:
            previous_name = previous['name']
            current['metadata'] = {'dependsOn': [previous_name]}
            previous = current


def consolidate_outputs(bundles, prefix):
    """
    Consolidates values of multiple outputs into one array for the new
    output.
    """

    res = {}
    outputs = [output for bundle in bundles for output in bundle.outputs]
    for output in outputs:
        output_name = output['name']
        new_name = prefix + output_name[0].upper() + output_name[1:] + 's'
        if not new_name in res:
            res[new_name] = {'name': new_name, 'value': []}
        res[new_name]['value'].append(output['value'])

    return [value for _, value in res.items()]


def get_resource_names_output(resources):
    """
    Creates the output dict with the names of all resources to be created.
    """

    names = [resource['name'] for resource in resources]

    return {'name': 'resources', 'value': names}


def generate_config(context):
    """ Creates the Cloud SQL instance, databases, and user. """

    properties = context.properties
    res_name = properties.get('name', context.env['name'])
    project_id = properties.get('project', context.env['project'])

    instance = get_instance(res_name, project_id, properties)
    instance_name = instance.outputs[0]['value']  # 'name' output

    users = get_users(instance_name, project_id, properties)
    dbs = get_databases(instance_name, project_id, properties)

    children = [user.resource for user in users] + [db.resource for db in dbs]
    create_sequentially(children)

    user_outputs = consolidate_outputs(users, 'user')
    db_outputs = consolidate_outputs(dbs, 'database')

    resources = [instance.resource] + children
    outputs = [get_resource_names_output(resources)] + instance.outputs + \
              db_outputs + user_outputs

    return {'resources': resources, 'outputs': outputs}
