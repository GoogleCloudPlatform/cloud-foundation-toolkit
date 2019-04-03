from collections import namedtuple
import io
import re
from six.moves.urllib.parse import urlparse

from apitools.base.py import exceptions as apitools_exceptions
from googlecloudsdk.api_lib.deployment_manager import dm_base
from ruamel.yaml import YAML

DM_OUTPUT_QUERY_REGEX = re.compile(
    r'!DMOutput\s+(?P<url>\bdm://[-/a-zA-Z0-9]+\b)|'
    r'\$\(out\.(?P<token>[-.a-zA-Z0-9]+)\)'
)

DMOutputQueryAttributes = namedtuple(
    'DMOutputQueryAttributes',
    ['project',
     'deployment',
     'resource',
     'name']
)


@dm_base.UseDmApi(dm_base.DmApiVersion.V2)
class DM_API(dm_base.DmCommand):
    """ Class representing the DM API

    This a proxy class only, so other modules in this project
    only import this local class instead of gcloud's. Here's the source:

    https://github.com/google-cloud-sdk/google-cloud-sdk/blob/master/lib/googlecloudsdk/api_lib/deployment_manager/dm_base.py
    """


API = DM_API()


def get_deployment(project, deployment):
    try:
        return API.client.deployments.Get(
            API.messages.DeploymentmanagerDeploymentsGetRequest(
                project=project,
                deployment=deployment
            )
        )
    except apitools_exceptions.HttpNotFoundError as _:
        return None


def get_manifest(project, deployment):
    deployment_rsp = get_deployment(project, deployment)

    return API.client.manifests.Get(
        API.messages.DeploymentmanagerManifestsGetRequest(
            project=project,
            deployment=deployment,
            manifest=deployment_rsp.manifest.split('/')[-1]
        )
    )


def parse_dm_output_url(url, project=''):
    error_msg = (
        'The url must look like '
        '"dm://${project}/${deployment}/${resource}/${name}" or'
        '"dm://${deployment}/${resource}/${name}"'
    )
    parsed_url = urlparse(url)
    if parsed_url.scheme != 'dm':
        raise ValueError(error_msg)
    path = parsed_url.path.split('/')[1:]

    # path == 2 if project isn't specified in the URL
    # path == 3 if project is specified in the URL
    if len(path) == 2:
        args = [project] + [parsed_url.netloc] + path
    elif len(path) == 3:
        args = [parsed_url.netloc] + path
    else:
        raise ValueError(error_msg)

    return DMOutputQueryAttributes(*args)


def parse_dm_output_token(token, project=''):
    error_msg = (
        'The url must look like '
        '$(out.${project}.${deployment}.${resource}.${name}" or '
        '$(out.${deployment}.${resource}.${name}"'
    )
    parts = token.split('.')

    # parts == 3 if project isn't specified in the token
    # parts == 4 if project is specified in the token
    if len(parts) == 3:
        return DMOutputQueryAttributes(project, *parts)
    elif len(parts) == 4:
        return DMOutputQueryAttributes(*parts)
    else:
        raise ValueError(error_msg)

def get_deployment_output(project, deployment, resource, name):
    manifest = get_manifest(project, deployment)
    layout = YAML().load(manifest.layout)
    return traverse_resource_output(layout, resource, name)

def traverse_resource_output(layout, resource, name):
    for _resource in layout.get('resources', []):
        if _resource['name'] == resource:
            for output in _resource.get('outputs', []):
                if output['name'] == name:
                    return output['finalValue']

        #recursive traversal of complex resources to search for outputs.
        output = traverse_resource_output(_resource, resource, name)
        if output != []:
            return output
    return []
