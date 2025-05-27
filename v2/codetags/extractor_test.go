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
	tests := []struct {
		name   string
		prefix string
		lines  []string
		want   map[string][]string
	}{
		{
			name:   "prefix match",
			prefix: "+k8s:",
			lines: []string{
				"Comment line without marker",
				"+k8s:required",
				"+listType=set",
				"+k8s:format=k8s-long-name",
			},
			want: map[string][]string{
				"required": {"required"},
				"format":   {"format=k8s-long-name"},
			},
		},
		{
			name:   "empty lines",
			prefix: "+k8s:",
			lines:  []string{},
			want:   map[string][]string{},
		},
		{
			name:   "no matching lines",
			prefix: "+k8s:",
			lines: []string{
				"Comment line without marker",
				"Another comment line",
			},
			want: map[string][]string{},
		},
		{
			name:   "different marker",
			prefix: "@k8s:",
			lines: []string{
				"Comment line without marker",
				"@k8s:required",
				"@validation:required",
				"+k8s:format=k8s-long-name",
			},
			want: map[string][]string{
				"required": {"required"},
			},
		},
		{
			name:   "no group",
			prefix: "+",
			lines: []string{
				"+k8s:required",
				"+validation:required",
				"+validation:format=special",
			},
			want: map[string][]string{
				"k8s:required":        {"k8s:required"},
				"validation:required": {"validation:required"},
				"validation:format":   {"validation:format=special"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Extract(tt.prefix, tt.lines)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got:\n%v\nwant:\n%v\n", got, tt.want)
			}
		})
	}
}

func TestExtractAndParse(t *testing.T) {
	mktags := func(t ...TypedTag) []TypedTag { return t }

	cases := []struct {
		name     string
		comments []string
		expect   map[string][]TypedTag
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
			expect: map[string][]TypedTag{
				"quoted": mktags(
					TypedTag{Name: "quoted", Args: []Arg{
						{Value: ArgString("value")},
					}},
				),
				"backticked": mktags(
					TypedTag{Name: "backticked", Args: []Arg{
						{Value: ArgString("value")},
					}},
				),
				"ident": mktags(
					TypedTag{Name: "ident", Args: []Arg{
						{Value: ArgString("value")},
					}},
				),
				"integer": mktags(
					TypedTag{Name: "integer", Args: []Arg{
						{Value: MustArgInt("2")},
					}},
				),
				"negative": mktags(
					TypedTag{Name: "negative", Args: []Arg{
						{Value: MustArgInt("-5")},
					}},
				),
				"hex": mktags(
					TypedTag{Name: "hex", Args: []Arg{
						{Value: MustArgInt("0xFF00B3")},
					}},
				),
				"octal": mktags(
					TypedTag{Name: "octal", Args: []Arg{
						{Value: MustArgInt("0o04167")},
					}},
				),
				"binary": mktags(
					TypedTag{Name: "binary", Args: []Arg{
						{Value: MustArgInt("0b10101")},
					}},
				),
				"true": mktags(
					TypedTag{Name: "true", Args: []Arg{
						{Value: ArgBool(true)},
					}},
				),
				"false": mktags(
					TypedTag{Name: "false", Args: []Arg{
						{Value: ArgBool(false)},
					}},
				),
			},
		},
		{
			name: "named params",
			comments: []string{
				"+strings(q: \"value\", b: `value`, i: value)",
				"+numbers(n1: 2, n2: -5, n3: 0xFF00B3, n4: 0o04167, n5: 0b10101)",
				"+bools(t: true, f:false)",
			},
			expect: map[string][]TypedTag{
				"strings": mktags(
					TypedTag{Name: "strings", Args: []Arg{
						{Name: "q", Value: ArgString("value")},
						{Name: "b", Value: ArgString(`value`)},
						{Name: "i", Value: ArgString("value")},
					}}),
				"numbers": mktags(
					TypedTag{Name: "numbers", Args: []Arg{
						{Name: "n1", Value: MustArgInt("2")},
						{Name: "n2", Value: MustArgInt("-5")},
						{Name: "n3", Value: MustArgInt("0xFF00B3")},
						{Name: "n4", Value: MustArgInt("0o04167")},
						{Name: "n5", Value: MustArgInt("0b10101")},
					}}),
				"bools": mktags(
					TypedTag{Name: "bools", Args: []Arg{
						{Name: "t", Value: ArgBool(true)},
						{Name: "f", Value: ArgBool(false)},
					}}),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := map[string][]TypedTag{}
			for name, matchedTags := range Extract("+", tc.comments) {
				parsed, err := ParseAll(matchedTags)
				if err != nil {
					t.Fatalf("case %q: unexpected error: %v", tc.name, err)
				}
				out[name] = parsed
			}
			if !reflect.DeepEqual(out, tc.expect) {
				t.Errorf("case %q: wrong result:\n%v", tc.name, cmp.Diff(tc.expect, out))
			}
		})
	}
}
