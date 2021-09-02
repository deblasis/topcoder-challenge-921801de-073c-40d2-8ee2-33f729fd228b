#!/bin/bash

REPO_ROOT="${REPO_ROOT:-$(cd "$(dirname "$0")/.." && pwd)}"

PROTO_FILE=${1:-"auth_dbsvc.proto"}
PB_PATH="${REPO_ROOT}/services/auth_dbsvc/pb"

echo "Generating pb files for ${PROTO_FILE} service"
protoc -I="${PB_PATH}"  "${PB_PATH}/${PROTO_FILE}" --go_out=plugins=grpc:.

