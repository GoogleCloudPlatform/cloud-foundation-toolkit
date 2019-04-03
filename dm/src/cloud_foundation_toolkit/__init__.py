import logging
import pkg_resources

from googlecloudsdk.core.credentials import store as creds_store

# Setup logging and expose Logger object to the rest of the project
LOG = logging.getLogger("cft")
LOG.addHandler(logging.StreamHandler())
LOG.propagate = False

__VERSION__ = pkg_resources.get_distribution(__name__).version

# Register credentials providers - for instance SA, etc
credential_providers = [
    creds_store.DevShellCredentialProvider(),
    creds_store.GceCredentialProvider(),
]
for provider in credential_providers:
    provider.Register()
