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
""" This template creates Forseti Security tools and resources. """

import collections
import random
import string
import copy

# The helper tuple for handling resources and their outputs.
DMResource = collections.namedtuple(
    'DMResource',
    'self_link resources outputs'
)

SUFFIX_LENGTH = 10
CHAR_CHOICE = string.digits + string.ascii_lowercase
FORSETI_APIS = [
    'admin.googleapis.com',
    'appengine.googleapis.com',
    'bigquery-json.googleapis.com',
    'cloudbilling.googleapis.com',
    'iam.googleapis.com',
    'cloudresourcemanager.googleapis.com',
    'sqladmin.googleapis.com',
    'sql-component.googleapis.com',
    'compute.googleapis.com',
    'deploymentmanager.googleapis.com'
]
PROJECT_REMOVE_SA = True
PROJECT_REMOVE_VPC = True
CLOUD_MAN = 'gcp-types/cloudresourcemanager-v1:cloudresourcemanager'
IAM = 'gcp-types/iam-v1:iam.projects'
STORAGE = 'gcp-types/storage-v1:storage'

# If True, the organization and project policies that were previously added for
# the Forseti service accounts are removed on deployment deletion.
# However, due to the nature of the IAM patches, this step fails if the
# underlying policy has changed after the deployment creation.

CLEANUP_POLICY_ON_DELETE = True

def get_random_string(length):
    """ Generates a random string of a given length. """

    return ''.join([random.choice(CHAR_CHOICE) for _ in range(length)])

def generate_project_id(prefix):
    """ Generates a new project ID. """

    return prefix + '-' + get_random_string(SUFFIX_LENGTH)

def create_forseti_project(deployment_name, properties):
    """ Generates a new project for the Forseti tools. """

    project_id = properties.get('id', generate_project_id(deployment_name))
    project_name = properties.get('name', project_id)

    project = {
        'type': 'project.py',
        'name': project_id,
        'properties': {
            'name': project_name,
            'parent': copy.deepcopy(properties['parent']),
            'billingAccountId': properties['billingAccountId'],
            'activateApis': FORSETI_APIS,
            'removeDefaultSA': PROJECT_REMOVE_SA,
            'removeDefaultVPC': PROJECT_REMOVE_VPC,
            'serviceAccounts': [],
            'groups': []
        }
    }

    return DMResource(
        self_link=get_ref(project_id, 'projectId'),
        resources=[project],
        outputs=[
            {
                'name': 'projectId',
                'value': project_id
            },
            {
                'name': 'resources',
                'value': get_ref(project_id, 'resources')
            }
        ]
    )

def wait_for_init_complete(project, *deps):
    """ Adds explicit dependsOn metadata to the project dependencies. """

    resources_output = find_output_value('resources', project.outputs)
    for dependency in deps:
        for resource in dependency.resources:
            if 'type' in resource or 'getIamPolicy' in resource['action']:
                resource['metadata'] = {'dependsOn': resources_output}

def get_forseti_project(deployment_name, properties):
    """ Gets a reference to a project for the Forseti resources. """

    create = properties['create']

    if create:
        return create_forseti_project(deployment_name, properties)

    return DMResource(self_link=properties['id'], resources=[], outputs=[])

def create_policy_bindings(member, roles):
    """ Converts member+roles args to a proper policy bindings object. """

    bindings = []

    for role in roles:
        bindings.append({'role': role, 'members': [member]})

    return bindings

def get_action_path(res_type):
    """
    Gets a proper type provider path for assigning the IAM policy for a given
    resource type.
    """

    if res_type == 'serviceAccount':
        type_provider = IAM
    elif res_type == 'bucket':
        type_provider = STORAGE
    else:
        type_provider = CLOUD_MAN

    return '{}.{}s'.format(type_provider, res_type)

def set_member_roles(member, roles, res_type, res_id, project_id):
    """ Sets the IAM policy of a given resource. """

    random_suffix = get_random_string(SUFFIX_LENGTH)

    bindings = create_policy_bindings(member, roles)

    action_path = get_action_path(res_type)

    if res_type == 'bucket':
        properties = {
            'bucket': res_id,
            'project': project_id,
            'bindings': bindings
        }
    else:
        properties = {
            'resource': res_id,
            'policy': {'bindings': bindings}
        }

    resources = [
        {
            'name': 'set-iam-policy-' + random_suffix,
            'action': '{}.setIamPolicy'.format(action_path),
            'properties': properties
        },
    ]

    return DMResource(None, resources, [])

def patch_member_roles(member, roles, res_type, res_id):
    """ Patches the IAM policy of a given resource. """

    random_suffix = get_random_string(SUFFIX_LENGTH)
    get_iam_policy_name = 'get-iam-policy-' + random_suffix
    set_iam_policy_name = 'set-iam-policy-' + random_suffix
    rem_iam_policy_name = 'rem-iam-policy-' + random_suffix

    bindings = create_policy_bindings(member, roles)

    action_path = get_action_path(res_type)

    resources = [
        {
            # Get the existing policy.
            'name': get_iam_policy_name,
            'action': '{}.getIamPolicy'.format(action_path),
            'properties': {
                'resource': res_id
            }
        },
        {
            # Patch the existing policy.
            'name': set_iam_policy_name,
            'action': '{}.setIamPolicy'.format(action_path),
            'properties': {
                'resource': res_id,
                'policy': '$(ref.{})'.format(get_iam_policy_name),
                'gcpIamPolicyPatch': {'add': bindings}
            }
        },
    ]

    if CLEANUP_POLICY_ON_DELETE:
        resources.append({
            'name': rem_iam_policy_name,
            'action': '{}.setIamPolicy'.format(action_path),
            'metadata': {'runtimePolicy': ['DELETE']},
            'properties':
                {
                    'resource': res_id,
                    'policy': '$(ref.' + set_iam_policy_name + ')',
                    'gcpIamPolicyPatch': {'remove': copy.deepcopy(bindings)}
                }
        })

    return DMResource(None, resources, [])

def get_service_account(
        default_id,
        properties,
        project_roles,
        project_id,
        org_roles,
        org_id,
        sa_roles
):
    """ Creates a new service account. """

    account_id = properties.get('accountId', default_id)
    display_name = properties.get('displayName', account_id)
    sa_res_name = account_id
    sa_res = {
        'name': sa_res_name,
        'type': 'iam.v1.serviceAccount',
        'properties':
            {
                'accountId': account_id,
                'displayName': display_name,
                'projectId': project_id
            }
    }

    self_link = '$(ref.{}.email)'.format(sa_res_name)
    sa_name = 'serviceAccount:{}'.format(self_link)

    sa_bundle = DMResource(self_link, [sa_res], [])

    if project_roles:
        project_policy = patch_member_roles(
            sa_name,
            project_roles,
            'project',
            project_id
        )
        sa_bundle = merge_dm_resources(sa_bundle, project_policy)

    if org_roles:
        org_policy = patch_member_roles(
            sa_name,
            org_roles,
            'organization',
            'organizations/{}'.format(org_id)
        )
        sa_bundle = merge_dm_resources(sa_bundle, org_policy)

    if sa_roles:
        sa_policy = set_member_roles(
            sa_name,
            sa_roles,
            'serviceAccount',
            'projects/{}/serviceAccounts/{}'.format(project_id, self_link),
            project_id
        )
        sa_bundle = merge_dm_resources(sa_bundle, sa_policy)

    return sa_bundle

def get_client_service_account(project_id, properties):
    """ Creates a new service account for the client instance. """

    client_sa_settings = properties.get('serviceAccount', {})
    default_sa_id = 'forseti-client-gcp-' + get_random_string(SUFFIX_LENGTH)
    project_roles = [
        'roles/logging.logWriter',
        'roles/storage.objectViewer'
    ]

    return get_service_account(
        default_sa_id,
        client_sa_settings,
        project_roles,
        project_id,
        None,
        None,
        None
    )

def merge_dm_resources(first, *argv):
    """ Merges an arbitrary number of DM resources into one. """

    if argv:
        merged_resource = merge_dm_resources(argv[0], *argv[1:])
        return DMResource(
            first.self_link,
            first.resources + merged_resource.resources,
            first.outputs + merged_resource.outputs
        )

    return first

def get_server_service_account(project_id, properties, org_id):
    """ Creates a new service account for the server instance. """

    server_sa_settings = properties.get('serviceAccount', {})
    default_sa_id = 'forseti-server-gcp-' + get_random_string(SUFFIX_LENGTH)
    project_roles = [
        'roles/cloudsql.client',
        'roles/logging.logWriter',
        'roles/storage.objectViewer',
        'roles/storage.objectCreator'
    ]
    org_roles = [
        'roles/appengine.appViewer',
        'roles/bigquery.dataViewer',
        'roles/browser',
        'roles/cloudasset.viewer',
        'roles/cloudsql.viewer',
        'roles/compute.networkViewer',
        'roles/compute.securityAdmin',
        'roles/iam.securityReviewer',
        'roles/servicemanagement.quotaViewer',
        'roles/serviceusage.serviceUsageConsumer'
    ]
    sa_roles = [
        'roles/iam.serviceAccountTokenCreator'
    ]
    sa_bundle = get_service_account(
        default_sa_id,
        server_sa_settings,
        project_roles,
        project_id,
        org_roles,
        org_id,
        sa_roles
    )

    return sa_bundle

def keep_first(collection):
    """ Removes from the collection all elements except for the first one. """

    while len(collection) > 1:
        collection.pop(1)

def squash_patch_policies(set_policies, policy_ref):
    """
    Optimizes the policy assignments, so that we can achieve the same
    result with a smaller number of steps.
    """

    if set_policies:
        first_set_policy = set_policies[0]
        source_bindings = first_set_policy['properties']['gcpIamPolicyPatch']
        # Merge the remaining set- policies into the first one.
        for other in set_policies[1:]:
            other_bindings = other['properties']['gcpIamPolicyPatch']
            for action in other_bindings:
                if not action in source_bindings:
                    source_bindings[action] = other_bindings[action]
                else:
                    source_bindings[action] += other_bindings[action]

        # Update the reference to the get- policy.
        first_set_policy['properties']['policy'] = policy_ref
        return DMResource(None, [first_set_policy], [])

    return DMResource(None, [], [])

def group_iam_policies_by_targets(policies):
    """ Groups the collection of IAM actions by target. """

    policies_by_targets = {}
    for policy in policies:
        target = policy['properties']['resource']
        if not target in policies_by_targets:
            policies_by_targets[target] = []
        policies_by_targets[target].append(policy)

    return policies_by_targets

def is_get_policy(policy):
    """ Checks if the current policy action is a get action. """

    return 'getIamPolicy' in policy['action']

def is_set_policy(policy):
    """ Checks if the current policy action is a set action. """

    return 'setIamPolicy' in policy['action'] and not 'metadata' in policy

def is_rem_policy(policy):
    """ Checks if the current policy action is a patch-on-delete action. """

    return 'setIamPolicy' in policy['action'] and 'metadata' in policy

def optimize_policies_creation(first_sa_bundle, second_sa_bundle):
    """
    Reorganizes the policy patches so they don't affect the same resource
    more than once, thus avoiding the race condition.
    """

    all_resources = first_sa_bundle.resources + second_sa_bundle.resources
    policies = [policy for policy in all_resources if 'action' in policy]

    # Group policies by target (organization, project).
    policies_by_targets = group_iam_policies_by_targets(policies)

    # A placeholder for the result.
    extracted_policies = DMResource(self_link=None, resources=[], outputs=[])

    # Leave only one get- and one set- IamPolicy (add, remove) for each target.
    for _, target_policies in policies_by_targets.items():
        get_policy = next(
            (p for p in target_policies if is_get_policy(p)),
            None
        )
        if get_policy:
            extracted_policies.resources.append(get_policy)

            set_policies = [p for p in target_policies if is_set_policy(p)]
            set_policy_bundle = squash_patch_policies(
                set_policies,
                '$(ref.{})'.format(get_policy['name'])
            )
            extracted_policies = merge_dm_resources(
                extracted_policies,
                set_policy_bundle
            )

            rem_policies = [p for p in target_policies if is_rem_policy(p)]

            rem_policy_bundle = squash_patch_policies(
                rem_policies,
                '$(ref.{})'.format(set_policy_bundle.resources[0]['name'])
            )
            extracted_policies = merge_dm_resources(
                extracted_policies,
                rem_policy_bundle
            )
        else:
            extracted_policies.resources.extend(target_policies)

    keep_first(first_sa_bundle.resources)
    keep_first(second_sa_bundle.resources)

    return extracted_policies

def get_ref(res_name, prop='selfLink'):
    """ Gets a Deployment Manager reference link. """

    return '$(ref.{}.{})'.format(res_name, prop)

def get_server_bucket(properties, project_id, server_sa_email):
    """ Configures and gets a link to the server configuration bucket. """

    name = properties['name']

    bucket = DMResource(name, [], [])

    roles = set_member_roles(
        'serviceAccount:' + server_sa_email,
        ['roles/storage.objectAdmin'],
        'bucket',
        name,
        project_id
    )

    return merge_dm_resources(bucket, roles)

def find_output_value(name, outputs):
    """ Finds a specific output within a collection. """

    return next(
        output['value'] for output in outputs if output['name'] == name
    )

def get_cloud_sql(properties, project_id):
    """ Creates a Cloud SQL instance with a database. """

    instance_name = properties.get('instanceName', 'forseti-sql-' + project_id)
    # Add random suffix when creating/recreating instances.
    # GCP keeps the names for up to a week.
    # instance_name += '-' + get_random_string(SUFFIX_LENGTH)

    self_link = '$(ref.{}.name)'.format(instance_name)

    sql = {
        'name': instance_name,
        'type': 'gcp-types/sqladmin-v1beta4:instances',
        'properties': {
            'name': instance_name,
            'project': project_id,
            'backendType': 'SECOND_GEN',
            'databaseVersion': 'MYSQL_5_7',
            'region': properties['region'],
            'settings': {
                'tier': 'db-n1-standard-1',
                'backupConfiguration': {
                    'enabled': True,
                    'binaryLogEnabled': True
                },
                'activationPolicy': 'ALWAYS',
                'ipConfiguration': {
                    'ipv4Enabled': True,
                    'authorizedNetworks': [],
                    'requireSsl': True
                },
                'dataDiskSizeGb': '25',
                'dataDiskType': 'PD_SSD',
            },
            'instanceType': 'CLOUD_SQL_INSTANCE',
        }
    }

    db_name = instance_name + '-db'
    database = {
        'name': db_name,
        'type': 'gcp-types/sqladmin-v1beta4:databases',
        'properties':
            {
                'name': db_name,
                'project': project_id,
                'instance': self_link
            }
    }

    return DMResource(self_link, [sql, database], [
        {
            'name': 'connectionName',
            'value': get_ref(instance_name, 'connectionName')
        },
        {
            'name': 'databaseName',
            'value': '$(ref.{}.name)'.format(db_name)
        },
    ])

def get_firewall_rule(name, properties, project_id, network):
    """ Creates a firewall rule. """

    resource = {
        'name': name,
        'type': 'gcp-types/compute-v1:firewalls',
        'properties': copy.deepcopy(properties)
    }

    resource['properties']['project'] = project_id
    resource['properties']['network'] = network

    return DMResource(get_ref(name), [resource], [])

def get_firewall_rules(prefix, project_id, network):
    """ Creates firewall rules required by the Forseti network. """

    icmp_desc = {
        'sourceRanges': ['0.0.0.0/0'],
        'allowed': [{'IPProtocol': 'icmp'}]
    }
    icmp_rule = get_firewall_rule(
        prefix + '-allow-icmp',
        icmp_desc,
        project_id,
        network
    )

    ssh_desc = {
        'sourceRanges': ['0.0.0.0/0'],
        'allowed': [{'IPProtocol': 'tcp', 'ports': [22]}]
    }
    ssh_rule = get_firewall_rule(
        prefix + '-allow-ssh',
        ssh_desc,
        project_id,
        network
    )

    client_to_server_desc = {
        'sourceTags': ['forseti-client'],
        'targetTags': ['forseti-server'],
        'allowed': [{'IPProtocol': 'all'}]
    }

    cts_name = prefix + '-allow-client-to-server'
    cts_bundle = get_firewall_rule(
        cts_name,
        client_to_server_desc,
        project_id,
        network
    )

    return merge_dm_resources(icmp_rule, ssh_rule, cts_bundle)

def get_network(project_id):
    """ Creates a Forseti VPC. """

    name = 'forseti-network'

    network = {
        'name': name,
        'type': 'gcp-types/compute-v1:networks',
        'properties': {
            'name': name,
            'project': project_id,
            'autoCreateSubnetworks': True
        }
    }

    self_link = get_ref(name)

    firewall_bundle = get_firewall_rules(name, project_id, self_link)

    network_bundle = DMResource(self_link, [network], [
        {
            'name': 'networkName',
            'value': get_ref(name, 'name')
        },
        {
            'name': 'networkSelfLink',
            'value': self_link
        },
    ])

    return merge_dm_resources(network_bundle, firewall_bundle)

def get_server(properties, project, network, sql, service_account, bucket):
    """ Creates a Forseti server instance. """

    name = properties['name']

    instance_properties = copy.deepcopy(properties)
    instance = {
        'name': name,
        'type': 'server.py',
        'properties': instance_properties
    }

    if 'serviceAccount' in instance_properties:
        del instance_properties['serviceAccount']

    instance_properties['name'] = name
    instance_properties['tags'] = ['forseti-server']

    for sql_prop in ['databaseName', 'connectionName']:
        instance_properties[sql_prop] = find_output_value(
            sql_prop,
            sql.outputs
        )

    instance_properties['bucket'] = bucket.self_link
    instance_properties['project'] = project.self_link
    instance_properties['serviceAccountEmail'] = service_account.self_link
    instance_properties['serviceAccountScopes'] = [
        'https://www.googleapis.com/auth/cloud-platform'
    ]
    instance_properties['network'] = network.self_link

    self_link = get_ref(name)

    return DMResource(self_link, [instance], [
        {
            'name': 'serverName',
            'value': get_ref(name, 'name')
        },
        {
            'name': 'serverSelfLink',
            'value': self_link
        },
        {
            'name': 'serverInternalIp',
            'value': get_ref(name, 'internalIp')
        }
    ])

def get_client(properties, project, network, service_account, server):
    """ Creates a Forseti client instance. """

    name = properties['name']

    instance_properties = copy.deepcopy(properties)
    instance = {
        'name': name,
        'type': 'client.py',
        'properties': instance_properties
    }

    if 'serviceAccount' in instance_properties:
        del instance_properties['serviceAccount']

    instance_properties['name'] = name
    instance_properties['tags'] = ['forseti-client']

    instance_properties['serverIp'] = find_output_value(
        'serverInternalIp',
        server.outputs
    )
    instance_properties['project'] = project.self_link
    instance_properties['serviceAccountEmail'] = service_account.self_link
    instance_properties['serviceAccountScopes'] = [
        'https://www.googleapis.com/auth/cloud-platform'
    ]
    instance_properties['network'] = network.self_link

    self_link = get_ref(name)

    return DMResource(self_link, [instance], [
        {
            'name': 'clientName',
            'value': get_ref(name, 'name')
        },
        {
            'name': 'clientSelfLink',
            'value': self_link
        }
    ])

def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    org_id = properties['organizationId']
    project = get_forseti_project(context.env['name'], properties['project'])

    # Service accounts
    client_sa = get_client_service_account(
        project.self_link,
        properties['client']
    )

    server_sa = get_server_service_account(
        project.self_link,
        properties['server'],
        org_id
    )

    # Avoid race conditions at policy creation.
    policies = optimize_policies_creation(server_sa, client_sa)

    # Network + firewalls
    network = get_network(project.self_link)

    # Configure the server config bucket.
    bucket = get_server_bucket(
        properties['bucket'],
        project.self_link,
        server_sa.self_link
    )

    # Cloud SQL.
    cloud_sql = get_cloud_sql(properties['cloudSql'], project.self_link)

    if project.outputs: # creates a project
        wait_for_init_complete(
            project,
            client_sa,
            server_sa,
            network,
            cloud_sql,
            policies
        )

    # Client/Server instances.
    server = get_server(
        properties['server'],
        project,
        network,
        cloud_sql,
        server_sa,
        bucket
    )
    client = get_client(
        properties['client'],
        project,
        network,
        client_sa, server
    )

    # Join all resources into one final collection.
    result = merge_dm_resources(
        project,
        client_sa,
        server_sa,
        policies,
        bucket,
        cloud_sql,
        network,
        server,
        client
    )

    return {
        'resources': result.resources,
        'outputs': result.outputs
    }
