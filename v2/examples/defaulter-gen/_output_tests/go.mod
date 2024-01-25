// Package must be named output_tests rather than _output tests to avoid
// ambiguous import conflict with main gengo package
module k8s.io/gengo/v2/examples/defaulter-gen/output_tests

go 1.21.3

require (
	github.com/google/go-cmp v0.5.9
	k8s.io/apimachinery v0.28.0
	k8s.io/gengo/v2 v2.0.0
)

require (
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/tools v0.16.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/klog/v2 v2.100.1 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)

replace k8s.io/gengo/v2 => ../../../
