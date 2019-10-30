#!/usr/bin/env python

# Copyright 2019 Google LLC
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
# Verifies that all source files contain the necessary copyright boilerplate
# snippet.

# Generates a dynamic bash script with a bunch of TF_VAR_* exports
# with names and values obtained from the provisioned terraform
# configuration at `--path` (defaults to current working directory).
#
# To omit intermediate file creation, the script is just being printed
# to stdout. You can source it in bash with the following command:
#
#   $ source <(python export_tf_outputs.py --path=test/setup)
#
import json
import os
import argparse
import sys


def get_args():
    """Parses command line arguments.

    Configures and runs argparse.ArgumentParser to extract command line
    arguments.

    Returns:
        An argparse.Namespace containing the arguments parsed from the
        command line
    """
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--path",
        default=os.path.abspath(os.getcwd()),
        help="path to the terraform configuration to get outputs from")
    return parser.parse_args()


def get_service_account(sa_key):
    """Decode service account from base64"""
    return os.popen("echo %s | base64 --decode" % sa_key).read().strip()


def main(args):
    """Utility to pipe outputs from one terraform configuration to another.

    Parses outputs of applied terraform configuration and prints them as
    TF_VAR_{name}={value} to be later supplied to another terraform
    configuration
    """

    # If path to terraform doesn't exist, notify and gracefully exit
    if not os.path.exists(args.path):
        script = "#!/usr/bin/env bash\n"
        script += ("echo 'Warning: Folder not found: %s"
                  ) % args.path
        print(script)
        sys.exit(0)

    # Get terraform outputs and parse them from json.
    terraform_command = "cd %s && terraform output -json"
    outputs_json = os.popen(terraform_command % args.path).read()
    outputs = json.loads(outputs_json)

    # If the specified folder was not a terraform configuration or if it was not
    # provisioned, terraform just silently ignores it and returns an empty object.
    # Notify user about this and exit.
    if not outputs:
        script = "#!/usr/bin/env bash\n"
        script += ("echo 'Warning: The terraform state file at %s either has no "
                   "outputs defined, or the terraform configuration has not"
                   "been provisioned yet.'"
                  ) % args.path
        print(script)
        sys.exit(0)

    # process value of each variable
    variables = {}
    plain_types = (int, float, "".__class__, u"".__class__)
    for name in outputs:
        if isinstance(outputs[name]['value'], plain_types):
            variables[name] = outputs[name]['value']
        else:
            variables[name] = json.dumps(outputs[name]['value'])

    # Generate a bash script which exports TF_VAR`s
    script = "#!/usr/bin/env bash\n"
    script += "echo 'Automatically setting inputs from outputs of %s'\n" % args.path
    for name in variables:
        script += "export TF_VAR_%s='%s'\n" % (name, variables[name])

    # Handle the magic `sa_key` variable in a special way
    if 'sa_key' in variables:
        service_account = get_service_account(variables['sa_key'])
        script += "export SERVICE_ACCOUNT_JSON='%s'" % service_account

    # Output to stdout
    print(script)


if __name__ == "__main__":
    ARGS = get_args()
    main(ARGS)
