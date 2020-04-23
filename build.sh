#!/bin/sh -e

SRC_DIR=$(dirname "$0")
BUILD_DIR="${PWD}/build"

rm -rf "${BUILD_DIR}"
cp -a "${SRC_DIR}" "${BUILD_DIR}"

cd "${BUILD_DIR}"

[ "$1" ] || scripts/assets.sh

# Generate sources
go generate ./...

# Build & install the binary
go install -v .
