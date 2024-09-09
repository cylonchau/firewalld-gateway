#!/bin/bash

PROJECT_ROOT=$(cd $(dirname ${BASH_SOURCE[0]})/..; pwd)
OUT_DIR="_output"
GO_CMD=`which go`
BINARY_NAME=$1

[ -d ${PROJECT_ROOT}/${OUT_DIR} ] || mkdir -pv ${PROJECT_ROOT}/${OUT_DIR}
cp -a ${PROJECT_ROOT}/dist server/app/router/ && \
    cd ${PROJECT_ROOT}/cmd && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${GO_CMD} build -ldflags "-s -w" -o ${PROJECT_ROOT}/${OUT_DIR}/${BINARY_NAME} main.go && \
    rm -fr ${PROJECT_ROOT}/server/app/router/dist