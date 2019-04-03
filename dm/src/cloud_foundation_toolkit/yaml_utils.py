from ruamel.yaml import YAML
from ruamel.yaml.compat import StringIO


class CFTBaseYAML(YAML):

    def dump(self, data, stream=None, **kwargs):
        inefficient = False
        if stream is None:
            inefficient = True
            stream = StringIO()
        YAML.dump(self, data, stream, **kwargs)
        if inefficient:
            return stream.getvalue()
