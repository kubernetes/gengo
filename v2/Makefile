all:
	go build ./...

test:
	go test -race ./...

# We verify for the maximum version of the go directive as 1.20
# here because the oldest go directive that exists on our supported
# release branches in k/k is 1.20.
verify:
	./hack/verify-examples.sh
	./hack/verify-go-directive.sh 1.20
