#!/bin/sh
# Test Dockerfile.scripts image: build and run without Postgres params.
# The scripts image must start with only DIAG_HTTP_STORAGE_HOST or DIAG_PV_MOUNT_PATH
# (no DIAG_POSTGRES_* or DIAG_DB_NAME). Run from repository root:
#   apps/dumps-collector/tests/scripts/test-dockerfile-scripts.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
COLLECTOR_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"
REPO_ROOT="$(cd "${COLLECTOR_DIR}/../.." && pwd)"
IMAGE_NAME="${1:-qubership-profiler-dumps-collector-scripts}"

log() { echo "[$(date +%FT%T%Z)][test-dockerfile-scripts] $*"; }

log "Building Dockerfile.scripts image: ${IMAGE_NAME}"
docker build -t "${IMAGE_NAME}" -f "${COLLECTOR_DIR}/Dockerfile.scripts" "${REPO_ROOT}"

log "Running container without any Postgres params (scripts image must not require them)"
# Use DIAG_HTTP_STORAGE_HOST so entrypoint does not start prf_dump_writer; only nginx runs.
# Container should start and nginx listen on 8080; we verify with a short run and curl or timeout.
CONTAINER_ID=$(docker run -d \
  -e DIAG_HTTP_STORAGE_HOST=http://backend:8080 \
  --name "dumps-scripts-test-$$" \
  "${IMAGE_NAME}")

cleanup() {
  docker rm -f "dumps-scripts-test-$$" 2>/dev/null || true
}
trap cleanup EXIT

log "Waiting for container to be ready..."
sleep 3

if docker exec "dumps-scripts-test-$$" sh -c "wget -q -O - http://127.0.0.1:8080/ || true" 2>/dev/null | head -1 | grep -q .; then
  log "PASS: Container responded on 8080 without Postgres params"
else
  # Container running without PG is enough for pass
  if docker inspect -f '{{.State.Running}}' "dumps-scripts-test-$$" 2>/dev/null | grep -q true; then
    log "PASS: Container is running (scripts image runs without PG params)"
  else
    log "FAIL: Container not running"
    docker logs "dumps-scripts-test-$$" 2>&1 || true
    exit 1
  fi
fi

log "All Dockerfile.scripts tests passed."
