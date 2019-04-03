from six import PY2

from apitools.base.py.exceptions import HttpNotFoundError
import pytest
from ruamel.yaml import YAML

from cloud_foundation_toolkit.dm_utils import API
from cloud_foundation_toolkit.dm_utils import get_deployment


if PY2:
    import mock
else:
    import unittest.mock as mock

class Message():
    def __init__(self, **kwargs):
        [setattr(self, k, v) for k, v in kwargs.items()]


def test_get_deployment():
    with mock.patch.object(API.client.deployments, 'Get') as m:
        m.side_effect = HttpNotFoundError('a', 'b', 'c')
        d = get_deployment('some-deployment', 'some-project')
        assert d is None

