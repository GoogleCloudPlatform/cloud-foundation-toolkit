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
""" Uploads a Cloud Function from the local source to a GS bucket. """

import zipfile
import hashlib
import base64
import uuid
from StringIO import StringIO
from copy import deepcopy

GS_SCHEMA_LENGTH = 5

def extract_source_files(imports, local_upload_path):
    """ Returns tuples of the imported sources files. """

    imported_files = []
    for imported_file in imports:
        if imported_file.startswith(local_upload_path):
            file_name = imported_file[len(local_upload_path):]
            file_content = imports[imported_file]
            imported_files.append((file_name, file_content))

    return imported_files

def archive_files(files):
    """ Archives input files and returns the result as a binary array. """

    output_file = StringIO()
    sources_zip = zipfile.ZipFile(output_file,
                                  mode='w',
                                  compression=zipfile.ZIP_DEFLATED)

    for source_file in files:
        sources_zip.writestr(*source_file)

    sources_zip.close()
    return output_file.getvalue()

def upload_source(function, imports, local_path, source_archive_url):
    """ Uploads the Cloud Function source code from the local machine 
    to a Cloud Storage bucket. If the bucket does not exist, creates it.
    """

    # Creates an in-memory archive of the Cloud Function source files.
    sources = extract_source_files(imports, local_path)
    archive_base64 = base64.b64encode(archive_files(sources))

    # The Cloud Function knows it was updated when MD5 changes.
    md5 = hashlib.md5()
    md5.update(archive_base64)

    # Splits the upload path into the bucket and archive names.
    bucket_name = source_archive_url[:source_archive_url.index('/', GS_SCHEMA_LENGTH)] # pylint: disable=line-too-long
    archive_name = source_archive_url[source_archive_url.rfind('/') + 1:]

    # Uses a Docker volume to pass the archive between the build steps.
    volume = '/cloud-function'
    volume_archive_path = volume + '/' + archive_name
    volumes = [
        {
            'name': 'cloud-function',
            'path': volume
        }
    ]

    # Saves the inline base64-ZIP to a file.
    cmd = "echo '{}' | base64 -d > {};".format(archive_base64,
                                               volume_archive_path)

    build_action = {
        'name': 'upload-task',
        'action': 'gcp-types/cloudbuild-v1:cloudbuild.projects.builds.create',
        'metadata':
            {
                'runtimePolicy': ['UPDATE_ON_CHANGE']
            },
        'properties':
            {
                'steps':
                    [
                        { # Saves a ZIP to a file.
                            'name': 'ubuntu',
                            'args': ['bash', '-c', cmd],
                            'volumes': volumes,
                        },
                        { # Creates a bucket if one does not exist.
                            'name': 'gcr.io/cloud-builders/gsutil',
                            'args': [
                                '-c',
                                'gsutil mb {} || true'.format(bucket_name)
                            ],
                            'entrypoint': '/bin/bash'
                        },
                        { # Uploads the ZIP to the bucket.
                            'name': 'gcr.io/cloud-builders/gsutil',
                            'args': [
                                'cp',
                                volume_archive_path, source_archive_url
                            ],
                            'volumes': deepcopy(volumes)
                        }
                    ],
                'timeout': '120s'
            }
    }

    function['properties']['labels'] = {'content-md5': md5.hexdigest()}

    return ([build_action], [])

def generate_bucket_name():
    """ Generates a bucket name for the Cloud Function. """

    return 'gs://cloud-functions-{}'.format(uuid.uuid4())

def generate_archive_name():
    """ Generates the Cloud Function's zip name. """

    return 'cloud-function-{}.zip'.format(uuid.uuid4())

def generate_upload_path():
    """ Generates the full upload path for the Cloud Function. """

    return generate_bucket_name() + '/' + generate_archive_name()
