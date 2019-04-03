from six import PY2

if PY2:
    import mock
else:
    import unittest.mock as mock


from apitools.base.py.exceptions import HttpNotFoundError
import pytest

from cloud_foundation_toolkit import actions
from cloud_foundation_toolkit.deployment import Config, ConfigGraph


ACTIONS = ['apply', 'create', 'delete', 'update']


class Args(object):

    def __init__(self, **kwargs):
        self.preview = False
        self.project = False
        self.show_stages = False
        self.format = 'human'
        [setattr(self, k, v) for k, v in kwargs.items()]


def get_number_of_elements(items):
    if isinstance(item, list):
        return sum(get_number_of_elements(subitem) for subitem in item)
    else:
        return 1


def test_execute(configs):
    args = Args(action='apply', config=[configs.directory])
    with mock.patch('cloud_foundation_toolkit.actions.Deployment') as m1:
        graph = ConfigGraph([v.path for k, v in configs.files.items()])
        n_configs = len(configs.files)

        r = actions.execute(args)
        assert r == None
        assert m1.call_count == n_configs

        args.show_stages = True
        r = actions.execute(args)
        assert r == None
        assert m1.call_count == n_configs

        with mock.patch('cloud_foundation_toolkit.actions.json.dumps') as m2:
            args.format = 'json'
            r = actions.execute(args)
            assert m1.call_count == n_configs
            assert m2.call_count == 1

        with mock.patch('cloud_foundation_toolkit.actions.YAML.dump') as m2:
            args.format = 'yaml'
            r = actions.execute(args)
            assert m1.call_count == n_configs
            assert m2.call_count == 1





def test_valid_actions():
    ACTUAL_ACTIONS = actions.ACTION_MAP.keys()
    ACTUAL_ACTIONS.sort()
    assert ACTUAL_ACTIONS == ACTIONS


def test_action(configs):
    args = Args(config=[configs.directory])
    for action in ACTIONS:
        args.action = action
        args.show_stages = False
        n_configs = len(configs.files)
        with mock.patch('cloud_foundation_toolkit.actions.Deployment') as m1:
            # Test the normal/expected flow of the function
            r = actions.execute(args)
            method = getattr(mock.call(), action)
            assert m1.call_count == n_configs
            if action == 'delete':
                assert m1.mock_calls.count(method()) == n_configs
            else:
                assert m1.mock_calls.count(method(preview=args.preview)) == n_configs

            # Test exception handling in the function
            m1.reset_mock()
            getattr(m1.return_value, action).side_effect = HttpNotFoundError('a', 'b', 'c')
            if action == 'delete':
                # if delete is called, execute() should catch the exception
                # and keep going as if nothing happens
                r = actions.execute(args)
                assert m1.mock_calls.count(method()) == n_configs
            else:
                # If exception is raised in any method other than delete,
                # something is really wrong, so exception in re-raised
                # by `execute()`, making the script exit
                # called onde
                with pytest.raises(HttpNotFoundError):
                    r = actions.execute(args)
                assert m1.mock_calls.count(method(preview=args.preview)) == 1

            # Test dry-run
            m1.reset_mock()
            args.show_stages = True
            r = actions.execute(args)
            method = getattr(mock.call(), action)
            m1.assert_not_called()


def test_get_config_files(configs):
    # Test only single directory
    r = actions.get_config_files([configs.directory])
    files = [v.path for k, v in configs.files.items()]
    files.sort()
    r.sort()
    assert files == r

    # Test only files
    files = [v.path for k, v in configs.files.items()]
    r = actions.get_config_files(files)
    files.sort()
    r.sort()
    assert files == r

    # Test files and directories
    confs = [configs.directory] + ['some_file.yaml']
    r = actions.get_config_files(confs)
    files = [v.path for k, v in configs.files.items()] + ['some_file.yaml']
    files.sort()
    r.sort()
    assert files == r


