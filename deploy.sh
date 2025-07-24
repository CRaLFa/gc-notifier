#!/bin/bash

set -eu

deploy_webhook () {
    gcloud --project=${1} functions deploy gc-webhook \
        --gen2 \
        --runtime=go123 \
        --region=asia-northeast1 \
        --source=./webhook \
        --entry-point=GetGroupID \
        --trigger-http \
        --allow-unauthenticated \
        --service-account=cloud-run@${1}.iam.gserviceaccount.com \
        --env-vars-file=./.env.yaml
}

deploy_notifier () {
    gcloud --project=${1} functions deploy gc-notifier \
        --gen2 \
        --runtime=go123 \
        --region=asia-northeast1 \
        --source=./notifier \
        --entry-point=PostNotification \
        --trigger-http \
        --no-allow-unauthenticated \
        --service-account=cloud-run@${1}.iam.gserviceaccount.com \
        --env-vars-file=./.env.yaml
}

main () {
    cd "$(dirname "$0")"
    local project_id=$(yq -r '.PROJECT_ID' < .env.yaml)

    [[ $# -lt 1 || "$1" = 'webhook' ]] && {
        echo -e "Deploying webhook...\n"
        deploy_webhook "$project_id"
        echo
    }
    [[ $# -lt 1 || "$1" = 'notifier' ]] && {
        echo -e "Deploying notifier...\n"
        deploy_notifier "$project_id"
    }
}

main "$@"
