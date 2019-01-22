#!/usr/bin/env bash

gcloud functions deploy router \
    --region=us-central1 \
    --trigger-http \
    --runtime=go111 \
    --set-env-vars=JWT_SECRET_KEY=${JWT_SECRET_KEY} \
    --entry-point=Route