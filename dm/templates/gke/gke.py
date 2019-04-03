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
    project_id = context.env['project']
    properties = context.properties
    cluster_type = properties.get('clusterLocationType')
    propc = properties['cluster']
    name = propc.get('name') or context.env['name']
    gke_cluster = {
        'name': name,
        'type': '',
        'properties':
            {
                'cluster':
                    {
                        'name':
                            name + '-cluster',
                        'initialNodeCount':
                            propc.get('initialNodeCount'),
                        'initialClusterVersion':
                            propc.get('initialClusterVersion')
                    }
            }
    }

    if cluster_type == 'Regional':
        provider = 'gcp-types/container-v1beta1:projects.locations.clusters'
        if not properties.get('region'):
            raise KeyError(
                "region is a required property for a {} Cluster."
                .format(cluster_type)
            )
        parent = 'projects/{}/locations/{}'.format(
            project_id,
            properties.get('region')
        )
        gke_cluster['properties']['parent'] = parent

    elif cluster_type == 'Zonal':
        provider = 'container.v1.cluster'
        if not properties.get('zone'):
            raise KeyError(
                "zone is a required property for a {} Cluster."
                .format(cluster_type)
            )
        gke_cluster['properties']['zone'] = properties.get('zone')

    gke_cluster['type'] = provider

    req_props = ['network', 'subnetwork']

    optional_props = [
        'description',
        'nodeConfig',
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
        'masterAuthorizedNetworksConfig'
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
        'clientCertificate',
        'clientKey',
        'currentMasterVersion',
        'currentNodeVersion',
        'nodeIpv4CidrSize',
        'servicesIpv4Cidr'
    ]

    for outprop in output_props:
        output_obj = {}
        output_obj['name'] = outprop
        ma_props = ['clusterCaCertificate', 'clientCertificate', 'clientKey']
        if outprop in ma_props:
            output_obj['value'] = '$(ref.' + name + \
                '.masterAuth.' + outprop + ')'
        elif outprop == 'instanceGroupUrls':
            output_obj['value'] = '$(ref.' + name + \
                '.nodePools[0].' + outprop + ')'
        else:
            output_obj['value'] = '$(ref.' + name + '.' + outprop + ')'

        outputs.append(output_obj)

    return {'resources': resources, 'outputs': outputs}
