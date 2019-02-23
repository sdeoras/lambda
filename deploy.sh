#!/usr/bin/env bash

export NAME="lambda"
export PROJECT=`gcloud config list 2>/dev/null | grep ^project | awk '{print $3}'`
export REGION=`gcloud config list 2>/dev/null | grep ^region | awk '{print $3}'`
export GOOGLE_GCF_DOMAIN="${REGION}-${PROJECT}.cloudfunctions.net"
export CLOUD_FUNCTIONS_BUCKET="${PROJECT}-gcf"

go mod vendor
zip -r payload-${NAME}.zip lambda.go bin api jwt email infer vendor
rm -rf vendor
gsutil cp payload-${NAME}.zip gs://${CLOUD_FUNCTIONS_BUCKET}
rm -rf payload-${NAME}.zip

gcloud functions deploy ${NAME} \
    --region=${REGION} \
    --trigger-http \
    --entry-point=Lambda \
    --runtime=go111 \
    --memory=2048MB \
    --set-env-vars=JWT_SECRET_KEY="${JWT_SECRET_KEY}" \
    --set-env-vars=GOOGLE_GCF_DOMAIN="${GOOGLE_GCF_DOMAIN}" \
    --set-env-vars=GOOGLE_CLIENT_ID="${GOOGLE_CLIENT_ID}" \
    --set-env-vars=GOOGLE_CLIENT_SECRET="${GOOGLE_CLIENT_SECRET}" \
    --set-env-vars=GCLOUD_PROJECT_NAME="${PROJECT}" \
    --set-env-vars=SENDGRID_API_KEY="${SENDGRID_API_KEY}" \
    --set-env-vars=EMAIL_FROM_NAME="${EMAIL_FROM_NAME}" \
    --set-env-vars=EMAIL_FROM_EMAIL="${EMAIL_FROM_EMAIL}" \
    --set-env-vars=EMAIL_TO_NAME="${EMAIL_TO_NAME}" \
    --set-env-vars=EMAIL_TO_EMAIL="${EMAIL_TO_EMAIL}" \
    --set-env-vars=LAMBDA_BUCKET="${LAMBDA_BUCKET}" \
    --source=gs://${CLOUD_FUNCTIONS_BUCKET}/payload-${NAME}.zip
