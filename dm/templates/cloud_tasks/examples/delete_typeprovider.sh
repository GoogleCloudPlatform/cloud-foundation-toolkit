#!/bin/bash
set -e
set -o pipefail

# Delete the type provider.
gcloud beta deployment-manager type-providers delete cloudtasks -q

exit 0
