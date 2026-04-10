#!/usr/bin/env bash
set -euo pipefail

PROJECT="trichner-212015"
REGION="europe-west4"
SERVICE="heiss"

ENV_VARS=$(grep -v '^\s*#' .env.prod | grep -v '^\s*$' | paste -sd, -)

gcloud run deploy "$SERVICE" \
  --project "$PROJECT" \
  --region "$REGION" \
  --source . \
  --set-env-vars "$ENV_VARS" \
  --allow-unauthenticated
