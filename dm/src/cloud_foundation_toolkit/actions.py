# Copyr ight2018 Google Inc. All rights reserved.
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
""" Deployment Actions """

import glob
import json
import os.path
import sys

from apitools.base.py import exceptions as apitools_exceptions
from ruamel.yaml import YAML

from cloud_foundation_toolkit import LOG
from cloud_foundation_toolkit.deployment import Config, ConfigGraph, Deployment

# To avoid code repetition this ACTION_MAP is used to translate the
# args provided to the cmd line to the appropriate method of the
# deployment object
ACTION_MAP = {
    'apply': {
        'preview': 'preview'
    },
    'create': {
        'preview': 'preview'
    },
    'delete': {},
    'update': {
        'preview': 'preview'
    }
}


def check_file(config):
    extensions = ['.yaml', '.yml', '.jinja']
    for ext in extensions:
        if ext == config[-len(ext):]:
            return True


def get_config_files(config):
    """ Build a list of config files
    List could have files directory or yaml strings

    Args(list): List of configs. Each item can be a file, a directory,
        or a yaml string

    Returns: A list of config files or strings
    """

    config_files = []

    for conf in config:
        if os.path.isdir(conf):
            config_files.extend(
                [f for f in glob.glob(conf + '/*') if check_file(f)]
            )
        else:
            config_files.append(conf)

    LOG.debug('Config files %s', config_files)
    return config_files


def execute(args):
    action = args.action

    if action == 'delete' or (hasattr(args, 'reverse') and args.reverse):
        graph = reversed(
            ConfigGraph(get_config_files(args.config),
                        project=args.project)
        )
    else:
        graph = ConfigGraph(get_config_files(args.config), project=args.project)

    arguments = {}
    for k, v in vars(args).items():
        if k in ACTION_MAP.get(action, {}):
            arguments[ACTION_MAP[action][k]] = v

    LOG.debug(
        'Excuting %s on %s with arguments %s',
        action,
        args.config,
        arguments
    )

    if args.show_stages:
        output = []
        for level in graph:
            configs = []
            for config in level:
                configs.append(
                    {
                        'project': config.project,
                        'deployment': config.deployment,
                        'source': config.source
                    }
                )
            output.append(configs)
        if args.format == 'yaml':
            YAML().dump(output, sys.stdout)
        elif args.format == 'json':
            print(json.dumps(output, indent=2))
        else:
            for i, stage in enumerate(output, start=1):
                print('---------- Stage {} ----------'.format(i))
                for config in stage:
                    print(
                        ' - project: {}, deployment: {}, source: {}'.format(
                            config['project'],
                            config['deployment'],
                            config['source']
                        )
                    )
            print('------------------------------')

    else:
        for i, stage in enumerate(graph, start=1):
            print('---------- Stage {} ----------'.format(i))
            for config in stage:
                LOG.debug('%s config %s', action, config.deployment)
                deployment = Deployment(config)
                method = getattr(deployment, action)
                try:
                    method(**arguments)
                except apitools_exceptions.HttpNotFoundError:
                    LOG.warn('Deployment %s does not exit', config.deployment)
                    if action != 'delete':
                        raise
        print('------------------------------')
