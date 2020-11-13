TOOL=defaulter-gen

.PHONY : gen test verify git_verify_clean git_verify_no_change
test:
	go test ./_output_tests/marker/...

gen:
	@go build -o /tmp/$(TOOL)
	@PKGS=$$(cd _output_tests; go list ./...  | paste -sd' ' -); \
	/tmp/$(TOOL) --logtostderr --v=4 -i $$(echo $$PKGS | sed 's/ /,/g') -O zz_generated

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
