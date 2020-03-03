#!/bin/bash
# this script assumes we have a list of fields as tab separated values in the
# directory where we want to create setters
# read directory into variable
dir=$(echo $1 | cut -f1 -d/);
# read fields from file
while read p; do
  # split fields into field and type
  IFS="	";
  read -r -a arr <<< $p;
  # change underscores to camel case
  arr[0]=$(echo ${arr[0]} | sed -e 's/_\(.\)/\U\1/g');
  # change bool to boolean
  [[ ${arr[2]} == "bool" ]] && arr[2]="boolean";
  # change empty string to "''"
  [[ ${arr[3]} == '' ]] && arr[3]="''";
  # get value from yaml
  arr[3]=$(~/.local/bin/yq .metadata.${arr[0]} $1 -r);
  if [ "${arr[3]}" = "" ]; then
    kpt cfg create-setter $dir ${arr[0]} '' --type \"${arr[2]}\" --field ${arr[0]} --description \"${arr[1]}\";
  else
    kpt cfg create-setter $dir ${arr[0]} ${arr[3]} --type \"${arr[2]}\" --field ${arr[0]} --description \"${arr[1]}\";
  fi
done <$dir/fields
