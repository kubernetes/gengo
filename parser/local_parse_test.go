/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package parser

import (
	"go/ast"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestImportBuildPackage(t *testing.T) {
	// get our original dir to restore
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Chdir(dir)
	}()
	// switch into the fake dependent module which knows where the fake module is
	os.Chdir("../testdata/dependent")

	b := New()
	if _, err := b.importBuildPackage("fake/dep"); err != nil {
		t.Fatal(err)
	}
	if _, ok := b.buildPackages["fake/dep"]; !ok {
		t.Errorf("missing expected, but got %v", b.buildPackages)
	}

	if len(b.buildPackages) > 1 {
		// this would happen if the canonicalization failed to normalize the path
		// you'd get a k8s.io/gengo/vendor/fake/dep key too
		t.Errorf("missing one, but got %v", b.buildPackages)
	}
}

func TestIsErrPackageNotFound(t *testing.T) {
	b := New()
	if _, err := b.importBuildPackage("fake/empty"); !isErrPackageNotFound(err) {
		t.Errorf("expected error like %s, but got %v", regexErrPackageNotFound.String(), err)
	}
}

func TestCanonicalizeImportPath(t *testing.T) {
	tcs := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "passthrough",
			input:  "github.com/foo/bar",
			output: "github.com/foo/bar",
		},
		{
			name:   "simple",
			input:  "github.com/foo/vendor/k8s.io/kubernetes/pkg/api",
			output: "k8s.io/kubernetes/pkg/api",
		},
		{
			name:   "deeper",
			input:  "github.com/foo/bar/vendor/k8s.io/kubernetes/pkg/api",
			output: "k8s.io/kubernetes/pkg/api",
		},
	}

	for _, tc := range tcs {
		actual := canonicalizeImportPath(tc.input)
		if string(actual) != tc.output {
			t.Errorf("%v: expected %q got %q", tc.name, tc.output, actual)
		}
	}
}

func TestExtractDirecrives(t *testing.T) {
	var tests = []struct {
		commentLines []string
		directives   []string
	}{
		{
			[]string{"// foo", "//   ", "//", "//", "//   bar   "},
			nil,
		},
		{
			[]string{"// foo", "/* bar */"},
			nil,
		},
		{
			[]string{"// foo", "// bar", "// notdirective:baz"},
			nil,
		},
		{
			[]string{"// foo", "/* notdirective0:baz */", "/*notdirective1:baz */"},
			nil,
		},
		{
			[]string{"// foo", "//go:noinline", "// bar", "// notdirective:baz", "//lint123:ignore"},
			[]string{"go:noinline", "lint123:ignore"},
		},
		{
			[]string{"// foo", "//go:noinline", "/* notdirective0:baz */", "/*notdirective1:baz */", "//lint123:ignore"},
			[]string{"go:noinline", "lint123:ignore"},
		},
	}

	for i, tt := range tests {
		list := make([]*ast.Comment, len(tt.commentLines))
		for i, line := range tt.commentLines {
			list[i] = &ast.Comment{Text: line}
		}

		directives := extractDirectives(&ast.CommentGroup{List: list})
		if diff := cmp.Diff(directives, tt.directives); diff != "" {
			t.Errorf(
				"Case: %d\nWanted, got:\n%v\n-----\n%v\nDiff:\n%s",
				i, tt.directives, directives, diff,
			)
		}
	}
}
