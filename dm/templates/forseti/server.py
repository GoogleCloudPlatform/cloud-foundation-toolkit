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
""" Creates a GCE server instance for Forseti Security. """

FORSETI_HOME = '$USER_HOME/forseti-security'

FORSETI_SERVER_CONF = '{}/configs/forseti_conf_server.yaml'.format(FORSETI_HOME)

EXPORT_FORSETI_VARS = (
    'export FORSETI_HOME={}\n'
    'export FORSETI_SERVER_CONF={}\n'
).format(FORSETI_HOME,
         FORSETI_SERVER_CONF)

STARTUP_SCRIPT_TEMPLATE = """#!/bin/bash
exec > /tmp/deployment.log
exec 2>&1
# Ubuntu update.
sudo apt-get update -y
sudo apt-get upgrade -y
sudo apt-get update && sudo apt-get --assume-yes install google-cloud-sdk
USER_HOME=/home/ubuntu
# Install fluentd if necessary.
FLUENTD=$(ls /usr/sbin/google-fluentd)
if [ -z "$FLUENTD" ]; then
      cd $USER_HOME
      curl -sSO https://dl.google.com/cloudagents/install-logging-agent.sh
      bash install-logging-agent.sh
fi
# Check whether Cloud SQL proxy is installed.
CLOUD_SQL_PROXY=$(which cloud_sql_proxy)
if [ -z "$CLOUD_SQL_PROXY" ]; then
        cd $USER_HOME
        wget https://dl.google.com/cloudsql/cloud_sql_proxy.{cloudsql_arch}
        sudo mv cloud_sql_proxy.{cloudsql_arch} /usr/local/bin/cloud_sql_proxy
        chmod +x /usr/local/bin/cloud_sql_proxy
fi
# Install Forseti Security.
cd $USER_HOME
rm -rf *forseti*
# Download Forseti source code
{download_forseti}
cd forseti-security
git fetch --all
{checkout_forseti_version}
# Forseti Host Setup
sudo apt-get install -y git unzip
# Forseti host dependencies
sudo apt-get install -y $(cat install/dependencies/apt_packages.txt | grep -v "#" | xargs)
# Forseti dependencies
pip install --upgrade pip==9.0.3
pip install -q --upgrade setuptools wheel
pip install -q --upgrade -r requirements.txt
# Setup Forseti logging
touch /var/log/forseti.log
chown ubuntu:root /var/log/forseti.log
cp {forseti_home}/configs/logging/fluentd/forseti.conf /etc/google-fluentd/config.d/forseti.conf
cp {forseti_home}/configs/logging/logrotate/forseti /etc/logrotate.d/forseti
chmod 644 /etc/logrotate.d/forseti
service google-fluentd restart
logrotate /etc/logrotate.conf
# Change the access level of configs/ rules/ and run_forseti.sh
chmod -R ug+rwx {forseti_home}/configs {forseti_home}/rules {forseti_home}/install/gcp/scripts/run_forseti.sh
# Install Forseti
python setup.py install
# Export variables required by initialize_forseti_services.sh.
{export_initialize_vars}
# Export variables required by run_forseti.sh
{export_forseti_vars}
# Store the variables in /etc/profile.d/forseti_environment.sh 
# so all the users will have access to them
echo "echo '{export_forseti_vars}' >> /etc/profile.d/forseti_environment.sh" | sudo sh
# Download server configuration from GCS
gsutil cp gs://{scanner_bucket}/configs/forseti_conf_server.yaml {forseti_server_conf}
gsutil cp -r gs://{scanner_bucket}/rules {forseti_home}/
# Start Forseti service depends on vars defined above.
bash ./install/gcp/scripts/initialize_forseti_services.sh
echo "Starting services."
systemctl start cloudsqlproxy
sleep 5
systemctl start forseti
echo "Success! The Forseti API server has been started."
# Create a Forseti env script
FORSETI_ENV="$(cat <<EOF
#!/bin/bash
export PATH=$PATH:/usr/local/bin
# Forseti environment variables
export FORSETI_HOME=/home/ubuntu/forseti-security
export FORSETI_SERVER_CONF=$FORSETI_HOME/configs/forseti_conf_server.yaml
export SCANNER_BUCKET={scanner_bucket}
EOF
)"
echo "$FORSETI_ENV" > $USER_HOME/forseti_env.sh
USER=ubuntu
# Use flock to prevent rerun of the same cron job when the previous job is still running.
# If the lock file does not exist under the tmp directory, it will create the file and put a lock on top of the file.
# When the previous cron job is not finished and the new one is trying to run, it will attempt to acquire the lock
# to the lock file and fail because the file is already locked by the previous process.
# The -n flag in flock will fail the process right away when the process is not able to acquire the lock so we won't
# queue up the jobs.
# If the cron job failed the acquire lock on the process, it will log a warning message to syslog.
(echo "{run_frequency} (/usr/bin/flock -n /home/ubuntu/forseti-security/forseti_cron_runner.lock $FORSETI_HOME/install/gcp/scripts/run_forseti.sh || echo '[forseti-security] Warning: New Forseti cron job will not be started, because previous Forseti job is still running.') 2>&1 | logger") | crontab -u $USER -
echo "Added the run_forseti.sh to crontab under user $USER"
echo "Execution of startup script finished"
"""

def get_export_initialize_vars(database, port, connection_string):
    """ Gets the shell script that persists the Forseti env variables. """

    template = """
        export SQL_PORT={}\n
        export SQL_INSTANCE_CONN_STRING="{}"\n
        export FORSETI_DB_NAME="{}"\n
    """
    return template.format(port, connection_string, database)


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
    """ Generates a configuration. """

    properties = context.properties
    instance_name = properties.get('name', context.env['name'])
    project_id = properties.get('project', context.env['project'])
    source_image = properties['sourceImage']
    source_path = properties['srcPath']
    version = properties['srcVersion']
    zone = properties['zone']
    bucket = properties['bucket']
    run_frequency = properties['frequency']
    machine_type = properties['machineType']
    machine_type_uri = get_full_machine_type(project_id, zone, machine_type)

    startup_script = STARTUP_SCRIPT_TEMPLATE.format(
        cloudsql_arch=properties['sqlOsArch'],
        download_forseti='git clone {}.git'.format(source_path),
        checkout_forseti_version='git checkout {}'.format(version),
        forseti_home=FORSETI_HOME,
        scanner_bucket=bucket,
        forseti_server_conf=FORSETI_SERVER_CONF,
        export_initialize_vars=get_export_initialize_vars(
            properties['databaseName'],
            properties['port'],
            properties['connectionName']
        ),
        export_forseti_vars=EXPORT_FORSETI_VARS,
        run_frequency=run_frequency,
    )

    resources = []

    resources.append(
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
                                        'sourceImage': source_image
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
    )
    return {
        'resources':
            resources,
        'outputs':
            [
                {
                    'name': 'name',
                    'value': '$(ref.{}.name)'.format(instance_name)
                },
                {
                    'name': 'selfLink',
                    'value': '$(ref.{}.selfLink)'.format(instance_name)
                },
                {
                    'name':
                        'internalIp',
                    'value':
                        '$(ref.{}.networkInterfaces[0].networkIP)'.
                        format(instance_name)
                }
            ]
    }
