#!/usr/bin/env bash
set -euo pipefail
trap 'echo "❌ Build failed at line $LINENO"' ERR

# =========================
# 参数校验
# =========================
if [[ $# -lt 2 ]]; then
  echo "Usage: $0 <GIT_BRANCH> <GIT_COMMIT>"
  exit 1
fi

GIT_BRANCH="$1"
GIT_COMMIT="$2"

# =========================
# 基本变量
# =========================
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CI_DIR="${PROJECT_ROOT}/ci"
OUTPUT_DIR="${CI_DIR}/dist"
BIN_NAME="go-demo"

GO_IMAGE="golang:1.25.5"
GOOS="${GOOS:-linux}"
GOARCH="${GOARCH:-amd64}"

BUILD_TIME="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"

mkdir -p "${OUTPUT_DIR}"

echo "▶ Building ${BIN_NAME}"
echo "  GOOS=${GOOS} GOARCH=${GOARCH}"
echo "  BRANCH=${GIT_BRANCH} COMMIT=${GIT_COMMIT}"

# =========================
# Docker Build
# =========================
docker run --rm \
  -e GOOS="${GOOS}" \
  -e GOARCH="${GOARCH}" \
  -e CGO_ENABLED=0 \
  -v "${PROJECT_ROOT}:/workspace" \
  -v "${HOME}/go/pkg/mod:/go/pkg/mod" \
  -v "${HOME}/.cache/go-build:/root/.cache/go-build" \
  -w /workspace \
  "${GO_IMAGE}" \
  bash -c "
    go mod download
    go test ./...
    go build \
      -trimpath \
      -ldflags \"-s -w \
        -X main.GitBranch=${GIT_BRANCH} \
        -X main.GitCommit=${GIT_COMMIT} \
        -X main.BuildTime=${BUILD_TIME}\" \
      -o ci/dist/${BIN_NAME}
  "

chmod +x "${OUTPUT_DIR}/${BIN_NAME}"
ls -lh "${OUTPUT_DIR}/${BIN_NAME}"

echo "✅ Build success"
