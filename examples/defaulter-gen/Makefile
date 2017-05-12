TOOL=defaulter-gen

test:
	@if ! git diff --quiet HEAD; then \
	    echo "FAIL: git client is not clean"; \
	    false; \
	fi
	@go build -o /tmp/$(TOOL)
	/tmp/$(TOOL) --logtostderr --v=6 -i ./output_tests/wholepkg -O zz_generated
	@if ! diff -q ./output_tests/wholepkg/zz_generated.expected ./output_tests/wholepkg/zz_generated.go; then \
		echo "FAIL: output files changed"; \
		diff ./output_tests/wholepkg/zz_generated.expected ./output_tests/wholepkg/zz_generated.go; \
		false; \
	fi
	rm ./output_tests/wholepkg/zz_generated.go
