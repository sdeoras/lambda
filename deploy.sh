#!/usr/bin/env bash

go mod vendor
zip -r payload.zip lambda.go api jwt email vendor
rm -rf vendor
gsutil cp payload.zip gs://${CLOUD_FUNCTIONS_BUCKET}
rm -rf payload.zip

gcloud functions deploy lambda \
    --region=us-central1 \
    --trigger-http \
    --entry-point=Lambda \
    --runtime go111 \
    --set-env-vars=JWT_SECRET_KEY="${JWT_SECRET_KEY}" \
    --set-env-vars=GOOGLE_GCF_DOMAIN="${GOOGLE_GCF_DOMAIN}" \
    --set-env-vars=GOOGLE_CLIENT_ID="${GOOGLE_CLIENT_ID}" \
    --set-env-vars=GOOGLE_CLIENT_SECRET="${GOOGLE_CLIENT_SECRET}" \
    --set-env-vars=GCLOUD_PROJECT_NAME="${GCLOUD_PROJECT_NAME}" \
    --set-env-vars=SENDGRID_API_KEY="${SENDGRID_API_KEY}" \
    --set-env-vars=EMAIL_FROM_NAME="${EMAIL_FROM_NAME}" \
    --set-env-vars=EMAIL_FROM_EMAIL="${EMAIL_FROM_EMAIL}" \
    --set-env-vars=EMAIL_TO_NAME="${EMAIL_TO_NAME}" \
    --set-env-vars=EMAIL_TO_EMAIL="${EMAIL_TO_EMAIL}" \
    --source=gs://${CLOUD_FUNCTIONS_BUCKET}/payload.zip
