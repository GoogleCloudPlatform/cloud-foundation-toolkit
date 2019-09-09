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
""" Creates a Cloud KMS KeyRing and cryptographic key resources. """


def generate_config(context):
    """
    Entry point for the deployment resources
    """

    resources = []
    properties = context.properties
    project_id = properties.get('project', context.env['project'])
    parent = 'projects/{}/locations/{}'.format(
        project_id,
        properties.get('region')
    )
    keyring_name = properties.get('keyRingName', context.env['name'])
    keyring_id = '{}/keyRings/{}'.format(parent, keyring_name)
    # https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations.keyRings
    provider = 'gcp-types/cloudkms-v1:projects.locations.keyRings'
    # keyring resource
    keyring = {
        'name': context.env['name'],
        'type': provider,
        'properties': {
            'parent': parent,
            'keyRingId': keyring_name
        }
    }
    resources.append(keyring)

    # cryptographic key resources
    for key in properties.get('keys', []):
        key_name = key['cryptoKeyName'].lower()
        key_resource = '{}-{}'.format(context.env['name'], key_name)
        crypto_key = {
            'name': key_resource,
            'type': provider + '.cryptoKeys',
            'properties':
                {
                    'parent': keyring_id,
                    'cryptoKeyId': key_name,
                    'purpose': key.get('cryptoKeyPurpose'),
                    'labels': key.get('labels',
                                      {})
                },
            'metadata': {
                'dependsOn': [context.env['name']]
            }
        }

        # crypto key optional properties
        for prop in ['versionTemplate', 'nextRotationTime', 'rotationPeriod']:
            if prop in key:
                crypto_key['properties'][prop] = key.get(prop)
        resources.append(crypto_key)

        # IAM policy bindings for the crypto key
        if 'iamPolicyBinding' in key:
            provider = 'gcp-types/cloudkms-v1:cloudkms.projects.locations.keyRings'
            key_resource_name = '{}/cryptoKeys/{}'.format(keyring_id, key_name)
            crypto_key_iam = {
                'name': '{}-iamPolicy'.format(key_resource),
                # https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations.keyRings.cryptoKeys/setIamPolicy
                # https://cloudkms.googleapis.com/$discovery/rest?version=v1
                'action': provider + '.cryptoKeys.setIamPolicy',
                'properties':
                    {
                        'resource': key_resource_name,
                        'policy': {
                            'bindings': key.get('iamPolicyBinding')
                        }
                    },
                'metadata': {
                    'dependsOn': [key_resource]
                }
            }
            resources.append(crypto_key_iam)

    return {
        'resources':
            resources,
        'outputs':
            [
                {
                    'name': 'keyRing',
                    'value': '$(ref.{}.name)'.format(context.env['name'])
                }
            ]
    }
