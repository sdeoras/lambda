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
    --source=gs://${CLOUD_FUNCTIONS_BUCKET}/payload-${NAME}.zip

# do health checks on all services
echo "Performing health checks..."
for SERVICE in ${SERVICES[@]}; do
    echo "${NAME}/${SERVICE}" \
    `curl "https://${GOOGLE_GCF_DOMAIN}/${NAME}/health/?format=mesg&service=${SERVICE}" 2>/dev/null`
done
