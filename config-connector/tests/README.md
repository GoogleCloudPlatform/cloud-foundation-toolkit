# Config Connector Solutions Testing CLI

## Introduction

This folder contains the Go CLI and testcases for testing the Config Connector
Solutions defined in [../solutions](../solutions) folder.

*  **[ccs-test/](./ccs-test/)** - Go CLI
*  **[testcases/](./testcases/)** - Testcases for each solution. If has
   the same folder structure as the solutions, i.e. if the solution is under 
   <code>../solutions/<b>iam/kpt/member-iam/</b></code>, then the corresponding
   testcases should be under <code>./testcases/<b>iam/kpt/member-iam/</b></code>

## Requirements

*  [Go](https://golang.org/doc/install)
*  [kpt](../solutions/README.md#kpt)
*  [helm](../solutions/README.md#helm)
*  a working Kubernetes cluster with Config Connector [installed and 
   configured](https://cloud.google.com/config-connector/docs/how-to/install-upgrade-uninstall)

## Consumption

1.  Clone GoogleCloudPlatform/cloud-foundation-toolkit repository under your `$GOPATH`:
  
    ```
    git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
    ```

1.  Go to the tests folder:

    ```
    cd cloud-foundation-toolkit/config-connector/tests
    ```

1.  Install ccs-test CLI:

    ```
    go install ./ccs-test
    ```

1.  Set the environment variables for the tests.

    1.  Make a copy of the environments template file
        ([./testcases/environments.template](./testcases/environments.template)) :

        ```
        cp ./testcases/environments.template ./testcases/environments.yaml
        ```

    1.  Edit the environments file
        ([./testcases/environments.yaml](./testcases/environments.yaml)) to
        update the environment variables. Use **any command or editing tool**
        you prefer. E.g. you can use `sed` command:

        ```
        # YOUR_PROJECT_ID should be the project ID that your default namespace
        # is annotated with.
        sed -i 's/${DEFAULT_PROJECT_ID?}/[YOUR_PROJECT_ID]/g' ./testcases/environments.yaml
        ```

        **Note:** Please remember to set **ALL** the environment variables.

1.  Follow the README of each solution to configure permissions for the
    cnrm-system service account, and enable necessary APIs.

    **Note:** This step will be automated in the upcoming changes.

## How to run the tests?

**Note:** Currently we only support testing one solution each time by setting
the relative path of the solution using `--path` or `-p` flag.

**Note:** Currently we only support testing kpt solutions specified under
[testcases folder](./testcases).

Under the [tests](.) folder, run a test by providing the relative path:
```
ccs-test run --path [RELATIVE_PATH]  # E.g. "iam/kpt/member-iam"
```

Each test should take less than 2 mins to finish. You'll find the detailed
output of the test after you run the command.

If you find the last line of the output is `======Successfully finished the test
for solution [RELATIVE_PATH]======`, it means the test run is successful.
Otherwise, you'll find the detailed error message of the failure.

## How to add new tests?

**Note:** Currently we only support adding tests for kpt solutions.

If you want to create tests for solution 
`[SOLUTION_AREA]/kpt/[SOLUTION_NAME]` (e.g. `iam/kpt/member-iam`):

1.  Under your local copy of your
    [forked](https://help.github.com/en/github/getting-started-with-github/fork-a-repo)
    cloud-foundation-toolkit repository, go to the testcases folder:

    ```
    cd cloud-foundation-toolkit/config-connector/tests/testcases
    ```

1.  Create the folder if it doesn't exist:

    ```
    mkdir -p [SOLUTION_AREA]/kpt/[SOLUTION_NAME]
    ```

1.  Create the testcase YAML file `required_fields_only.yaml`:

    ```
    touch [SOLUTION_AREA]/kpt/[SOLUTION_NAME]/required_fields_only.yaml
    ```

    **Note:** Currently we only support one testcase, which only set required
    kpt setters (setters set by PLACEHOLDER). Setting optional kpt setters in
    tests is not necessary except for the SQLInstance name. This issue will be
    addressed in the upcoming changes.

1.  Check if the solution requires any PLACEHOLDERs to be set:

    ```
    kpt cfg list-setters ../../solutions/[SOLUTION_AREA]/kpt/[SOLUTION_NAME]
    ```

1.  For each setter that is a placeholder, append the follow key-value pair in
    the testcase YAML file:

    ```
    # \$ENV_VAR is the placeholder to reference to ENV_VAR you've set in the
    # environments file (./environments.yaml). E.g. `\$PROJECT_ID`.
    echo "[SETTER_NAME]: \$ENV_VAR" >> \
    [SOLUTION_AREA]/kpt/[SOLUTION_NAME]/required_fields_only.yaml
    ```

    **Note:** Please don't use $ENV_VAR directly in the command. The back slash
    ("\\") is necessary because here, it is a string, but not a variable. We
    don't want to set the value of ENV_VAR in the testcase YAML file.

1.  Check the environments template file
    ([./environments.template](./environments.template)). For each environment
    variable you need but doesn't exist in the environments template file, add
    it:

    ```
    echo "ENV_VAR: \${ENV_VAR?}" >> ./environments.template
    ```

1.  Create the YAML file for the original values of the setters:

    ```
    touch [SOLUTION_AREA]/kpt/[SOLUTION_NAME]/original_values.yaml
    ```

1.  For each placeholder setter and its original value (you can find them by
    running
    `kpt cfg list-setters ../../solutions/[SOLUTION_AREA]/kpt/[SOLUTION_NAME]`),
    append the key-pairs to the original values YAML file:

    ```
    # You need to add the back slash ("\") in front of the original value
    # because it is placeholder starting with "$".
    echo "[SETTER_NAME]: \[ORIGINAL_VALUE]" >> \
    [SOLUTION_AREA]/kpt/[SOLUTION_NAME]/original_values.yaml
    ```
