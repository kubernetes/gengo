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

	for _, tc := range []struct {
		name     string
		original string
		alias    string
		warn     bool
	}{
		{
			name:     "different prose",
			original: "Original implements Tracker.",
			alias:    "Alias implements Tracker.\n\nAlias is preserved for backwards compatibility.",
		},
		{
			name:     "same directives with different prose",
			original: "Original implements Tracker.\n+k8s:deepcopy-gen=true",
			alias:    "Alias is a compatibility name.\n+k8s:deepcopy-gen=true",
		},
		{
			name:     "conflicting directives",
			original: "+k8s:deepcopy-gen=true",
			alias:    "+k8s:deepcopy-gen=false",
			warn:     true,
		},
		{
			name:     "directive only on original",
			original: "+custom:generate=true",
			alias:    "Alias is a compatibility name.",
			warn:     true,
		},
		{
			name:     "directive only on alias",
			original: "Original implements Tracker.",
			alias:    "+custom:generate=true",
			warn:     true,
		},
	} {
		for _, detached := range []bool{false, true} {
			for _, aliasFirst := range []bool{false, true} {
				name := tc.name
				if detached {
					name += "/detached"
				}
				if aliasFirst {
					name += "/alias-first"
				}
				t.Run(name, func(t *testing.T) {
					p := New()
					comment := func(s string) string {
						return "// " + strings.ReplaceAll(s, "\n", "\n// ") + "\n"
					}
					gap := ""
					if detached {
						gap = "\n"
					}
					source := "package p\n" + comment(tc.original) + gap + "type Original struct{}\n" +
						comment(tc.alias) + gap + "type Alias = Original\n"
					f, err := parser.ParseFile(p.fset, "comments.go", source, parser.ParseComments)
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
					original, alias := pkg.Scope().Lookup("Original"), pkg.Scope().Lookup("Alias")
					typ := &types.Type{GoType: original.Type()}
					first, second := original, alias
					if aliasFirst {
						first, second = alias, original
					}
					var logs bytes.Buffer
					klog.SetOutput(&logs)
					p.addCommentsToType(first, typ)
					p.addCommentsToType(second, typ)
					klog.Flush()
					if warned := strings.Contains(logs.String(), "Mismatched"); warned != tc.warn {
						t.Errorf("warning = %v, want %v; logs:\n%s", warned, tc.warn, logs.String())
					}
					got := typ.CommentLines
					if detached {
						got = typ.SecondClosestCommentLines
					}
					if want := strings.Split(tc.original, "\n"); !reflect.DeepEqual(got, want) {
						t.Errorf("comments = %q, want %q", got, want)
					}
				})
			}
		}
	}
}
