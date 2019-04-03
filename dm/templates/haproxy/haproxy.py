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
""" This template creates a Compute Instance with an HAProxy installed
    and configured to load-balance traffic between instance groups.
"""

DISK_IMAGE = 'projects/debian-cloud/global/images/family/debian-9'
SETUP_HAPROXY_SH = """#!/bin/bash

set -euf -o pipefail

apt-get update && apt-get install -y haproxy

METADATA_SERVER="http://metadata.google.internal/computeMetadata/v1/instance/attributes"
function get_metadata() {
  curl -s "$METADATA_SERVER/$1" -H "Metadata-Flavor: Google"
}

# Set up an HAProxy update script.
CONFIG_UPDATER="/sbin/haproxy-conf-updater"
get_metadata "haproxy-conf-updater" > $CONFIG_UPDATER
REFRESH_RATE=`get_metadata "refresh-rate"`
chmod +x $CONFIG_UPDATER

# Set up an HAProxy config.
$CONFIG_UPDATER

# Keep the HAProxy config up to date.
CRONFILE=$(mktemp)
crontab -l > "${CRONFILE}" || true
echo "${REFRESH_RATE} * * * * ${CONFIG_UPDATER}" >> "${CRONFILE}"
crontab "${CRONFILE}"
service cron start
"""

HAPROXY_CONF_UPDATER_SH = """#!/bin/bash

set -euf -o pipefail

METADATA_SERVER="http://metadata.google.internal/computeMetadata/v1/instance/attributes"

function get_metadata() {
  curl -s "$METADATA_SERVER/$1" -H "Metadata-Flavor: Google"
}

CONFIG_FILE=/etc/haproxy/haproxy.cfg
BASE_CONFIG_FILE=${CONFIG_FILE}.bak

if [ ! -f $BASE_CONFIG_FILE ]; then
    cp $CONFIG_FILE $BASE_CONFIG_FILE
fi

TEMP_CONFIG_FILE=`mktemp`
LB_ALGORITHM=`get_metadata lb-algorithm`
IG_PORT=`get_metadata ig-port`
LB_PORT=`get_metadata lb-port`
LB_MODE=`get_metadata lb-mode`

# Build a server list.
SERVERS=
GCLOUD=`which gcloud`
for g in $(get_metadata groups); do
  if [[ "${g}" =~ zones/([^/]+)/instanceGroups/(.*)$ ]]; then
    GROUP="${BASH_REMATCH[2]}"
    ZONE="${BASH_REMATCH[1]}"
    SERVERS=${SERVERS}$'\\n'$($GCLOUD compute instance-groups list-instances \\
      $GROUP --zone $ZONE | grep -v NAME | \
      sed "s/^\\([^ ]*\\) .*\\$/  server \\1 \\1:$IG_PORT check/")
  else
    echo "Invalid group: ${g}"
  fi
done

# Set up the config file.
cp ${BASE_CONFIG_FILE} ${TEMP_CONFIG_FILE}

echo "
# Internal load balancing config

frontend tcp-in
	bind *:$LB_PORT
	mode $LB_MODE
	option ${LB_MODE}log
	default_backend instances

backend instances
	mode $LB_MODE
	balance $LB_ALGORITHM
${SERVERS}" >> ${TEMP_CONFIG_FILE}

# Update the config and restart if the config has changed.
ret=0
diff ${TEMP_CONFIG_FILE} ${CONFIG_FILE} || ret=$?
if [ ${ret} -ne 0 ]; then
  mv ${TEMP_CONFIG_FILE} ${CONFIG_FILE}
  service haproxy restart
fi
"""

def append_metadata_entry(metadata, new_key, new_value):
    """ Appends a new key-value pair to the existing metadata. """

    metadata['items'].append({
        'key': new_key,
        'value': new_value
    })

def configure_haproxy_frontend(properties, metadata):
    """ Sets up user-facing HAProxy parameters. """

    lb_properties = properties['loadBalancer']
    lb_algorithm = lb_properties['algorithm']
    lb_mode = lb_properties['mode']
    lb_port = lb_properties['port']

    append_metadata_entry(metadata, 'lb-algorithm', lb_algorithm)
    append_metadata_entry(metadata, 'lb-port', lb_port)
    append_metadata_entry(metadata, 'lb-mode', lb_mode)

def configure_haproxy_backend(home_zone, properties, metadata):
    """ Sets up instance-facing HAProxy parameters. """

    instances_properties = properties['instances']
    append_metadata_entry(metadata, 'ig-port', instances_properties['port'])
    groups = ' '.join(['zones/{}/instanceGroups/{}'.format(home_zone, group)
                       if 'zones/' not in group
                       else group
                       for group
                       in instances_properties['groups']])

    append_metadata_entry(metadata, 'groups', groups)
    cron_refresh_rate = instances_properties['refreshIntervalMin']
    cron_minutes_value = '*/' + str(cron_refresh_rate)
    append_metadata_entry(metadata, 'refresh-rate', cron_minutes_value)

def configure_haproxy_setup(metadata):
    """ Adds metadata that installs and configures the HAProxy. """

    append_metadata_entry(metadata, 'startup-script', SETUP_HAPROXY_SH)
    append_metadata_entry(metadata, 'haproxy-conf-updater',
                          HAPROXY_CONF_UPDATER_SH)


def generate_config(context):
    """ Entry point for the deployment resources. """

    properties = context.properties
    lb_name = properties.get('name', context.env['name'])
    zone = properties['zone']
    metadata = properties.get('metadata', {'items':[]})

    configure_haproxy_frontend(properties, metadata)
    configure_haproxy_backend(zone, properties, metadata)
    configure_haproxy_setup(metadata)

    service_account = properties['serviceAccountEmail']

    load_balancer = {
        'name': lb_name,
        'type': 'instance.py',
        'properties':
            {
                'machineType': properties['machineType'],
                'diskImage': DISK_IMAGE,
                'zone': zone,
                'network': properties['network'],
                'metadata': metadata,
                'serviceAccounts': [
                    {
                        'email': service_account,
                        'scopes': [
                            'https://www.googleapis.com/auth/compute.readonly'
                        ]
                    }
                ]
            }
    }

    return {
        'resources': [load_balancer],
        'outputs': [
            {
                'name': 'internalIp',
                'value': '$(ref.{}.internalIp)'.format(lb_name)
            },
            {
                'name': 'externalIp',
                'value': '$(ref.{}.externalIp)'.format(lb_name)
            },
            {
                'name': 'name',
                'value': '$(ref.{}.name)'.format(lb_name)
            },
            {
                'name': 'selfLink',
                'value': '$(ref.{}.selfLink)'.format(lb_name)
            },
            {
                'name': 'port',
                'value':  properties['loadBalancer']['port']
            }
        ]
    }
