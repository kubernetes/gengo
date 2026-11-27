/*
Copyright 2026 The Kubernetes Authors.

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
	"bytes"
	"flag"
	"go/ast"
	"go/parser"
	gotypes "go/types"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"k8s.io/gengo/v2/types"
	"k8s.io/klog/v2"
)

func TestAliasCommentWarnings(t *testing.T) {
	// klog output is global, so these tests must not run in parallel.
	state := klog.CaptureState()
	defer state.Restore()
	klog.LogToStderr(false)
	flags := flag.NewFlagSet("klog", flag.ContinueOnError)
	klog.InitFlags(flags)
	if err := flags.Set("stderrthreshold", "FATAL"); err != nil {
		t.Fatal(err)
	}

	variants := []struct {
		name       string
		original   string
		alias      string
		aliasFirst bool
	}{
		{name: "attached/original-first", original: "Original", alias: "Alias"},
		{name: "attached/alias-first", original: "Original", alias: "Alias", aliasFirst: true},
		{name: "detached/original-first", original: "DetachedOriginal", alias: "DetachedAlias"},
		{name: "detached/alias-first", original: "DetachedOriginal", alias: "DetachedAlias", aliasFirst: true},
	}
	for _, tc := range []struct {
		fixture string
		warn    bool
	}{
		{fixture: "different-prose"},
		{fixture: "same-directives"},
		{fixture: "conflicting-directives", warn: true},
		{fixture: "directive-only-on-original", warn: true},
		{fixture: "directive-only-on-alias", warn: true},
	} {
		for _, variant := range variants {
			t.Run(tc.fixture+"/"+variant.name, func(t *testing.T) {
				p, scope := loadAliasCommentFixture(t, tc.fixture)
				original, alias := scope.Lookup(variant.original), scope.Lookup(variant.alias)
				first, second := original, alias
				if variant.aliasFirst {
					first, second = alias, original
				}

				// Processing an alias must preserve the original's comments.
				want := &types.Type{}
				if typeAliasEnabled {
					p.addCommentsToType(original, want)
				} else {
					// Without alias information, the last declaration wins.
					p.addCommentsToType(second, want)
				}

				var logs bytes.Buffer
				klog.SetOutput(&logs)
				got := &types.Type{GoType: original.Type()}
				p.addCommentsToType(first, got)
				p.addCommentsToType(second, got)
				klog.Flush()

				if warned := strings.Contains(logs.String(), "Mismatched"); warned != tc.warn {
					t.Errorf("warning = %v, want %v; logs:\n%s", warned, tc.warn, logs.String())
				}
				if !reflect.DeepEqual(got.CommentLines, want.CommentLines) {
					t.Errorf("attached comments = %q, want %q", got.CommentLines, want.CommentLines)
				}
				if !reflect.DeepEqual(got.SecondClosestCommentLines, want.SecondClosestCommentLines) {
					t.Errorf("detached comments = %q, want %q", got.SecondClosestCommentLines, want.SecondClosestCommentLines)
				}
			})
		}
	}
}

func loadAliasCommentFixture(t *testing.T, name string) (*Parser, *gotypes.Scope) {
	t.Helper()
	p := New()
	filename := filepath.Join("testdata", "alias-comments", name, "types.go")
	f, err := parser.ParseFile(p.fset, filename, nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	pkg, err := new(gotypes.Config).Check("p", p.fset, []*ast.File{f}, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, group := range f.Comments {
		pos := p.fset.Position(group.End())
		p.endLineToCommentGroup[fileLine{pos.Filename, pos.Line}] = group
	}
	return p, pkg.Scope()
}
