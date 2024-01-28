#!/usr/bin/env bash

# hack/verify-examples.sh
# Pre-submit script to verify:
#   1.) committed generated code matches source
#   2.) code generation tests pass

# Exit immediately if any command fails
set -e

# Silence pushd/popd
pushd () {
    command pushd "$@" > /dev/null
}

popd () {
    command popd "$@" > /dev/null
}

# Ensure all files are committed
if ! git diff --quiet HEAD; then
    echo "FAIL: git client is not clean"
    exit 1
fi

echo "Running tests..."
go test ./examples/...

rm ./examples/tracer/testdata/simple/out.txt
go run ./examples/tracer ./examples/tracer/testdata/simple > ./examples/tracer/testdata/simple/out.txt
if ! git diff --quiet HEAD; then
    echo "FAIL: output files changed"
    git diff
    exit 1
fi

rm ./examples/kilroy/testdata/simple/generated.kilroy.go
go run ./examples/kilroy/ ./examples/kilroy/testdata/simple
if ! git diff --quiet HEAD; then
    echo "FAIL: output files changed"
    git diff
    exit 1
fi

rm -rf ./examples/pointuh/testdata/results
go run ./examples/pointuh/ \
    --output-dir ./examples/pointuh/testdata/results \
    --output-pkg k8s.io/gengo/examples/pointuh/testdata/results \
    ./examples/pointuh/testdata/simple
if ! git diff --quiet HEAD; then
    echo "FAIL: output files changed"
    git diff
    exit 1
fi
