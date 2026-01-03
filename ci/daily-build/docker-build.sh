#!/usr/bin/env bash
set -euo pipefail
trap 'echo "❌ Build failed at line $LINENO"' ERR

if [[ $# -lt 1 ]]; then
  echo "Usage: $0 <IMAGE_TAG>"
  exit 1
fi

IMAGE_TAG="$1"

APP_NAME="go-demo"
REGISTRY="registry.example.com/go-demo"
IMAGE_NAME="${REGISTRY}/${APP_NAME}:${IMAGE_TAG}"

echo "▶ Docker build: ${IMAGE_NAME}"

docker build -t ${IMAGE_NAME} -f ci/daily-build/Dockerfile .

docker push ${IMAGE_NAME}

echo "✅ Docker image pushed"
