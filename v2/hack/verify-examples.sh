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

echo "Removing generated code"

# Defaulter-gen and deepcopy-gen both generate types of this format
find ./examples -name "zz_generated.go" -type f -delete

echo "Generating example output..."
go generate ./examples/...
pushd ./examples/defaulter-gen/_output_tests; go generate ./...; popd

# If there are any differences with committed files, fail
if ! git diff --quiet HEAD; then
    echo "FAIL: output files changed"
    git diff
    exit 1
fi

echo "Running tests..."
go test ./examples/...
go run ./examples/import-boss/main.go -i $(go list k8s.io/gengo/v2/... | grep -v import-boss/tests | paste -sd',' -) --verify-only
pushd ./examples/defaulter-gen/_output_tests; go test ./...; popd

rm ./examples/tracer/testdata/simple/out.txt
go run ./examples/tracer -i ./examples/tracer/testdata/simple > ./examples/tracer/testdata/simple/out.txt
if ! git diff --quiet HEAD; then
    echo "FAIL: output files changed"
    git diff
    exit 1
fi

rm ./examples/kilroy/testdata/simple/generated.kilroy.go
go run ./examples/kilroy/ -i ./examples/kilroy/testdata/simple -o .
if ! git diff --quiet HEAD; then
    echo "FAIL: output files changed"
    git diff
    exit 1
fi
