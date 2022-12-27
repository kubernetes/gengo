/*
Copyright 2016 The Kubernetes Authors.

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

package generators

import (
	"path/filepath"
	"testing"

	"k8s.io/gengo/types"
)

func TestRemoveLastDir(t *testing.T) {
	table := map[string]struct{ newPath, removedDir string }{
		"a/b/c": {"a/c", "b"},
	}
	for slashInput, expect := range table {
		input := filepath.FromSlash(slashInput)

		gotPath, gotRemoved := removeLastDir(input)
		if e, a := filepath.FromSlash(expect.newPath), gotPath; e != a {
			t.Errorf("%v: wanted %v, got %v", input, e, a)
		}
		if e, a := filepath.FromSlash(expect.removedDir), gotRemoved; e != a {
			t.Errorf("%v: wanted %v, got %v", input, e, a)
		}
	}
}

func TestInputIncludes(t *testing.T) {
	inputs := []string{"a", "a/b", "c/d/..."}

	if !inputIncludes(inputs, &types.Package{Path: "a"}) {
		t.Errorf("Expected 'a' to match")
	}
	if !inputIncludes(inputs, &types.Package{Path: "a/b"}) {
		t.Errorf("Expected 'a/b' to match")
	}
	if inputIncludes(inputs, &types.Package{Path: "a/b/c"}) {
		t.Errorf("Expected 'a/b/c' to not match")
	}
	if inputIncludes(inputs, &types.Package{Path: "c"}) {
		t.Errorf("Expected 'c' to not match")
	}
	if !inputIncludes(inputs, &types.Package{Path: "c/d"}) {
		t.Errorf("Expected 'c/d' to match")
	}
	if !inputIncludes(inputs, &types.Package{Path: "c/d/e"}) {
		t.Errorf("Expected 'c/d/e' to match via /... syntax")
	}
	if inputIncludes(inputs, &types.Package{Path: "z"}) {
		t.Errorf("Expected 'z' to not match")
	}
	if inputIncludes(inputs, &types.Package{Path: "c/z"}) {
		t.Errorf("Expected 'c/z' to not match")
	}
}
