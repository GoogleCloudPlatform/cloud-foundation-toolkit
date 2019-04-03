from collections import namedtuple
import io
import jinja2
import os
import os.path
import pytest

FIXTURES_DIR = '../fixtures'

ConfigType = namedtuple('ConfigType', ['path', 'content', 'jinja'])

def get_fixtures_fullpath():
    """ Returns the full path for the fixture directory

    Args:

    Returns: The full path to the fixture diretory, Eg:
        /home/vagrant/git/deploymentmanager-samples/community/cloud-foundation/tests/unit/../fixtures/configs/

    """
    return '{}/{}'.format(
        os.path.dirname(os.path.realpath(__file__)),
        FIXTURES_DIR,
    )

def get_configsdir_fullpath():
    """ Returns the full path for a config file fixture

    Args:

    Returns: The full path to the config directory, Eg:
        /home/vagrant/git/deploymentmanager-samples/community/cloud-foundation/tests/unit/../fixtures/configs/config-test-1.yaml

    """
    return '{}/{}'.format(
        get_fixtures_fullpath(),
        'configs'
    )

def get_config_fullpath(config):
    """ Returns the full path for a config file fixture

    Args:
        config (string): The config file name. Eg, config-test-1.yaml

    Returns: The full path to the config file, Eg:
        /home/vagrant/git/deploymentmanager-samples/community/cloud-foundation/tests/unit/../fixtures/configs/config-test-1.yaml

    """
    return '{}/{}'.format(
        get_configsdir_fullpath(),
        config
    )


class Configs():
    directory = get_configsdir_fullpath()

    @property
    def files(self):
        if not hasattr(self, '_files'):
            self._files = {}
            files = [f for f in os.listdir(self.directory) if '.yaml' == f[-5:]]
            for f in files:
                fullpath = get_config_fullpath(f)
                content = io.open(fullpath).read()
                self._files[f] = ConfigType(
                    path=fullpath,
                    content=content,
                    jinja=jinja2.Template(content).render()
                )
        return self._files


@pytest.fixture
def configs():
    return Configs()


if __name__ == '__main__':
    c = Configs()
    print(c.directory)
    for f in c.files:
        print(f.path, f.content)
