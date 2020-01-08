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
""" This template creates a Google Kubernetes Engine cluster. """

def generate_config(context):
    """ Entry point for the deployment resources. """

    resources = []
    outputs = []
    properties = context.properties
    name = properties['cluster'].get('name', context.env['name'])
    project_id = properties.get('project', context.env['project'])
    propc = properties['cluster']
    gke_cluster = {
        'name': context.env['name'],
        'type': '',
        'properties':
            {
                'parent': 'projects/{}/locations/{}'.format(
                    project_id,
                    properties.get('zone', properties.get('location', properties.get('region')))
                ),
                'cluster':
                    {
                        'name': name,
                    }
            }
    }

    if properties.get('zone'):
        # https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1beta1/projects.zones.clusters
        gke_cluster['type'] = 'gcp-types/container-v1beta1:projects.zones.clusters'
        # TODO: remove, this is a bug
        gke_cluster['properties']['zone'] = properties.get('zone')
    else:
        # https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1beta1/projects.locations.clusters
        gke_cluster['type'] = 'gcp-types/container-v1beta1:projects.locations.clusters'

    req_props = ['network', 'subnetwork']

    optional_props = [
        'initialNodeCount',
        'initialClusterVersion',
        'description',
        'nodeConfig',
        'nodePools',
        'privateClusterConfig',
        'binaryAuthorization',
        'binaryAuthorization',
        'networkConfig',
        'defaultMaxPodsConstraint',
        'resourceUsageExportConfig',
        'authenticatorGroupsConfig',
        'verticalPodAutoscaling',
        'tierSettings',
        'enableTpu',
        'databaseEncryption',
        'workloadIdentityConfig',
        'masterAuth',
        'loggingService',
        'monitoringService',
        'clusterIpv4Cidr',
        'addonsConfig',
        'locations',
        'enableKubernetesAlpha',
        'resourceLabels',
        'labelFingerprint',
        'legacyAbac',
        'networkPolicy',
        'ipAllocationPolicy',
        'masterAuthorizedNetworksConfig',
        'maintenancePolicy',
        'podSecurityPolicyConfig',
        'privateCluster',
        'masterIpv4CidrBlock'
    ]

    cluster_props = gke_cluster['properties']['cluster']

    for prop in req_props:
        cluster_props[prop] = propc.get(prop)
        if prop not in propc:
            raise KeyError(
                "{} is a required cluster property for a {} Cluster."
                .format(prop,
                        cluster_type)
            )

    for oprop in optional_props:
        if oprop in propc:
            cluster_props[oprop] = propc[oprop]

    resources.append(gke_cluster)

    # Output variables
    output_props = [
        'selfLink',
        'endpoint',
        'instanceGroupUrls',
        'clusterCaCertificate',
        'currentMasterVersion',
        'currentNodeVersion',
        'servicesIpv4Cidr'
    ]

    if (
        # https://github.com/GoogleCloudPlatform/deploymentmanager-samples/issues/463
        propc.get('enableDefaultAuthOutput', False) and (
            propc.get('masterAuth', {}).get('clientCertificateConfig', False)
        )
    ):
        output_props.append('clientCertificate')
        output_props.append('clientKey')

    for outprop in output_props:
        output_obj = {}
        output_obj['name'] = outprop
        ma_props = ['clusterCaCertificate', 'clientCertificate', 'clientKey']
        if outprop in ma_props:
            output_obj['value'] = '$(ref.' + context.env['name'] + \
                                  '.masterAuth.' + outprop + ')'
        elif outprop == 'instanceGroupUrls':
            output_obj['value'] = '$(ref.' + context.env['name'] + \
                '.nodePools[0].' + outprop + ')'
        else:
            output_obj['value'] = '$(ref.' + context.env['name'] + '.' + outprop + ')'

        outputs.append(output_obj)

    return {'resources': resources, 'outputs': outputs}
