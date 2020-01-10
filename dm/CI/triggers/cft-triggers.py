# This is a wrapper tamplete to crawl through the dm/templates folder and
# create a trigger for each.
# This is not a generic template, used for CFT GitHub PR testing

import copy


def generate_config(context):

    tests = []
    resources = []
    for test in context.imports:
        if '/tests/integration/' in test:
            testData = test.split('/')
            testFolder = testData[3]
            batsFile = testData[6]

            props = copy.deepcopy(context.properties)
            props['description'] = props['description'].replace('#template#', batsFile[:-5])
            props['substitutions']['_BATS_TEST_FILE'] = \
                props['substitutions']['_BATS_TEST_FILE'].replace(
                '#template#', testFolder).replace(
                '#templatetest#', batsFile)
            for i in range(len(props['includedFiles'])):
                props['includedFiles'][i] = props['includedFiles'][i].replace(
                    '#template#', testFolder)
            resources.append({
                'type': "cft-trigger.py",
                'name': context.env['name'] + "-" + batsFile[:-5],
                'properties': props})

    return {'resources': resources}
