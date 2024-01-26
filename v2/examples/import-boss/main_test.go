package main_test

import (
	"bytes"
	"context"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"testing"
)

const pkgName string = "k8s.io/gengo/v2/examples/import-boss"

func runImportBoss(inputPackages ...string) error {
	for i, v := range inputPackages {
		inputPackages[i] = pkg(v)
	}

	_, filename, _, _ := runtime.Caller(1)
	dir := path.Dir(filename)
	cmd := exec.CommandContext(context.Background(),
		"go", "run", pkgName, "--logtostderr", "--v=4",
		"-i", strings.Join(inputPackages, ","),
	)
	cmd.Dir = dir

	errBuf := &bytes.Buffer{}
	outBuf := &bytes.Buffer{}
	cmd.Stderr = errBuf
	cmd.Stdout = outBuf

	err := cmd.Run()
	print(errBuf.String())
	print(outBuf.String())
	return err
}

func pkg(name string) string {
	return pkgName + "/tests/" + name
}

type importBossTestCase struct {
	packageName string
	expectError bool
}

func TestRules(t *testing.T) {
	cases := []importBossTestCase{
		{
			packageName: "a",
			expectError: false,
		},
		{
			packageName: "b",
			expectError: true,
		},
		{
			packageName: "c",
			expectError: false,
		},
		{
			packageName: "nested",
			expectError: true,
		},
		{
			packageName: "nested/nested",
			expectError: false,
		},
		{
			packageName: "nested/nested/nested",
			expectError: true,
		},
		{
			packageName: "nested/nested/nested/inherit",
			expectError: true,
		},
	}

	for _, v := range cases {
		t.Run(v.packageName, func(t *testing.T) {
			err := runImportBoss("rules/" + v.packageName)
			if err != nil != v.expectError {
				t.Errorf("expected error: %v, returned error: %v", v.expectError, err != nil)
			}
		})
	}
}

func TestInverse(t *testing.T) {
	libPackages := []string{
		"inverse/lib",
		"inverse/lib/nonprod",
		"inverse/lib/private",
		"inverse/lib/public",
	}

	cases := []importBossTestCase{

		{
			packageName: "a",
			expectError: false,
		},
		{
			packageName: "b",
			expectError: true,
		},
		{
			packageName: "c",
			expectError: false,
		},
		{
			packageName: "d",
			expectError: true,
		},
		{
			packageName: "lib/quarantine",
			expectError: true,
		},
	}

	for _, v := range cases {
		t.Run(v.packageName, func(t *testing.T) {
			err := runImportBoss(append([]string{"inverse/" + v.packageName}, libPackages...)...)
			if err != nil != v.expectError {
				t.Errorf("expected error: %v, returned error: %v", v.expectError, err != nil)
			}
		})
	}
}
