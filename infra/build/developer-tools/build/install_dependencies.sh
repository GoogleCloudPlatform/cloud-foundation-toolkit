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

set -e
set -u

# coreutils provides xargs and other utilities necessary for lint checks
apk add --no-cache coreutils

# curl is used by unit tests and is nice to have
apk add --no-cache curl

# findutils provides find which is used by lint checks
apk add --no-cache findutils

# git is used to clone repositories
apk add --no-cache git

# grep is used by lint checks
apk add --no-cache grep

# g++ is probably used to install dependencies like psych, but unsure
apk add --no-cache g++

# jq is useful for parsing JSON data
apk add --no-cache jq

# make is used for executing make tasks
apk add --no-cache make

# musl-dev provides the standard C headers
apk add --no-cache python3-dev musl-dev

# openssh is used for ssh-ing into bastion hosts
apk add --no-cache openssh

# unclear why perl is needed, but is good to have
apk add --no-cache perl

# python 3 is needed for python linting
apk add --no-cache python3

# py-pip is needed for installing pip packages
apk add --no-cache py-pip

# ca-certificates is needed to verify the authenticity of artifacts downloaded
# from the internet
apk add --no-cache ca-certificates

# diffutils contains the full version of diff needed for the -exclude argument.
# That argument is needed for check_documentation in task_helper_functions.sh
apk add --no-cache diffutils

# rsync is needed for check_documentation in task_helper_functions.sh
apk add --no-cache rsync

# install gcompat, because protoc needs a real glibc or compatible layer
apk add --no-cache gcompat

# flake8 and jinja2 are used for lint checks, cookiecutter is used for terraform-google-module-template tests
# requests~=2.28 for https://github.com/psf/requests/pull/6179
# Remaining items moved here from requirements.txt
pip install flake8 "Jinja2~=3.1" cookiecutter "requests~=2.28" PyGithub==1.51 google-auth~=1.6 google-api-python-client~=1.7 crcmod
