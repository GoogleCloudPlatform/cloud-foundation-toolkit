#!/bin/bash

IFS=$'\n\t'
set -eou pipefail
MODE="DRYRUN"

if [[ "$#" -lt 1 || "${1}" == '-h' || "${1}" == '--help' ]]; then
  cat >&2 <<"EOF"
cft-image-cleanup.sh cleans up untagged cft-dev-tool images.
USAGE:
  cft-image-cleanup.sh REPOSITORY [DELETE]
  e.g. $ ./cft-image-cleanup.sh gcr.io/cloud-foundation-cicd/cft/developer-tools DELETE
  would delete all image digests that do not have a tag in the gcr.io/cloud-foundation-cicd/cft/developer-tools repository
EOF
  exit 1
fi

main(){
  local C=0
  IMAGE="${1}"
  for digest in $(gcloud container images list-tags "${IMAGE}" --limit=999999 --sort-by=TIMESTAMP \
    --format='get(digest)' --filter='-tags:*'); do
    if [[ "$MODE" == "DRYRUN" ]]; then
      echo "to delete: $digest"
    elif [[ "$MODE" == "DELETE" ]]; then
      (
        set -x
        gcloud container images delete -q --force-delete-tags "${IMAGE}@${digest}"
      )
    fi
    (( C=C+1 ))
  done
  echo "Deleted ${C} images in ${IMAGE}." >&2
}

if [[ "$#" -eq 1 ]]; then
  echo ">>> executing in DRY RUN mode; use the DELETE arg for deleting the images <<<"
elif [[ "$#" -eq 2 && "${2}" == 'DELETE' ]]; then
  MODE="DELETE"
fi
main "${1}"
