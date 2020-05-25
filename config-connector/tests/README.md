# Config Connector Solutions Testing CLI

## Introduction

This folder contains the Go CLI and testcases for testing the Config Connector
Solutions defined in [../solutions](../solutions) folder.

*   **[ccs-test/](./ccs-test/)** - Go code for the solutions test CLI
*   **[testcases/](./testcases/)** - Testcases for each solution. If has
    the same folder structure as the solutions, i.e. if the solution is under
    <code>../solutions/<b>iam/kpt/member-iam/</b></code>, then the corresponding
    testcases should be under <code>./testcases/<b>iam/kpt/member-iam/</b>
    </code>

## Requirements

*   [gsutil](https://cloud.google.com/storage/docs/gsutil_install)
*   [kpt](../solutions/README.md#kpt)
*   [helm](../solutions/README.md#helm)
*   a working Kubernetes cluster with Config Connector [installed and
    configured](
    https://cloud.google.com/config-connector/docs/how-to/install-upgrade-uninstall)
    *   [Default namespace](
        https://cloud.google.com/config-connector/docs/how-to/install-upgrade-uninstall#setting_your_default_namespace)
        should be [configured to the **project** where you want to manage the GCP
        resources](
        https://cloud.google.com/config-connector/docs/how-to/install-upgrade-uninstall#specify_where_to_create_your_resources).

## Consumption

1.  Clone GoogleCloudPlatform/cloud-foundation-toolkit repository:
  
    ```
    git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
    ```

1.  Go to the tests folder:

    ```
    cd cloud-foundation-toolkit/config-connector/tests
    ```

1.  Download the `test-cli` executable file:

    ```
    gsutil cp gs://kcc-solutions-test/test-cli test-cli
    ```
1.  Change the file ACL to make it executable:

    ```
    chmod +x test-cli
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
./test-cli run --path [RELATIVE_PATH]  # E.g. "iam/kpt/member-iam"
```

Most test should take a few minutes to finish. But you'll need to specify the
timeout using the optional `--timeout` or `-t` flag for special test cases:

*   [projects/kpt/shared-vpc](../solutions/projects/kpt/shared-vpc): 10m
    ```
    ./test-cli run --path projects/kpt/shared-vpc --timeout 10m
    ```

After you run the command, detailed output will be printed out. If you find the
last line of the output is `======Successfully finished the test for solution
RELATIVE_PATH]======`, it means the test run is successful. Otherwise, you'll
find the detailed error message for the failure.

### Exceptions

Solutions that require manual steps can't be tested using our `test-cli`. Here
is the list of exceptions:

*   [projects/kpt/project-hierarchy](
    ../solutions/projects/kpt/project-hierarchy) - need to manually figure out
    the folder ID before creating projects ([GitHub issue](
    https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/104))

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

1.  Create the testcase YAML file `required_fields_only.yaml` from
    `test_values.template` file:

    ```
    cp test_values.template [SOLUTION_AREA]/kpt/[SOLUTION_NAME]/required_fields_only.yaml
    ```

    **Note:** Currently we only support one testcase, which only set required
    kpt setters (setters set by PLACEHOLDER). Setting optional kpt setters in
    tests is not necessary except for the SQLInstance name. This issue will be
    addressed in the upcoming changes.

1.  Check if the solution requires any PLACEHOLDERs to be set:

    ```
    kpt cfg list-setters ../../solutions/[SOLUTION_AREA]/kpt/[SOLUTION_NAME]
    ```

1.  For each setter that is a placeholder, decide if the value should be a **new
    globally unique** value. E.g., the value of a new project ID.

    1.  If the value **MUST** be globally unique, append the following key-value
        pair in the testcase YAML file:

        ```
        # \$ENV_VAR is the placeholder to reference to ENV_VAR you've set in the
        # environments file (./environments.yaml).
        # In order to create globally unique resource names, you need to append
        # `-\$RANDOM_ID` after the `\$ENV_VAR`. E.g. `\$PROJECT_ID-\$RANDOM_ID`.
        echo "[SETTER_NAME]: \$ENV_VAR-\$RANDOM_ID" >> \
        [SOLUTION_AREA]/kpt/[SOLUTION_NAME]/required_fields_only.yaml
        ```

    1.  If the value doesn't need to be globally unique, append the following
        key-value pair in the testcase YAML file:
        ```
        # \$ENV_VAR is the placeholder to reference to ENV_VAR you've set in the
        # environments file (./environments.yaml). E.g. `\$PROJECT_ID`.
        echo "[SETTER_NAME]: \$ENV_VAR" >> \
        [SOLUTION_AREA]/kpt/[SOLUTION_NAME]/required_fields_only.yaml
        ```

    **Note:** Please don't use $ENV_VAR directly in the command. The back slash
    ("\\") is necessary because here, it is a string, but not a variable. We
    don't want to set the value of ENV_VAR in the testcase YAML file.

    **Note:** `$RANDOM_ID` is a placeholder for the autogen randomized suffix,
    and `RANDOM_ID` shouldn't be the name of the env var.

1.  Check the environments template file
    ([./environments.template](./environments.template)). For each environment
    variable you need but does **NOT** exist in the environments template file,
    add it:

    ```
    echo "ENV_VAR: \${ENV_VAR?}" >> ./environments.template
    ```

1.  Create the YAML file for the original values of the setters from
    `test_values.template` file:

    ```
    cp test_values.template [SOLUTION_AREA]/kpt/[SOLUTION_NAME]/original_values.yaml
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

# License

  Apache 2.0 - See [LICENSE](/LICENSE) for more information.
