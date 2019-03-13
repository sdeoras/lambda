#!/usr/bin/env bash

export NAME="gan"
export TOOL="gangen"

if [[ -n "$1" ]]; then
  export NAME="${NAME}-$1"
else
  echo "Deploying to production. To deploy to staging use command: ./deploy.sh staging"
fi

# check for dependencies that are not part of git distribution
if [[ ! -f "src/bin/src/${TOOL}/a.out" ]]; then
    echo "you need to build a.out for linux/amd64. pl. run src/bin/src/${TOOL}/deploy.sh"
    exit 1
fi

if [[ ! -f "src/bin/src/${TOOL}/lib/libtensorflow.so" ]]; then
    echo "pl. download src/bin/src/${TOOL}/lib/libtensorflow.so for Linux for TF v1.12.0"
    exit 1
fi

if [[ ! -f "src/bin/src/${TOOL}/lib/libtensorflow_framework.so" ]]; then
    echo "pl. download src/bin/src/${TOOL}/lib/libtensorflow_framework.so for Linux for TF v1.12.0"
    exit 1
fi

export PROJECT=`gcloud config list 2>/dev/null | grep ^project | awk '{print $3}'`
export REGION=`gcloud config list 2>/dev/null | grep ^region | awk '{print $3}'`
export GOOGLE_GCF_DOMAIN="${REGION}-${PROJECT}.cloudfunctions.net"
export CLOUD_FUNCTIONS_BUCKET="${PROJECT}-gcf"

go mod vendor
zip -r payload-${NAME}.zip lambda.go src vendor
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
    --set-env-vars=CLOUD_FUNCTIONS_BUCKET="${CLOUD_FUNCTIONS_BUCKET}" \
    --source=gs://${CLOUD_FUNCTIONS_BUCKET}/payload-${NAME}.zip
