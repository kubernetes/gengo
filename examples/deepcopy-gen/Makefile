TOOL=deepcopy-gen

test:
	@if ! git diff --quiet HEAD; then \
	    echo "FAIL: git client is not clean"; \
	    false; \
	fi
	@go build -o /tmp/$(TOOL)
	@PKGS=$$(go list ./output_tests/...  | paste -sd' ' -); \
	/tmp/$(TOOL) --logtostderr --v=4 -i $$(echo $$PKGS | sed 's/ /,/g') -O zz_generated --go-header-file ../../boilerplate/boilerplate.go.txt
	@if ! git diff --quiet HEAD; then \
	    echo "FAIL: output files changed"; \
	    git diff; \
	    false; \
	fi
