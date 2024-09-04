all:
	go build ./...
	make -C v2 all

test:
	go test -race ./...
	go -C ./examples/defaulter-gen/output_tests test ./...
	make -C v2 test

# We verify for the maximum version of the go directive as 1.20
# here because the oldest go directive that exists on our supported
# release branches in k/k is 1.20.
verify:
	./hack/verify-examples.sh
	./hack/verify-go-directive.sh 1.20
	make -C v2 verify
