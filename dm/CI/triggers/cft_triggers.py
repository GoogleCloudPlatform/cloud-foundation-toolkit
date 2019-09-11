# This is a wrapper tamplete to crawl through the dm/templates folder and
# create a trigger for each.
# This is not a generic template, used for CFT GitHub PR testing
""" This template creates a Cloud Build trigger for each folder under the CFT templates folder. """

import copy


def generate_config(context):
    """ Entry point for the deployment resources. """

    tests = []
    for test in context.imports:
        if '../../templates/' in test:
            tests.append(test[16:-10])

    resources = []
    for test in tests:
        props = copy.deepcopy(context.properties)
        props['description'] = props['description'].replace('#template#', test)
        props['substitutions']['_BATS_TEST_FILE'] = \
            props['substitutions']['_BATS_TEST_FILE'].replace('#template#', test)
        for i in range(len(props['includedFiles'])):
            props['includedFiles'][i] = props['includedFiles'][i].replace(
                '#template#', test)
        resources.append({
            'type': "cft-trigger.py",
            'name': "trigger-" + test,
            'properties': props})

    return {'resources': resources}
