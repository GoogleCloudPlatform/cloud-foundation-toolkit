#!/bin/bash
set -e
set -o pipefail

cat <<- EOF > ./options.yaml
options:
  inputMappings:
  - fieldName: Authorization
    location: HEADER
    value: >
      $.concat("Bearer ", $.googleOauth2AccessToken())
EOF

# Create the type-provider.
gcloud beta deployment-manager type-providers create cloudtasks \
      --api-options-file=options.yaml \
      --descriptor-url="https://cloudtasks.googleapis.com/\$discovery/rest?version=v2beta3"

exit 0
