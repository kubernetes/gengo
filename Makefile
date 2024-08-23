all:
	go build ./...
	go -C v2 build ./...

test:
	go test -race ./...
	go -C v2 test -race ./...
	(cd v2 && ./hack/verify-examples.sh)
	(cd ./examples/defaulter-gen/_output_tests && go test -v -race ./...)

# We verify for the maximum version of the go directive as 1.20
# here because the oldest go directive that exists on our supported
# release branches in k/k is 1.20.
verify:
	./hack/verify-examples.sh
	./hack/verify-go-directive.sh 1.20
