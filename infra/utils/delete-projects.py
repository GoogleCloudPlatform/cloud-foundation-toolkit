#!/usr/bin/env python3

# Copyright 2018 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import google.api_core
from google.cloud import resource_manager
from googleapiclient.discovery import build
import sys
from pprint import pprint
import argparse
from googleapiclient import discovery
from oauth2client.client import GoogleCredentials

client = resource_manager.Client()
credentials = GoogleCredentials.get_application_default()
client2 = build('cloudresourcemanager', 'v2', credentials=credentials)

def delete_liens(project_id):
    service = discovery.build('cloudresourcemanager', 'v1', credentials=credentials)

    parent = 'projects/{}'.format(project_id)
    request = service.liens().list(parent=parent)
    response = request.execute()

    liens_deleted = 0

    for lien in response.get('liens', []):
        print("Deleting lien:", lien)
        d_request = service.liens().delete(name=lien.get('name'))
        d_request.execute()
        liens_deleted += 1

    return liens_deleted

def delete_project(project):
    try:
        project.delete()
    except google.api_core.exceptions.BadRequest as e:
        liens_deleted = delete_liens(project.project_id)
        if liens_deleted >= 1:
            delete_project(project)
        else:
            print("Bad request and no liens found.")
            print(e)
    except (google.api_core.exceptions.Forbidden) as e:
        print("Failed to delete {}".format(project.project_id))
        print(e)

def delete_children(parent_type, parent_id):
    print("Deleting children of {} {}".format(parent_type, parent_id))

    project_filter = {
        'parent.type': parent_type,
        'parent.id': parent_id
    }
    for project in client.list_projects(project_filter):
        if (project.status != "ACTIVE"):
            print("  Skipping deletion of inactive project {}...".format(project.project_id))
            continue
        print("  Deleting project {} (status={})...".format(project.project_id, project.status))
        delete_project(project)

    name = "{}s/{}".format(parent_type, parent_id)
    res = client2.folders().list(parent=name).execute()
    for folder in res.get('folders', []):
        delete_children("folder", folder.get('name').split('/')[-1])

    deletion = client2.folders().delete(name=name).execute()
    if deletion.get('lifecycleState') == 'DELETE_REQUESTED':
        print("Deleted {}".format(name))
    else:
        print(deletion)

def main(argv):
    parser = argparser()
    args = parser.parse_args(argv[1:])

    (parent_type, parent_id) = args.parent_id.split('/')
    
    delete_children(parent_type.strip('s'), parent_id)

def argparser():
    parser = argparse.ArgumentParser(description='Delete projects within a folder')
    parser.add_argument('parent_id')
    return parser


if __name__ == "__main__":
    main(sys.argv)
