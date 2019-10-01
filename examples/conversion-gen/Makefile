TOOL=conversion-gen

test:
	@if ! git diff --quiet HEAD .; then \
	    echo "FAIL: git client is not clean"; \
	    false; \
	fi
	@go build -o /tmp/$(TOOL)
	@PKGS=$$(find . -name doc.go -exec grep --color=never -l '+k8s:conversion-gen' {} \; | xargs dirname | awk '{print "k8s.io/gengo/examples/conversion-gen" substr($$1,2); }') ; \
	/tmp/$(TOOL) --logtostderr --v=4 -i $$(echo $$PKGS | sed 's/ /,/g') -O zz_generated
	@if ! git diff --quiet HEAD .; then \
	    echo "FAIL: output files changed"; \
	    git diff; \
	    false; \
	fi

