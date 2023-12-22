// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package is copied from Go library net.
// https://golang.org/src/go/ast/ast.go
// The original private function isDirective
// is exported as public function.

package ast

import "strings"

// IsDirective reports whether c is a comment directive.
// This code is also in go/printer.
func IsDirective(c string) bool {
	// "//line " is a line directive.
	// "//extern " is for gccgo.
	// "//export " is for cgo.
	// (The // has been removed.)
	if strings.HasPrefix(c, "line ") || strings.HasPrefix(c, "extern ") || strings.HasPrefix(c, "export ") {
		return true
	}

	// "//[a-z0-9]+:[a-z0-9]"
	// (The // has been removed.)
	colon := strings.Index(c, ":")
	if colon <= 0 || colon+1 >= len(c) {
		return false
	}
	for i := 0; i <= colon+1; i++ {
		if i == colon {
			continue
		}
		b := c[i]
		if !('a' <= b && b <= 'z' || '0' <= b && b <= '9') {
			return false
		}
	}
	return true
}
