#! /bin/bash
# Copyright 2023 Google LLC
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

CURRENT_DIR=$1
WORKING_FOLDER=".working"
BLUPRINT_FOLDER=".blueprint"
GIT_FOLDER="git"
GOLDENS_FOLDER="goldens"
GOLDEN_METADATA="golden-metadata.yaml"
GOLDEN_DISPLAY_METADATA="golden-metadata.display.yaml"
WORKING_METADATA="metadata.yaml"
WORKING_DISPLAY_METADATA="metadata.display.yaml"

if [[ -n $CURRENT_DIR ]]; then
  WORKING_FOLDER="$CURRENT_DIR/.working"
fi

if [ -d $WORKING_FOLDER ]
then
  rm -r -f $WORKING_FOLDER
fi

mkdir $WORKING_FOLDER && pushd $WORKING_FOLDER
kpt pkg get https://github.com/terraform-google-modules/terraform-google-cloud-storage.git/@v4.0.0 "./$BLUPRINT_FOLDER/"
../../../bin/cft blueprint metadata -d -p "$BLUPRINT_FOLDER/" -q

mkdir $GIT_FOLDER
cp "../$GOLDENS_FOLDER/$GOLDEN_METADATA" "$GIT_FOLDER/$WORKING_METADATA"
cp "../$GOLDENS_FOLDER/$GOLDEN_DISPLAY_METADATA" "$GIT_FOLDER/$WORKING_DISPLAY_METADATA"

pushd "$GIT_FOLDER"

# Confirm if the goldens are still valid with the blueprint schema
../../../../bin/cft blueprint metadata -v
rval=$?

if [ $rval -ne 0 ]; then
  echo "Error! Unable to validate the golden metadata(s)."
  exit $rval
fi

git init
git add .

cp "../$BLUPRINT_FOLDER/$WORKING_METADATA" "$WORKING_METADATA"
cp "../$BLUPRINT_FOLDER/$WORKING_DISPLAY_METADATA" "$WORKING_DISPLAY_METADATA"

git diff --exit-code --quiet
rval=$?

if [ $rval -eq 1 ]; then
  echo "Error! Generated metadata(s) do not match the golden(s)."
  git diff > diff.txt
  cat diff.txt
  exit $rval
elif [ $rval -gt 1 ]; then
  echo "Error occurred while comparaing metadata(s) to golden."
  exit $rval
fi

echo "Success: generated metadata(s) match the golden(s)."
