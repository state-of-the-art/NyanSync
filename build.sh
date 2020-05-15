#!/bin/sh -e

SRC_DIR=$(dirname "$0")
BUILD_DIR="${PWD}/build"

rm -rf "${BUILD_DIR}"
cp -a "${SRC_DIR}" "${BUILD_DIR}"

cd "${BUILD_DIR}"

# Do not download/generate assets if $1 is not empty
[ "$1" ] || scripts/assets.sh

reformat=$(gofmt -l .)
[ -z "${reformat}" ] || (echo "Please run 'go fmt': ${reformat}"; exit 1)

# Generate sources
go generate ./...

# Build & install the binary
go install -v ./cmd/NyanShare
