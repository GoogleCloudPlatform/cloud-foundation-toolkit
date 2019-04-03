# Copyright 2017 The Forseti Security Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
""" Creates a GCE client instance for Forseti Security. """

FORSETI_HOME = '$USER_HOME/forseti-security'

FORSETI_CLIENT_CONFIG = '$FORSETI_HOME/forseti_conf_client.yaml'

EXPORT_VARS = ('export FORSETI_HOME={}\n'
               'export FORSETI_CLIENT_CONFIG={}\n'
              ).format(FORSETI_HOME,
                       FORSETI_CLIENT_CONFIG)


STARTUP_SCRIPT_TEMPLATE = """#!/bin/bash
exec > /tmp/deployment.log
exec 2>&1
# Ubuntu update.
sudo apt-get update -y
sudo apt-get upgrade -y
# Forseti setup.
sudo apt-get install -y git unzip
# Forseti dependencies
sudo apt-get install -y libffi-dev libssl-dev libmysqlclient-dev python-pip python-dev build-essential
USER=ubuntu
USER_HOME=/home/ubuntu
# Install fluentd if necessary.
FLUENTD=$(ls /usr/sbin/google-fluentd)
if [ -z "$FLUENTD" ]; then
      cd $USER_HOME
      curl -sSO https://dl.google.com/cloudagents/install-logging-agent.sh
      bash install-logging-agent.sh
fi
# Install Forseti Security.
cd $USER_HOME
rm -rf *forseti*
# Download Forseti source code
{download_forseti}
cd forseti-security
git fetch --all
{checkout_forseti_version}
# Forseti dependencies
pip install --upgrade pip==9.0.3
pip install -q --upgrade setuptools wheel
pip install -q --upgrade -r requirements.txt
# Install Forseti
python setup.py install
# Set ownership of the forseti project to $USER
chown -R $USER {forseti_home}
# Export variables
{persist_forseti_vars}
# Store the variables in /etc/profile.d/forseti_environment.sh 
# so all the users will have access to them
echo "echo '{persist_forseti_vars}' >> /etc/profile.d/forseti_environment.sh" | sudo sh
echo "server_ip: {server_ip}" > $FORSETI_CLIENT_CONFIG
chmod ugo+r $FORSETI_CLIENT_CONFIG
echo "Execution of startup script finished"
"""

def get_full_machine_type(project_id, zone, machine_type):
    """ Gets a full URL to the specified machine type. """

    prefix = 'https://www.googleapis.com/compute/v1'

    return '{}/projects/{}/zones/{}/machineTypes/{}'.format(
        prefix,
        project_id,
        zone,
        machine_type
    )

def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    instance_name = properties.get('name', context.env['name'])
    project_id = properties.get('project', context.env['project'])
    source_image = properties['sourceImage']
    source_path = properties['srcPath']
    version = properties['srcVersion']
    zone = properties['zone']
    machine_type = properties['machineType']
    machine_type_uri = get_full_machine_type(project_id, zone, machine_type)
    server_ip = properties['serverIp']
    startup_script = STARTUP_SCRIPT_TEMPLATE.format(
        download_forseti='git clone {}.git'.format(source_path),
        checkout_forseti_version='git checkout {}'.format(version),
        server_ip=server_ip,
        forseti_home=FORSETI_HOME,
        persist_forseti_vars=EXPORT_VARS,
    )

    resources = [
        {
            'name': instance_name,
            'type': 'gcp-types/compute-v1:instances',
            'properties':
                {
                    'project':
                        project_id,
                    'zone':
                        zone,
                    'machineType':
                        machine_type_uri,
                    'disks':
                        [
                            {
                                'deviceName': 'boot',
                                'type': 'PERSISTENT',
                                'boot': True,
                                'autoDelete': True,
                                'initializeParams':
                                    {
                                        'sourceImage': source_image,
                                    }
                            }
                        ],
                    'networkInterfaces':
                        [
                            {
                                'network':
                                    properties['network'],
                                'accessConfigs':
                                    [
                                        {
                                            'name': 'External NAT',
                                            'type': 'ONE_TO_ONE_NAT'
                                        }
                                    ]
                            }
                        ],
                    'serviceAccounts':
                        [
                            {
                                'email': properties['serviceAccountEmail'],
                                'scopes': properties['serviceAccountScopes'],
                            }
                        ],
                    'tags': {
                        'items': properties.get('tags',
                                                [])
                    },
                    'metadata':
                        {
                            'items':
                                [
                                    {
                                        'key': 'startup-script',
                                        'value': startup_script
                                    }
                                ]
                        }
                }
        }
    ]

    outputs = [
        {
            'name': 'name',
            'value': '$(ref.{}.name)'.format(instance_name)
        },
        {
            'name': 'selfLink',
            'value': '$(ref.{}.selfLink)'.format(instance_name)
        }
    ]

    return {'resources': resources, 'outputs': outputs}
