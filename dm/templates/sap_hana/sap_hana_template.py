# Copyright 2020 Google Inc. All rights reserved.
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
""" This template deploys SAP HANA and all required infrastructure resources (network, firewall rules, NAT, etc). """

def generate_config(context):
    
    properties = context.properties
    # Creating the network (VPC + subnets) resource
    network = {
        'name': 'sap-poc-vpc',
        'type': 'network.py',
        'properties': {
            'autoCreateSubnetworks': False,
            'subnetworks': [{
                'name': 'subnetwork-1',
                'region': properties['region'],
                'ipCidrRange': '10.0.0.0/24',

            }, {
                'name': 'subnetwork-2',
                'region': properties['region'],
                'ipCidrRange': '192.168.0.0/24',
            }]
        }
    }
   
    #Create a Cloud NAT Gateway
    cloud_router = {
        'name': 'cloud-nat-gateway',
        'type': 'cloud_router.py',
        'properties': {
            'name': 'cloud-nat-router',
            'network': '$(ref.sap-poc-vpc.name)',
            'region': properties['region'],
            'nats': [{
                'name': 'cloud-nat',
                'sourceSubnetworkIpRangesToNat': 'LIST_OF_SUBNETWORKS',
                'natIpAllocateOption': 'AUTO_ONLY',
                'subnetworks': [{
                    'name': '$(ref.subnetwork-1.selfLink)'
                }]
            }]
        }
    }

    #Create a windows bastion host to be used for installing HANA Studio and connecting to HANA DB
    windows_bastion_host = {
        'name': 'windows-bastion-host',
        'type': 'instance.py',
        'properties': {
            'zone': properties['primaryZone'],
            'diskImage': 'projects/windows-cloud/global/images/family/windows-2019',
            'machineType': 'n1-standard-1',
            'diskType': 'pd-ssd',
            'networks': [{
                'network': "$(ref.sap-poc-vpc.selfLink)",
                'subnetwork': "$(ref.subnetwork-2.selfLink)",
                'accessConfigs': [{
                    'type': 'ONE_TO_ONE_NAT'
                }]
            }],
            'tags': {
                'items': ['jumpserver']
            }
        }
    }

    #Create a linux bastion host which will be used to connect to HANA DB CLI and run commands.
    linux_bastion_host = {
        'name': 'linux-bastion-host',
        'type': 'instance.py',
        'properties': {
            'zone': properties['primaryZone'],
            'diskImage': 'projects/ubuntu-os-cloud/global/images/family/ubuntu-1804-lts',
            'machineType': 'f1-micro',
            'diskType': 'pd-ssd',
            'networks': [{
                'network': '$(ref.sap-poc-vpc.selfLink)',
                'subnetwork': '$(ref.subnetwork-2.selfLink)',
                'accessConfigs': [{
                    'type': 'ONE_TO_ONE_NAT'
                }]
            }],
            'tags': {
                'items': ['jumpserver']
            }
        }
    }
    # Create necessary Firewall rules to allow connectivity to HANA DB from both bastion hosts.
    firewall_rules = {
        'name': 'firewall-rules',
        'type': 'firewall.py',
        'properties': {
            'network': '$(ref.sap-poc-vpc.selfLink)',
            'rules': [{
                'name': 'allow-ssh-and-rdp',
                'allowed':[ {
                    'IPProtocol': 'tcp',
                    'ports': ['22', '3389']
                }],
                'description': 'Allow SSH and RDP from outside to bastion host.',
                'direction': 'INGRESS',
                'targetTags': ["jumpserver"],
                'sourceRanges': [
                    '0.0.0.0/0'
                ]
              }, {
                'name': 'jumpserver-to-hana',
                'allowed': [{
                    'IPProtocol': 'tcp',
                    'ports': ['22', '30015'] # In general,the port to open is 3 <Instance number> 15 to allow HANA Studio to Connect to HANA DB < -- -- -TODO
                }],
                'description': 'Allow SSH from bastion host to HANA instances',
                'direction': 'INGRESS',
                'targetTags': ["hana-db"],
                'sourceRanges': ['$(ref.subnetwork-2.ipCidrRange)']
              }
            ]
        }
    }

    sap_hana_resource = {}
    if properties.get('secondaryZone'):  # HANA HA deployment
        sap_hana_resource = {
            'name': 'sap_hana',
            "type": 'sap_hana_ha.py',
            'properties': {
                'primaryInstanceName': properties['primaryInstanceName'],
                'secondaryInstanceName': properties['secondaryInstanceName'],
                'primaryZone': properties['primaryZone'],
                'secondaryZone': properties['secondaryZone'],
                'sap_vip': '10.1.0.10', #TO DO: improve this by reserving an internal IP address in advance.
            }
        }
        
    else:
        sap_hana_resource = { #HANA standalone setup
            'name': 'sap_hana',
            "type": 'sap_hana.py',
            'properties': {
                'instanceName': properties['primaryInstanceName'],
                'zone': properties['primaryZone'],
                'sap_hana_scaleout_nodes': 0
            }
        }
    
    #Add the rest of "common & manadatory" properties
    sap_hana_resource['properties']['dependsOn'] = ['$(ref.cloud-nat-gateway.selfLink)']
    sap_hana_resource['properties']['instanceType'] = properties['instanceType']
    sap_hana_resource['properties']['subnetwork'] = 'subnetwork-1'
    sap_hana_resource['properties']['linuxImage'] = properties['linuxImage']
    sap_hana_resource['properties']['linuxImageProject'] = properties['linuxImageProject']
    sap_hana_resource['properties']['sap_hana_deployment_bucket'] = properties['sap_hana_deployment_bucket']
    sap_hana_resource['properties']['sap_hana_sid'] = properties['sap_hana_sid']
    sap_hana_resource['properties']['sap_hana_instance_number'] = 00
    sap_hana_resource['properties']['sap_hana_sidadm_password'] = properties['sap_hana_sidadm_password']
    sap_hana_resource['properties']['sap_hana_system_password'] = properties['sap_hana_system_password']
    sap_hana_resource['properties']['networkTag'] = 'hana-db'
    sap_hana_resource['properties']['publicIP'] = False

    
    # Define optional properties.
    optional_properties = [
        'serviceAccount',
        'sap_hana_backup_size',
        'sap_hana_double_volume_size',
        'sap_hana_sidadm_uid',
        'sap_hana_sapsys_gid',
        'sap_deployment_debug',
        'post_deployment_script'
    ]
    # Add optional properties if there are any.
    for prop in optional_properties:
        append_optional_property(sap_hana_resource, properties, prop)


    resources = [network, cloud_router, windows_bastion_host, linux_bastion_host, firewall_rules, sap_hana_resource]
    
    return { 'resources': resources}


def append_optional_property(resource, properties, prop_name):
    """ If the property is set, it is added to the resource. """

    val = properties.get(prop_name)
    if val:
        resource['properties'][prop_name] = val
    return