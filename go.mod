module k8s.io/gengo

go 1.18

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/google/go-cmp v0.4.0
	github.com/google/gofuzz v1.1.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/tools v0.0.0-20200505023115-26f46d2f7ef8
	k8s.io/klog/v2 v2.2.0
	sigs.k8s.io/yaml v1.2.0
)

require (
	github.com/go-logr/logr v0.2.0 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace (
	golang.org/x/sys => golang.org/x/sys v0.0.0-20190813064441-fde4db37ae7a // pinned to release-branch.go1.13
	golang.org/x/tools => golang.org/x/tools v0.0.0-20190821162956-65e3620a7ae7 // pinned to release-branch.go1.13
)
