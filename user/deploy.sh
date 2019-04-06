#!/usr/bin/env bash

export NAME="user"
SERVICES=("register" "query")

if [[ -n "$1" ]]; then
  export NAME="${NAME}-$1"
else
  echo "Deploying to production. To deploy to staging use command: ./deploy.sh staging"
fi

export PROJECT=`gcloud config list 2>/dev/null | grep ^project | awk '{print $3}'`
export REGION=`gcloud config list 2>/dev/null | grep ^region | awk '{print $3}'`
export GOOGLE_GCF_DOMAIN="${REGION}-${PROJECT}.cloudfunctions.net"
export CLOUD_FUNCTIONS_BUCKET="${PROJECT}-gcf"

go mod vendor 2> /dev/null
zip -r payload-${NAME}.zip lambda.go src vendor
gsutil cp payload-${NAME}.zip gs://${CLOUD_FUNCTIONS_BUCKET}
rm -rf vendor payload-${NAME}.zip

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

# do health checks on all services
echo "Performing health checks..."
for SERVICE in ${SERVICES[@]}; do
    echo "${NAME}/${SERVICE}" \
    `curl "https://${GOOGLE_GCF_DOMAIN}/${NAME}/health/?format=mesg&service=${SERVICE}" 2>/dev/null`
done