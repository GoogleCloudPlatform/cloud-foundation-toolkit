#! /bin/bash
# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -u

# Import helper functions from Docker container
# shellcheck disable=SC1091
source /usr/local/bin/task_helper_functions.sh

# All tests should execute regardless of failures so we can report at the end.
# Perform 'set +e' after sourcing task_helper_functions.sh in case someone has
# 'set -e' inside a module-specific task_helper_functions.sh file.
set +e

# constants
MARKDOWN=0
MARKDOWN_STR=""
CONTRIBUTING_GUIDE=""
# shellcheck disable=SC2089,SC2016  # Quotes/backslashes will be treated literally, expressions don't expand
messages='{
  "check_generate_modules": "The modules need to be regenerated. Please run `make_build`.",
  "check_documentation": "The documentation needs to be regenerated. Please run `make generate_docs`.",
  "check_whitespace": "Failed whitespace check. More details below.",
  "check_shell": "Failed shell check. More info on running shellcheck locally [here](https://www.shellcheck.net).",
  "check_headers": "All files need a license header. Please make sure all your files include the appropriate header. A helper tool available [here](https://github.com/google/addlicense).",
  "check_python": "Failed flake8 Python lint check.",
  "check_terraform": "Failed Terraform check. More details below."
}'
rval=0
failed_tests=()
tests=(
  check_generate_modules
  check_documentation
  check_whitespace
  check_shell
  check_headers
  check_python
  check_terraform
)

# parse args
for arg in "$@"
do
  case $arg in
    -m|--markdown)
      MARKDOWN=1
      shift
      ;;
    -c=*|--contrib-guide=*)
      CONTRIBUTING_GUIDE="${arg#*=}"
      shift
      ;;
      *) # end argument parsing
      shift
      ;;
  esac
done

for test in "${tests[@]}"; do
  # if not in markdown mode, pipe test output to stdout tty
  # nested if condition is a workaround for test[[]] not echoing some outputs from check_* tests even with subshell
  if [[ $MARKDOWN -eq 0 ]]; then
    if ! "${test}"; then
      failed_tests+=("${test}")
      ((rval++))
    fi
  # if control reaches here - in markdown mode, pipe test stderr to stdout for capture
  elif ! output=$(${test} 2>&1); then
    # add test name to list of failed_tests
    failed_tests+=("${test}")
    ((rval++))
    # clean output color, sqash multiple empty blank lines
    output=$(echo "$output" | sed -r "s/\x1b\[[0-9;]*m/\n/g" | tr -s '\n')
    # try to get a helpful error message, otherwise unknown
    error_help_message=$(echo "$messages" | jq  --arg check_name "$test" -r '.[$check_name] // "🦖 An unknown error has occurred" ')
    #construct markdown body
    MARKDOWN_STR+="- ⚠️${test}\n ${error_help_message} \n \`\`\`bash \n${output}\n \`\`\` \n"
  fi
done

# if any tests have failed
if [[ "${#failed_tests[@]}" -ne 0 ]]; then
  # echo output in markdown
  if [[ $MARKDOWN -eq 1 ]]; then
  header="Thanks for the PR! 🚀\nUnfortunately it looks like some of our CI checks failed. See the [Contributing Guide](${CONTRIBUTING_GUIDE}) for details.\n"
  echo -e "${header}${MARKDOWN_STR}"
  else
    # shellcheck disable=SC2145  # Output all elements of the array
    echo "Error: The following tests have failed: ${failed_tests[@]}"
    exit "${rval}"
  fi
fi
