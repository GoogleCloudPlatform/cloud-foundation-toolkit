from helper import config_merger


def generate_config(context):

    # Using helper functions to load external configurations.
    # The deployment YAML only contains the minimal context of the deployment.
    # (Module name, environment)
    # This way the wrapper template injects information to the target template
    # without overloading the starting YAML.

    local_properties = config_merger.ConfigContext(
        context.properties['environment'],
        context.properties['module'])

    # Passing values forward to template

    return {
        'resources': [{
            'type': "target-template.py",
            'name': context.env['name'],
            'properties': local_properties}]
    }
