set -e

BUILD_ID=$(gcloud beta builds triggers run gcp-org---terraform-apply --branch=production --project prj-b-cicd-9f11 --format="value(metadata.build.id)")
gcloud beta builds log $BUILD_ID --project prj-b-cicd-9f11 --stream

repos=("gcp-environments" "gcp-networks" "gcp-projects")
envs=("development" "non-production" "production")

for repo in ${repos[@]}; do
    TRIGGER_NAME="${repo}---terraform-apply"
    for env in ${envs[@]}; do
        echo "Reconciling ${repo} in ${env}"
        BUILD_ID=$(gcloud beta builds triggers run ${TRIGGER_NAME} --branch=${env} --project prj-b-cicd-9f11 --format="value(metadata.build.id)")
        gcloud beta builds log $BUILD_ID --project prj-b-cicd-9f11 --stream
    done
done