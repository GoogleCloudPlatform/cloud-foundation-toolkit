from six import PY2

from apitools.base.py.exceptions import HttpNotFoundError
import jinja2
import pytest
from ruamel.yaml import YAML

from cloud_foundation_toolkit.deployment import Config
from cloud_foundation_toolkit.deployment import ConfigGraph
from cloud_foundation_toolkit.deployment import Deployment

if PY2:
    import mock
else:
    import unittest.mock as mock

class Message():
    def __init__(self, **kwargs):
        [setattr(self, k, v) for k, v in kwargs.items()]


@pytest.fixture
def args():
    return Args()

def test_config(configs):
    c = Config(configs.files['my-networks.yaml'].path)
    assert c.as_string == configs.files['my-networks.yaml'].jinja


def test_config_list(configs):
    config_paths = [v.path for k, v in configs.files.items()]
    config_list = ConfigGraph(config_paths)
    for level in config_list:
        assert isinstance(level, list)
        for c in level:
            assert isinstance(c, Config)


def test_deployment_object(configs):
    config = Config(configs.files['my-networks.yaml'].path)
    deployment = Deployment(config)
    assert deployment.config['name'] == 'my-networks'


def test_deployment_get(configs):
    config = Config(configs.files['my-networks.yaml'].path)
    deployment = Deployment(config)
    with mock.patch.object(deployment.client.deployments, 'Get') as m:
        m.return_value = Message(
            name='my-networks',
            fingerprint='abcdefgh'
        )
        d = deployment.get()
        assert d is not None
        assert deployment.current == d


def test_deployment_get_doesnt_exist(configs):
    config = Config(configs.files['my-networks.yaml'].path)
    deployment = Deployment(config)
    with mock.patch('cloud_foundation_toolkit.deployment.get_deployment') as m:
        m.return_value = None
        d = deployment.get()
        assert d is None
        assert deployment.current == d


def test_deployment_create(configs):
    config = Config(configs.files['my-networks.yaml'].path)
    patches = {
        'client': mock.DEFAULT,
        'wait': mock.DEFAULT,
        'get': mock.DEFAULT,
        'print_resources_and_outputs': mock.DEFAULT
    }

    with mock.patch.multiple(Deployment, **patches) as mocks:
        deployment = Deployment(config)
        mocks['client'].deployments.Insert.return_value = Message(
            name='my-network-prod',
            fingerprint='abcdefgh'
        )
        mocks['client'].deployments.Get.return_value = Message(
            name='my-network-prod',
            fingerprint='abcdefgh'
        )

        d = deployment.create()
        assert deployment.current == d
