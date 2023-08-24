//go:generate go run k8s.io/gengo/examples/deepcopy-gen -i k8s.io/gengo/examples/deepcopy-gen/output_tests/... -O zz_generated --go-header-file=../../../boilerplate/boilerplate.go.txt -o . --trim-path-prefix=k8s.io/gengo/examples/deepcopy-gen/output_tests
package output_tests
