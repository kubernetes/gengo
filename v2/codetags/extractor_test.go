/*
Copyright 2025 The Kubernetes Authors.

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

package codetags

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExtract(t *testing.T) {
	k8sGroup := "k8s"
	emptyGroup := ""
	tests := []struct {
		name   string
		marker string
		group  *string
		lines  []string
		want   map[TagIdentifier][]string
	}{
		{
			name:   "example from documentation",
			marker: "+",
			group:  &k8sGroup,
			lines: []string{
				"Comment line without marker",
				"+k8s:required",
				"+listType=set",
				"+k8s:format=k8s-long-name",
			},
			want: map[TagIdentifier][]string{
				{Group: "k8s", Name: "required"}: {"required"},
				{Group: "k8s", Name: "format"}:   {"format=k8s-long-name"},
			},
		},
		{
			name:   "empty lines",
			marker: "+",
			group:  &k8sGroup,
			lines:  []string{},
			want:   map[TagIdentifier][]string{},
		},
		{
			name:   "no matching lines",
			marker: "+",
			group:  &k8sGroup,
			lines: []string{
				"Comment line without marker",
				"Another comment line",
			},
			want: map[TagIdentifier][]string{},
		},
		{
			name:   "different marker",
			marker: "@",
			group:  &k8sGroup,
			lines: []string{
				"Comment line without marker",
				"@k8s:required",
				"@validation:required",
				"+k8s:format=k8s-long-name",
			},
			want: map[TagIdentifier][]string{
				{Group: "k8s", Name: "required"}: {"required"},
			},
		},
		{
			name:   "empty group",
			marker: "+",
			group:  &emptyGroup,
			lines: []string{
				"+k8s:required",
				"+required",
				"+format=special",
			},
			want: map[TagIdentifier][]string{
				{Group: "", Name: "required"}: {"required"},
				{Group: "", Name: "format"}:   {"format=special"},
			},
		},
		{
			name:   "no group",
			marker: "+",
			group:  nil,
			lines: []string{
				"+k8s:required",
				"+validation:required",
				"+validation:format=special",
			},
			want: map[TagIdentifier][]string{
				{Group: "k8s", Name: "required"}:        {"required"},
				{Group: "validation", Name: "required"}: {"required"},
				{Group: "validation", Name: "format"}:   {"format=special"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Extract(tt.marker, tt.group, tt.lines)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractAndParseTagsWithArgs(t *testing.T) {
	mktags := func(t ...TypedTag) []TypedTag { return t }

	cases := []struct {
		name     string
		comments []string
		group    *string
		expect   map[TagIdentifier][]TypedTag
	}{
		{
			name: "positional params",
			comments: []string{
				"+quoted(\"value\")",
				"+backticked(`value`)",
				"+ident(value)",
				"+integer(2)",
				"+negative(-5)",
				"+hex(0xFF00B3)",
				"+octal(0o04167)",
				"+binary(0b10101)",
				"+true(true)",
				"+false(false)",
			},
			expect: map[TagIdentifier][]TypedTag{
				{Name: "quoted"}: mktags(
					TypedTag{Name: "quoted", Args: []Arg{
						{Value: "value"},
					}},
				),
				{Name: "backticked"}: mktags(
					TypedTag{Name: "backticked", Args: []Arg{
						{Value: "value"},
					}},
				),
				{Name: "ident"}: mktags(
					TypedTag{Name: "ident", Args: []Arg{
						{Value: "value"},
					}},
				),
				{Name: "integer"}: mktags(
					TypedTag{Name: "integer", Args: []Arg{
						{Value: int64(2)},
					}}),
				{Name: "negative"}: mktags(
					TypedTag{Name: "negative", Args: []Arg{
						{Value: int64(-5)},
					}}),
				{Name: "hex"}: mktags(
					TypedTag{Name: "hex", Args: []Arg{
						{Value: int64(0xFF00B3)},
					}}),
				{Name: "octal"}: mktags(
					TypedTag{Name: "octal", Args: []Arg{
						{Value: int64(0o04167)},
					}}),
				{Name: "binary"}: mktags(
					TypedTag{Name: "binary", Args: []Arg{
						{Value: int64(0b10101)},
					}}),
				{Name: "true"}: mktags(
					TypedTag{Name: "true", Args: []Arg{
						{Value: true},
					}}),
				{Name: "false"}: mktags(
					TypedTag{Name: "false", Args: []Arg{
						{Value: false},
					}}),
			},
		},
		{
			name: "named params",
			comments: []string{
				"+strings(q: \"value\", b: `value`, i: value)",
				"+numbers(n1: 2, n2: -5, n3: 0xFF00B3, n4: 0o04167, n5: 0b10101)",
				"+bools(t: true, f:false)",
			},
			expect: map[TagIdentifier][]TypedTag{
				{Name: "strings"}: mktags(
					TypedTag{Name: "strings", Args: []Arg{
						{Name: "q", Value: "value"},
						{Name: "b", Value: `value`},
						{Name: "i", Value: "value"},
					}}),
				{Name: "numbers"}: mktags(
					TypedTag{Name: "numbers", Args: []Arg{
						{Name: "n1", Value: int64(2)},
						{Name: "n2", Value: int64(-5)},
						{Name: "n3", Value: int64(0xFF00B3)},
						{Name: "n4", Value: int64(0o04167)},
						{Name: "n5", Value: int64(0b10101)},
					}}),
				{Name: "bools"}: mktags(
					TypedTag{Name: "bools", Args: []Arg{
						{Name: "t", Value: true},
						{Name: "f", Value: false},
					}}),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ExtractAndParse("+", tc.group, tc.comments)
			if err != nil {
				t.Fatalf("case %q: unexpected error: %v", tc.name, err)
			}
			if !reflect.DeepEqual(result, tc.expect) {
				t.Errorf("case %q: wrong result:\n%v", tc.name, cmp.Diff(tc.expect, result))
			}
		})
	}
}

func TestParseTagKey(t *testing.T) {
	mkss := func(s ...string) []Arg {
		var args []Arg
		for _, v := range s {
			args = append(args, Arg{Value: v})
		}
		return args
	}

	cases := []struct {
		input      string
		expectKey  string
		expectArgs []Arg
		err        bool
	}{
		{"simple", "simple", nil, false},
		{"parens()", "parens", nil, false},
		{"withArgLower(arg)", "withArgLower", mkss("arg"), false},
		{"withArgUpper(ARG)", "withArgUpper", mkss("ARG"), false},
		{"withArgMixed(ArG)", "withArgMixed", mkss("ArG"), false},
		{"withArgs(arg1, arg2)", "", nil, true},
		{"argWithDash(arg-name) ", "", nil, true},
		{"withRaw(`a = b`)", "withRaw", mkss("a = b"), false},
		{"badRaw(missing`)", "", nil, true},
		{"badMix(arg,`raw`)", "", nil, true},
		{`quoted(s: "value \" \\")`, "quoted", []Arg{
			{Name: "s", Value: "value \" \\"},
		}, false},
		{"backticks(s: `value`)", "backticks", []Arg{
			{Name: "s", Value: `value`},
		}, false},
		{"ident(k: value)", "ident", []Arg{
			{Name: "k", Value: "value"},
		}, false},
		{"numbers(n1: 2, n2: -5, n3: 0xFF00B3, n4: 0o04167, n5: 0b10101)", "numbers", []Arg{
			{Name: "n1", Value: int64(2)},
			{Name: "n2", Value: int64(-5)},
			{Name: "n3", Value: int64(0xFF00B3)},
			{Name: "n4", Value: int64(0o04167)},
			{Name: "n5", Value: int64(0b10101)},
		}, false},
		{"bools(t: true, f:false)", "bools", []Arg{
			{Name: "t", Value: true},
			{Name: "f", Value: false},
		}, false},
		{"mixed(s: `value`, i: 2, b: true)", "mixed", []Arg{
			{Name: "s", Value: "value"},
			{Name: "i", Value: int64(2)},
			{Name: "b", Value: true},
		}, false},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := parseTagKey(tc.input)
			if err != nil && tc.err == false {
				t.Errorf("[%q]: expected success, got: %v", tc.input, err)
				return
			}
			if err == nil {
				if tc.err == true {
					t.Errorf("[%q]: expected failure, got: %v(%v)", tc.input, parsed.name, parsed.args)
					return
				}
				if parsed.name != tc.expectKey {
					t.Errorf("[%q]\nexpected key: %q, got: %q", tc.input, tc.expectKey, parsed.name)
				}
				if len(parsed.args) != len(tc.expectArgs) {
					t.Errorf("[%q]: expected %d args, got: %q", tc.input, len(tc.expectArgs), parsed.args)
					return
				}
				for i := range tc.expectArgs {
					if want, got := tc.expectArgs[i], parsed.args[i]; got != want {
						t.Errorf("[%q]\nexpected %q, got %q", tc.input, want, got)
					}
				}
			}
		})
	}
}
