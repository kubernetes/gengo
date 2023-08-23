TOOL=defaulter-gen

.PHONY : gen test verify git_verify_clean git_verify_no_change
test:
	cd _output_tests && go test ./...

gen:
	@go build -o /tmp/$(TOOL)
	@cd _output_tests && \
	PKGS=$$(go list ./...  | paste -sd' ' -) && \
	/tmp/$(TOOL) --logtostderr --v=4 \
		-i $$(echo $$PKGS | sed 's/ /,/g') \
		-O zz_generated \
		-o . \
		--trim-path-prefix=k8s.io/gengo/examples/defaulter-gen/output_tests \
		--go-header-file=../../../boilerplate/boilerplate.go.txt

git_verify_clean:
	@if ! git diff --quiet HEAD; then \
	    echo "FAIL: git client is not clean"; \
	    false; \
	fi

git_verify_no_change:
	@if ! git diff --quiet HEAD; then \
	    echo "FAIL: output files changed"; \
	    git diff; \
	    false; \
	fi

verify: git_verify_clean gen git_verify_no_change
