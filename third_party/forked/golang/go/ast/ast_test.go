// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package is copied from Go library net.
// https://golang.org/src/go/ast/ast.go
// The original private function isDirective
// is exported as public function.

package ast

import "testing"

var isDirectiveTests = []struct {
	in string
	ok bool
}{
	{"abc", false},
	{"go:inline", true},
	{"Go:inline", false},
	{"go:Inline", false},
	{":inline", false},
	{"lint:ignore", true},
	{"lint:1234", true},
	{"1234:lint", true},
	{"go: inline", false},
	{"go:", false},
	{"go:*", false},
	{"go:x*", true},
	{"export foo", true},
	{"extern foo", true},
	{"expert foo", false},
}

func TestIsDirective(t *testing.T) {
	for _, tt := range isDirectiveTests {
		if ok := IsDirective(tt.in); ok != tt.ok {
			t.Errorf("isDirective(%q) = %v, want %v", tt.in, ok, tt.ok)
		}
	}
}
