all:
	go build ./...
	go -C v2 build ./...

test:
	go test -race ./...
	go -C v2 test -race ./...
	(cd v2 && ./hack/verify-examples.sh)

verify: verify-go-directive
	./hack/verify-examples.sh

# We set the maximum version of the go directive as 1.20 here
# because the oldest go directive that exists on our supported
# release branches in k/k is 1.20.
.PHONY: verify-go-directive
verify-go-directive:
	./hack/verify-go-directive.sh -g 1.20
