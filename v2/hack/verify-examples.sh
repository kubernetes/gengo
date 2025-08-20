#!/usr/bin/env bash

# hack/verify-examples.sh
# Pre-submit script to verify:
#   1.) committed generated code matches source
#   2.) code generation tests pass

# Exit immediately if any command fails
set -e

VERBOSE=${VERBOSE:-0}

# Make sure we run from the v2 root.
V2_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
cd "${V2_ROOT}"

# Ensure all files are committed
if ! git diff --quiet HEAD; then
    echo "FAIL: git client is not clean"
    exit 1
fi

echo "Running ./examples/tracer"
rm ./examples/tracer/testdata/simple/out.txt
go run ./examples/tracer -v"${VERBOSE}" ./examples/tracer/testdata/simple > ./examples/tracer/testdata/simple/out.txt
if ! git diff --quiet HEAD; then
    echo "FAIL: output files changed"
    git diff
    exit 1
fi

echo "Running ./examples/kilroy"
rm ./examples/kilroy/testdata/simple/generated.kilroy.go
go run ./examples/kilroy/ -v"${VERBOSE}" ./examples/kilroy/testdata/simple
if ! git diff --quiet HEAD; then
    echo "FAIL: output files changed"
    git diff
    exit 1
fi

echo "Running ./examples/pointuh"
rm -rf ./examples/pointuh/testdata/results
go run ./examples/pointuh/ -v"${VERBOSE}" \
    --output-dir ./examples/pointuh/testdata/results \
    --output-pkg k8s.io/gengo/examples/pointuh/testdata/results \
    ./examples/pointuh/testdata/simple
if ! git diff --quiet HEAD; then
    echo "FAIL: output files changed"
    git diff
    exit 1
fi
