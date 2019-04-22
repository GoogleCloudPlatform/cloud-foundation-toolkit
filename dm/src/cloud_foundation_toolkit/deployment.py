import base64
from collections import namedtuple
import io
import os
import os.path
import re
from six.moves import input
import sys
import tempfile

from apitools.base.py import exceptions as apitools_exceptions
from googlecloudsdk.api_lib.deployment_manager import dm_api_util
from googlecloudsdk.api_lib.deployment_manager import dm_base
from googlecloudsdk.api_lib.deployment_manager import exceptions as dm_exceptions
from googlecloudsdk.command_lib.deployment_manager import dm_util
from googlecloudsdk.command_lib.deployment_manager import dm_write
from googlecloudsdk.command_lib.deployment_manager import flags
from googlecloudsdk.command_lib.deployment_manager.importer import BuildConfig
from googlecloudsdk.command_lib.deployment_manager.importer import BuildTargetConfig
from googlecloudsdk.core.resource import resource_printer
from googlecloudsdk.third_party.apis.deploymentmanager.v2 import deploymentmanager_v2_messages as messages
import jinja2
import networkx as nx

from cloud_foundation_toolkit import LOG
from cloud_foundation_toolkit.dm_utils import DM_API
from cloud_foundation_toolkit.dm_utils import DM_OUTPUT_QUERY_REGEX
from cloud_foundation_toolkit.dm_utils import DMOutputQueryAttributes
from cloud_foundation_toolkit.dm_utils import get_deployment
from cloud_foundation_toolkit.dm_utils import get_deployment_output
from cloud_foundation_toolkit.dm_utils import parse_dm_output_url
from cloud_foundation_toolkit.dm_utils import parse_dm_output_token
from cloud_foundation_toolkit.yaml_utils import CFTBaseYAML

Node = namedtuple('Node', ['project', 'deployment'])


def ask():
    """Function that asks for user input from stdin."""
    answer = input("Update(u), Skip (s), or Abort(a) Deployment? ")
    while answer not in ['u', 's', 'a']:
        answer = input("Update(u), Skip (s), or Abort(a) Deployment? ")
    return answer


class Config(object):
    """Class representing a CFT config.

    Attributes:
        as_file (io.StringIO): A file-like interface to the
            jinja-rendered config.
        as_string (string): the jinja-rendered config.
        id (string): A base64-encoded id representing the path or raw
            content of the config. Could be used as a dict key.
        source (string): The path or the raw content of config (obtained
            by base64-decoding the 'id' attribute
    """
    yaml = CFTBaseYAML()

    def __init__(self, item, project=None):
        """ Contructor """

        self.source = item
        if project:
            self._project = project

        if os.path.exists(item):
            with io.open(item) as _fd:
                self.as_string = jinja2.Template(_fd.read()
                                                ).render(env=os.environ)
        else:
            self.as_string = jinja2.Template(item).render(env=os.environ)

        # YAML gets parsed twice:
        # 1. Here, to figure out deployment name, project and dependency list.
        # 2. When the Deployment() obj gets instantiated (to get the value of
        #    the output from the DM API)
        # This approach takes more CPU, but it's less error prone than
        # scanning the file ourselves.
        self.as_dict = self.yaml.load(self.as_string)

    @property
    def as_file(self):
        return io.StringIO(self.as_string)

    @property
    def id(self):
        return Node(self.project, self.deployment)

    @property
    def deployment(self):
        return self.as_dict.get(
            'name',
            os.path.basename(self.source).split('.')[0]
        )

    @property
    def project(self):
        """ Sets the project for the config

        This is a bit complicated but allows for quite a bit of
        flexibility. The project can be defined in a few different
        places, and this is the order on precedence:

        1- Command line
        2- Config file
        3- CLOUD_FOUNDATION_PROJECT_ID environment variable
        4- The GCP SDK configuration
        """
        if not hasattr(self, '_project'):
            self._project = self.as_dict.get('project') or \
                os.environ.get('CLOUD_FOUNDATION_PROJECT_ID') or \
                dm_base.GetProject()
        return self._project

    @property
    def dependencies(self):
        """
        """
        if hasattr(self, '_dependencies'):
            return self._dependencies

        self._dependencies = set()
        for line in self.as_file.readlines():
            # Ignore comments
            if re.match(r'^\s*#', line):
                continue

            # Match !DMOutput, $(out.x.y.w.z), etc tokens
            for match in DM_OUTPUT_QUERY_REGEX.finditer(line):
                for k, v in match.groupdict().items():
                    if not v:
                        continue
                    if k == 'url':
                        url = parse_dm_output_url(v, self.project)
                    elif k == 'token':
                        url = parse_dm_output_token(v, self.project)
                    self._dependencies.add(Node(url.project, url.deployment))

        return self._dependencies

    def __repr__(self):
        return '{}({}:{})'.format(self.__class__, self.deployment, self.project)


class ConfigGraph(object):
    """ Class representing the dependency graph between configs

    This is a container class holding the dependencies between configs.
    An instance of this class be be used as an iterator over the
    "levels" of dependencies.

    ```
    graph = ConfigGraph(["config-1.yaml", "config-2.yaml"])
    for level in graph:
        for config in level:
            deployment = Deployment(config)
            ...
    ```

    Attributes:
        graph(networkx.DiGraph): A networkx DiGraph()
        roots(list): List of all root nodes in the graph
        levels(list): List of dependency levels. Each element in the
           list is another list of nodes that can be processed in
           parallel.

    """

    def __init__(self, configs, project=None):
        """ Constructor """

        # Populate the config dict
        self.configs = {
            c.id: c for c in (Config(x,
                                     project=project) for x in configs)
        }

    @property
    def graph(self):
        if hasattr(self, '_graph'):
            return self._graph
        self._graph = nx.DiGraph()
        for _, config in self.configs.items():
            node = Node(config.project, config.deployment)
            self._graph.add_node(node)
            for dependency in config.dependencies:
                self.graph.add_edge(dependency, node)

            if not nx.is_directed_acyclic_graph(self._graph):
                raise SystemExit('Cyclic dependency in the graph')
        return self._graph

    @property
    def roots(self):
        if not hasattr(self, '_roots'):
            self._roots = [
                n for n in self.sort() if not list(self.graph.predecessors(n))
            ]
        return self._roots

    @property
    def levels(self):
        if hasattr(self, '_levels'):
            return self._levels

        graph = self.graph.copy()
        remaining_nodes = list(self.sort())
        self._levels = []

        while remaining_nodes:
            level_nodes, level_configs = [], []

            # Find the nodes in the level
            for node in remaining_nodes:
                if not nx.ancestors(graph, node):
                    level_nodes.append(node)

            # Find and load configs in the level
            # If a node is not in a provided config, it must be a
            # dependency, so we make sure it exists in DM, without
            # attempting to load an unexisting config
            for node in level_nodes:
                remaining_nodes.remove(node)
                graph.remove_node(node)

                if node in self.configs:
                    level_configs.append(self.configs[node])
                else:
                    deployment = get_deployment(node.project, node.deployment)
                    if not deployment:
                        raise SystemExit(
                            'Unresolved dependency. Resource {}, on which'
                            'other resources depended, neither was specified'
                            'in the submitted congigs nor existed in'
                            'Deployment Manager'.format(node)
                        )

            if level_configs:
                self._levels.append(level_configs)

        return self._levels

    def __iter__(self):
        """ Makes this class an iterator.

        Notice the iterator over `self.levels` not `self`
        """
        return iter(self.levels)

    def __reversed__(self):
        """ Class can be iterated in reverse order. """
        return reversed(self.levels)

    def sort(self, reverse=False):
        """ Sorts the graph in topological order.


        Args:
            reverse(boolean): Whether to return the nodes in reverse
                order on not.

        Returns: A generator of nodes sorted by topology, ie the
            elements are returned sequentially in order of dependency (
            independent nodes come first, unless 'reverse' is used.
        """
        generator = nx.topological_sort(self.graph)
        if reverse:
            return reversed(list(generator))
        return generator


class Deployment(DM_API):
    """Class representing a CFT deployment.

    This class makes extensive use of the Google Cloud SDK. Relevant files to
    understand some of this code:
    https://github.com/google-cloud-sdk/google-cloud-sdk/tree/master/lib/surface/deployment_manager/deployments
    https://github.com/google-cloud-sdk/google-cloud-sdk/blob/master/lib/googlecloudsdk/third_party/apis/deploymentmanager/v2/deploymentmanager_v2_messages.py
    https://github.com/google-cloud-sdk/google-cloud-sdk/blob/master/lib/googlecloudsdk/third_party/apis/deploymentmanager/v2/deploymentmanager_v2_client.py

    Attributes:
        config(dict): A dict holding the config for this deployment.
        current(Deployment): A Deployment object from the SDK, or None.
            This attribute is None until self.get() called. If the
            deployment doesn't exist in DM, it remains None.
        dm_config(dict): A dict built from the CFT config holding keys
            that DM can handle.
        target_config(TargetConfiguration): A TargetConfiguration object from
            the SDK.
    """

    # Number of seconds to wait for a create/update/delete operation
    OPERATION_TIMEOUT = 20 * 60  # 20 mins. Same as gcloud

    # The keys required by a DM config (not CFT config)
    DM_CONFIG_KEYS = ['imports', 'resources', 'outputs', 'configVersion']

    def __init__(self, config):
        """ The class constructor

        Args:
            config_item (configItem): A dict representing CFT config.
                Normally provided when creating/updating a deployment.
        """

        # Resolve custom yaml tags only during deployment instantiation
        # because if parsed earlier, the DM queries implemented for the
        # tags would likely fail with 404s
        self.yaml = CFTBaseYAML()
        self.yaml.Constructor.add_constructor(
            '!DMOutput',
            self.yaml_dm_output_constructor
        )
        self._config = config

        # Regex search/replace before loading the yaml
        self.config = self.yaml.load(config.as_string)
        self.yaml_walk(self.config)

        self.config['project'] = self._config.project
        self.config['name'] = self._config.deployment

        self.tmp_file_path = None

        LOG.debug('==> %s', self.config)
        self.current = None

    def yaml_walk(self, yaml_tree):
        """ Custom function for walking through the config and checking every string if its a regexp match

        In place walk over the config yaml. In case of a string togen it replaces the token with the complex
        YAML value of the reference.

        The function is able to walk through lists and dictionarries. It ignores boolm, int and double values.
        """
        if isinstance(yaml_tree, dict):
            for k, v in yaml_tree.items():  ## Walk each element in dictionary
                yaml_tree[k] = self.yaml_replace(v)
        elif isinstance(yaml_tree, list):
            for i, v in enumerate(yaml_tree):  ## Walk each element in list
                yaml_tree[i] = self.yaml_replace(v)

    def yaml_replace(self, v):
        if isinstance(v, str):
            match = DM_OUTPUT_QUERY_REGEX.match(v)
            if match is not None:
                return self.get_dm_output(match)
        else:
            self.yaml_walk(v)  ## Not string, recursive walk
        return v

    def get_dm_output(self, match):
        """ Custom function for the regex.match()

        This function gets executed everytime there's a match on one
        tokens used to represent the cross-deployment references (
        !DMOutput, $(out.x.y.w.z), etc.

        Args:
            match (re.MatchObject): A regex matche object

        Returns: A string with the value of the deployment output
        """

        for k, v in match.groupdict().items():
            if not v:
                continue
            if k == 'url':
                query_attributes = parse_dm_output_url(v, self._config.project)
            elif k == 'token':
                query_attributes = parse_dm_output_token(
                    v,
                    self._config.project
                )
            return get_deployment_output(
                query_attributes.project,
                query_attributes.deployment,
                query_attributes.resource,
                query_attributes.name
            )

    def yaml_dm_output_constructor(self, loader, node):
        """ Implements the !DMOutput yaml tag

        The tag takes string represeting an DM item URL.

        Example:
          network: !DMOutput dm://${project}/${deployment}/${resource}/${name}
        """

        data = loader.construct_scalar(node)
        url = parse_dm_output_url(data, self._config.project)
        return get_deployment_output(
            url.project,
            url.deployment,
            url.resource,
            url.name
        )

    @property
    def dm_config(self):
        """Returns a dict with keys that DM can handle.

        Args:

        Return: A dict representing a valid DM config (not CFT config)

        TODO (gus): Could a dictview be used here?
        """

        return {
            k: v for k,
            v in self.config.items() if k in self.DM_CONFIG_KEYS
        }

    @property
    def target_config(self):
        """Returns the 'target config' for the deployment.

        The 'import code' is very complex and error prone. Instead
        of rewriting it here, the code from the SDK/gcloud is being
        reused.
        The SDK code only works with actual files, not strings, so
        the processed configs to are written to temporary files then
        fed to the SDK code to handle the imports.

        Args:

        Returns: None
        """
        self.write_tmp_file()
        target = BuildTargetConfig(messages, config=self.tmp_file_path)
        self.delete_tmp_file()
        return target

    def write_tmp_file(self):
        """ Writes the yaml dump of the deployment to a temp file.

        This temporary file is always created in the current directory,
        not in the directory where the config file is.

        Args:

        Returns: None
        """

        with tempfile.NamedTemporaryFile(dir=os.getcwd(), delete=False) as tmp:
            self.yaml.dump(self.dm_config, tmp)
            self.tmp_file_path = tmp.name

    def delete_tmp_file(self):
        """ Delete the temporary config file """

        os.remove(self.tmp_file_path)

    def get(self):
        """ Returns a Deployment() message(obj) from the DM API.

        Shortcut to deployments.Get() that doesn't raise an exception
        When deployment doesn't exit.

        This method also updates the 'current' attribute with the latest
        data from the DM API.

        Args:

        Returns: A Deployment object from the SDK or None
        """

        self.current = get_deployment(
            project=self.config['project'],
            deployment=self.config['name']
        )
        return self.current

    def delete(self, delete_policy=None):
        """Deletes this deployment from DM.

        Args:
            delete_policy (str): The strings 'ABANDON' or 'DELETE'.
                The default (None), doesn't include the policy in the
                request obj, which translates 'DELETE' as default.

        Returns: None
        """

        message = self.messages.DeploymentmanagerDeploymentsDeleteRequest
        request = message(
            deployment=self.config['name'],
            project=self.config['project']
        )

        if delete_policy:
            request['deletePolicy'
                   ] = message.DeletePolicyValueValuesEnum(delete_policy)

        LOG.debug('Deleting deployment %', self.config['name'], request)

        # The actual operation.
        # No exception handling is done here to allow higher level
        # functions to do so.
        operation = self.client.deployments.Delete(request)

        # Wait for operation to finish
        self.wait(operation)

    def create(self, preview=False, create_policy=None):
        """Creates this deployment in DM.

        Args:
            preview (boolean): If True, create is done with preview.
            create_policy (str): The strings 'ACQUIRE' or 'CREATE_OR_ACQUIRE'.
                The default (None), doesn't include the policy in the
                request obj, which translates 'CREATE_OR_ACQUIRE' as default.

        Returns: None
        """

        deployment = self.messages.Deployment(
            name=self.config['name'],
            target=self.target_config
        )

        message = self.messages.DeploymentmanagerDeploymentsInsertRequest
        request = message(
            deployment=deployment,
            project=self.config['project'],
            preview=preview
        )
        if create_policy:
            request['createPolicy'
                   ] = message.CreatePolicyValueValuesEnum(create_policy)
        LOG.debug(
            'Creating deployment %s with data %s',
            self.config['name'],
            request
        )

        # The actual operation.
        # No exception handling is done here to allow higher level
        # functions to do so.
        operation = self.client.deployments.Insert(request)

        # Wait for operation to finish
        self.wait(operation)
        self.print_resources_and_outputs()
        return self.current


#
#        if preview:
#            func = self.confirm_preview()
#            func()
#        elif getattr(self.current, 'update', False):
#            self.update_preview()
#

    def update(self, preview=False, create_policy=None, delete_policy=None):
        """Updates this deployment in DM.

        If the deployment is already in preview mode in DM, the existing
        preview operation will be overwritten by this one.

        Args:
            preview (boolean): If True, update is done with preview.
            create_policy (str): The strings 'ACQUIRE' or 'CREATE_OR_ACQUIRE'.
                The default (None), doesn't include the policy in the
                request obj, which translates 'CREATE_OR_ACQUIRE' as default.
            delete_policy (str): The strings 'ABANDON' or 'DELETE'.
                The default (None), doesn't include the policy in the
                request obj, which translates 'DELETE' as default.

        Returns: None
        """

        # Get current deployment to figure out the fingerprint
        self.get()
        if not self.current:
            raise SystemExit(
                'Error updating {}: Deployment does not exist'.format(
                    self.config['name']
                )
            )

        new_deployment = self.messages.Deployment(
            name=self.config['name'],
            target=self.target_config,
            fingerprint=self.current.fingerprint or b''
        )

        message = self.messages.DeploymentmanagerDeploymentsUpdateRequest

        # getattr() below overwrites existing preview mode as targets
        # cannot be sent when deployment is already in preview mode
        request = message(
            deployment=self.config['name'],
            deploymentResource=new_deployment,
            project=self.config['project'],
            preview=preview or bool(getattr(self.current,
                                            'update',
                                            False))
        )
        if delete_policy:
            request['deletePolicy'
                   ] = message.DeletePolicyValueValuesEnum(delete_policy)
        if create_policy:
            request['createPolicy'
                   ] = message.CreatePolicyValueValuesEnum(create_policy)

        LOG.debug(
            'Updating deployment %s with data %s',
            self.config['name'],
            request
        )

        # The actual operation.
        # No exception handling is done here to allow higher level
        # functions to do so.
        operation = self.client.deployments.Update(request)

        # Wait for operation to finish
        self.wait(operation)

        self.print_resources_and_outputs()

        if preview:
            func = self.confirm_preview()
            func()
        elif getattr(self.current, 'update', False):
            self.update_preview()

    def confirm_preview(self):
        answer = ask()

        if answer == 'u':
            return self.update_preview
        elif answer == 's':
            return self.cancel_preview
        elif answer == 'a':
            raise SystemExit('Aborting deployment run!')
        else:
            raise SystemExit('Not a valid answer: {}'.format(answer))

    def update_preview(self):
        """Confirms an update preview.

        The request to the API doesn't include the target

        Args:

        Returns:
        """
        deployment = self.messages.Deployment(
            name=self.config['name'],
            fingerprint=self.current.fingerprint or b''
        )
        request = self.messages.DeploymentmanagerDeploymentsUpdateRequest(
            deployment=self.config['name'],
            deploymentResource=deployment,
            project=self.config['project'],
            preview=False
        )
        operation = self.client.deployments.Update(request)
        self.wait(operation, 'update preview')
        self.print_resources_and_outputs()

    def wait(self, operation, action=None, get=True):
        """Waits for a DM operation to be completed.

        Args:
            operation (Operation): An Operation object from the SDK.
            action (string): Any operation name to be used in the
                ticker. If not specified, the operation type is used.
            get (boolean): wether to retrieve the latest deployment
                info from the API to obtain the current fingerprint.
        """
        # This saves an API call if the self.get() was called just
        # before calling this method
        if get:
            self.get()

        action = action or operation.operationType

        dm_write.WaitForOperation(
            self.client,
            self.messages,
            operation.name,
            project=self.config['project'],
            timeout=self.OPERATION_TIMEOUT,
            operation_description='{} {} (fingerprint {})'.format(
                action,
                self.config['name'],
                base64.urlsafe_b64encode(self.current.fingerprint)
            )
        )
        return self.get()

    def cancel_preview(self):
        """Cancels a deployment preview.

        If a deployment is in preview mode, the update is cancelled and
        no resourced are changed

        Args:

        Returns:
        """
        cancel_msg = self.messages.DeploymentsCancelPreviewRequest(
            fingerprint=self.current.fingerprint or b''
        )
        req = self.messages.DeploymentmanagerDeploymentsCancelPreviewRequest(
            deployment=self.config['name'],
            deploymentsCancelPreviewRequest=cancel_msg,
            project=self.config['project']
        )
        operation = self.client.deployments.CancelPreview(req)
        self.wait(operation)

    def apply(self, preview=False, create_policy=None, delete_policy=None):
        """Creates or updates this deployment in DM.

        Args:
            preview (boolean): If True, update is done with preview.
            create_policy (str): The strings 'ACQUIRE' or 'CREATE_OR_ACQUIRE'.
                The default (None), doesn't include the policy in the
                request obj, which translates 'CREATE_OR_ACQUIRE' as default.
            delete_policy (str): The strings 'ABANDON' or 'DELETE'.
                The default (None), doesn't include the policy in the
                request obj, which translates 'DELETE' as default.

        Returns: None
        """
        try:
            self.create()
        except apitools_exceptions.HttpConflictError as err:
            self.update(preview=preview)

    def print_resources_and_outputs(self):
        """Prints the Resources and Outputs of this deployment."""

        rsp = dm_api_util.FetchResourcesAndOutputs(
            self.client,
            self.messages,
            self.config['project'],
            self.config['name'],
            #           self.ReleaseTrack() is base.ReleaseTrack.ALPHA
        )

        printer = resource_printer.Printer(flags.RESOURCES_AND_OUTPUTS_FORMAT)
        printer.AddRecord(rsp)
        printer.Finish()
        return rsp
